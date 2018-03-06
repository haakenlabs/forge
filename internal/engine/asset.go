/*
Copyright (c) 2017 HaakenLabs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/haakenlabs/forge/internal/builtin"
)

var _ System = &Asset{}

// ErrAssetNotFound reports that the asset was not found in the handler.
type ErrAssetNotFound string

func (e ErrAssetNotFound) Error() string {
	return "asset: no such asset: " + string(e)
}

// ErrAssetExists reports that the asset was already found in the handler.
type ErrAssetExists string

func (e ErrAssetExists) Error() string {
	return "asset: duplicate asset: " + string(e)
}

// ErrAssetType reports that the asset has an invalid underlying type.
type ErrAssetType string

func (e ErrAssetType) Error() string {
	return "asset: type assertion error for asset: " + string(e)
}

// ErrAssetNotFound reports that the handler is not registered.
type ErrHandlerNotFound string

func (e ErrHandlerNotFound) Error() string {
	return "asset: no such handler: " + string(e)
}

const SysNameAsset = "asset"

type AssetManifest struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Assets      map[string][]string `json:"assets,required"`
}

type AssetMetadata struct {
	Name  string   `json:"name,required"`
	Type  string   `json:"type,required"`
	Files []string `json:"files,required"`
}

type AssetHandler interface {
	// Load will allocate the asset. If OpenGL is required for allocation, this
	// function should not be called from a goroutine.
	Load(*Resource) error

	// GetAsset gets an asset by name.
	GetAsset(string) (Object, error)

	// MustGetAsset is like GetAsset, but panics if an error occurs.
	MustGetAsset(string) Object

	// Name returns the name of this handler.
	Name() string

	// Count returns the number of assets tracked by this handler.
	Count() int
}

type BaseAssetHandler struct {
	Items map[string]uint32
	Mu    *sync.RWMutex
}

type Asset struct {
	handlers map[string]AssetHandler
	packages map[string]*Package
	mu       *sync.RWMutex
}

// Setup sets up the System.
func (a *Asset) Setup() error {
	return nil
}

// Teardown tears down the System.
func (a *Asset) Teardown() {
	a.ReleaseAll()
	a.UnmountAllPackages()
}

// Name returns the name of the System.
func (a *Asset) Name() string {
	return SysNameAsset
}

// MountPackage mounts a new package by name.
func (a *Asset) MountPackage(name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, dup := a.packages[name]; dup {
		return ErrPackageMounted(name)
	}

	p := NewPackage(name)
	if err := p.Mount(); err != nil {
		return err
	}

	a.packages[name] = p

	return nil
}

// UnmountPackage unmounts a mounted package given by name.
func (a *Asset) UnmountPackage(name string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, ok := a.packages[name]; !ok {
		return ErrPackageNotMounted(name)
	}

	if err := a.packages[name].Unmount(); err != nil {
		return err
	}

	delete(a.packages, name)

	return nil
}

// UnmountAllPackages unmounts all mounted packages.
func (a *Asset) UnmountAllPackages() {
	for p := range a.packages {
		if err := a.UnmountPackage(p); err != nil {
			logrus.Error(err)
		}
	}
}

// Get gets an asset by name from a handler by kind.
func (a *Asset) Get(kind, name string) (Object, error) {
	return a.GetAsset(kind, name)
}

// MustGet is like Get, but panics if an error is encountered.
func (a *Asset) MustGet(kind, name string) Object {
	if a, err := a.Get(kind, name); err != nil {
		panic(err)
	} else {
		return a
	}
}

// LoadManifest loads a manifest of assets.
func (a *Asset) LoadManifest(files ...string) error {
	for _, v := range files {
		m := NewAssetManifest()

		r, err := NewResource(v)
		if err != nil {
			logrus.Error("Error reading manifest: ", err)
			continue
		}
		if err := a.ReadResource(r); err != nil {
			logrus.Error("Error reading manifest: ", err)
			continue
		}

		if err := json.Unmarshal(r.Bytes(), m); err != nil {
			return err
		}

		logrus.Debug(m.Assets)

		// Load assets.
		for t := range m.Assets {
			h, err := a.GetHandler(t)
			if err != nil {
				logrus.Error(err)
				continue
			}

			logrus.Debug(">>>>>>>> Loading asset type: ", t)

			// Read and load assets.
			for n := range m.Assets[t] {
				ar, err := NewResource(path.Join(r.DirPrefix(), m.Assets[t][n]))
				if err != nil {
					return err
				}

				if err := a.ReadResource(ar); err != nil {
					return err
				}

				logrus.Debug("Read asset: ", m.Assets[t][n])

				if err := h.Load(ar); err != nil {
					return err
				}

				logrus.Debug("Loaded asset: ", m.Assets[t][n])
			}
		}
	}

	return nil
}

func (a *Asset) ReadResource(r *Resource) error {
	if r == nil {
		return nil
	}

	switch r.resType {
	case ResourceFile:
		f, err := os.Open(r.location)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(r.buffer, f)

		return err
	case ResourcePackage:
		p, ok := a.packages[r.container]
		if !ok {
			return ErrPackageNotMounted(r.container)
		}

		return p.Read(r.location, r.buffer)
	case ResourceBindata:
		data, err := builtin.Asset(r.location)
		if err != nil {
			panic(err)
			return err
		}

		_, err = r.buffer.Write(data)

		return err
	default:
		return fmt.Errorf("resource: unknown resource type for resource: %d", int(r.resType))
	}
}

// Register registers an asset handler.
func (a *Asset) RegisterHandler(h AssetHandler) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if _, dup := a.handlers[h.Name()]; dup {
		return ErrHandlerNotFound(h.Name())
	}

	a.handlers[h.Name()] = h

	logrus.Debug("registered handler: ", h.Name())

	return nil
}

// HandlerRegistered returns true if there is a handler registered with the given name.
func (a *Asset) HandlerRegistered(name string) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	_, ok := a.handlers[name]

	return ok
}

// GetHandler gets an asset handler by name.
func (a *Asset) GetHandler(name string) (AssetHandler, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if h, ok := a.handlers[name]; ok {
		return h, nil
	}

	return nil, ErrHandlerNotFound(name)
}

// GetAsset gets an asset by name from a handler by kind.
func (a *Asset) GetAsset(kind, name string) (Object, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var asset Object
	var err error

	if _, ok := a.handlers[kind]; !ok {
		return nil, ErrHandlerNotFound(kind)
	}

	asset, err = a.handlers[kind].GetAsset(name)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

// MustGetAsset is like GetAsset, but panics if an error occurs.
func (a *Asset) MustGetAsset(kind, name string) Object {
	asset, err := a.GetAsset(kind, name)
	if err != nil {
		panic(err)
	}

	return asset
}

// ReleaseAll releases all builtin managed by this asset store.
func (a *Asset) ReleaseAll() {

}

// Count reports the total number of assets managed by this asset store.
func (a *Asset) Count() int {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var count int

	for h := range a.handlers {
		count += a.handlers[h].Count()
	}

	return count
}

// CountByKind reports the total number of assets by kind managed by this asset store.
func (a *Asset) CountByKind(kind string) (int, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if _, ok := a.handlers[kind]; !ok {
		return 0, ErrHandlerNotFound(kind)
	}

	return a.handlers[kind].Count(), nil
}

// GetAsset gets an asset by name.
func (h *BaseAssetHandler) GetAsset(name string) (Object, error) {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	if id, ok := h.Items[name]; !ok {
		return nil, ErrAssetNotFound(name)
	} else {
		obj, err := GetInstance().Get(id)
		if err != nil {
			return nil, err
		}

		return obj, nil
	}
}

// MustGetAsset is like GetAsset, but panics if an error occurs.
func (h *BaseAssetHandler) MustGetAsset(name string) Object {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	a, err := h.GetAsset(name)
	if err != nil {
		panic(err)
	}

	return a
}

// Count returns the number of assets tracked by this handler.
func (h *BaseAssetHandler) Count() int {
	h.Mu.RLock()
	defer h.Mu.RUnlock()

	return len(h.Items)
}

func NewAsset() *Asset {
	return &Asset{
		handlers: make(map[string]AssetHandler),
		packages: make(map[string]*Package),
		mu:       &sync.RWMutex{},
	}
}

// GetAsset gets the asset system from the current app.
func GetAsset() *Asset {
	return CurrentApp().MustSystem(SysNameAsset).(*Asset)
}

func NewAssetManifest() *AssetManifest {
	m := &AssetManifest{
		Assets: make(map[string][]string),
	}

	return m
}
