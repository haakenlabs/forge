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
	"io/ioutil"
	"path/filepath"
	"sync"
)

const (
	AssetNameShader = "shader"
)

var _ AssetHandler = &ShaderHandler{}

type ShaderHandler struct {
	BaseAssetHandler
}

type ShaderMetadata struct {
	Name     string   `json:"name"`
	Deferred bool     `json:"deferred"`
	Files    []string `json:"files"`
}

// Load will load data from the reader.
func (h *ShaderHandler) Load(r *Resource) error {
	m := &ShaderMetadata{}
	s := NewShader()

	data, err := ioutil.ReadAll(r.Reader())
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	name := m.Name
	if _, dup := h.Items[name]; dup {
		return ErrAssetExists(name)
	}

	s.SetName(m.Name)
	s.deferredCapable = m.Deferred

	// Populate shader data.
	for i := range m.Files {
		r, err := NewResource(filepath.Join(r.DirPrefix(), m.Files[i]))
		if err != nil {
			return err
		}
		if err := GetAsset().ReadResource(r); err != nil {
			return err
		}

		s.AddData(r.Bytes())
	}

	return h.Add(name, s)
}

func (h *ShaderHandler) Add(name string, shader *Shader) error {
	if _, dup := h.Items[name]; dup {
		return ErrAssetExists(name)
	}

	if err := shader.Alloc(); err != nil {
		return err
	}

	h.Items[name] = shader.ID()

	return nil
}

// Get gets an asset by name.
func (h *ShaderHandler) Get(name string) (*Shader, error) {
	a, err := h.GetAsset(name)
	if err != nil {
		return nil, err
	}

	a2, ok := a.(*Shader)
	if !ok {
		return nil, ErrAssetType(name)
	}

	return a2, nil
}

// MustGet is like GetAsset, but panics if an error occurs.
func (h *ShaderHandler) MustGet(name string) *Shader {
	a, err := h.Get(name)
	if err != nil {
		panic(err)
	}

	return a
}

func (h *ShaderHandler) Name() string {
	return AssetNameShader
}

func NewShaderHandler() *ShaderHandler {
	h := &ShaderHandler{}
	h.Items = make(map[string]uint32)
	h.Mu = &sync.RWMutex{}

	return h
}
