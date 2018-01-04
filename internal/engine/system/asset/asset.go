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

package asset

import (
	"github.com/haakenlabs/forge/internal/engine"
)

// RegisterHandler registers an asset handler.
func RegisterHandler(h engine.AssetHandler) error {
	return engine.GetAsset().RegisterHandler(h)
}

// GetHandler gets an asset handler by name.
func GetHandler(name string) (engine.AssetHandler, error) {
	return engine.GetAsset().GetHandler(name)
}

// Get gets an asset by name from a handler by kind.
func Get(kind, name string) (engine.Object, error) {
	return engine.GetAsset().Get(kind, name)
}

// // MustGet is like Get, but panics if an error is encountered.
func MustGet(kind, name string) engine.Object {
	return engine.GetAsset().MustGet(kind, name)
}

func MountPackage(name string) error {
	return engine.GetAsset().MountPackage(name)
}

// UnmountPackage unmounts a mounted package given by name.
func UnmountPackage(name string) error {
	return engine.GetAsset().UnmountPackage(name)
}

// UnmountAllPackages unmounts all mounted packages.
func UnmountAllPackages() {
	engine.GetAsset().UnmountAllPackages()
}

// LoadManifest loads a manifest of assets.
func LoadManifest(files ...string) error {
	return engine.GetAsset().LoadManifest(files...)
}

func ReadResource(r *engine.Resource) error {
	return engine.GetAsset().ReadResource(r)
}
