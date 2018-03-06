/*
Copyright (c) 2018 HaakenLabs

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
	"sync"

	"github.com/golang/freetype/truetype"
)

const (
	AssetNameFont = "font"
)

var _ AssetHandler = &FontHandler{}

type FontHandler struct {
	BaseAssetHandler
}

// Load will load data from the reader.
func (h *FontHandler) Load(r *Resource) error {
	name := r.Base()

	if _, dup := h.Items[name]; dup {
		return ErrAssetExists(name)
	}

	ttf, err := truetype.Parse(r.Bytes())
	if err != nil {
		return err
	}

	f := NewFont(ttf, ASCII)
	f.SetName(name)

	return h.Add(name, f)
}

func (h *FontHandler) Add(name string, font *Font) error {
	if _, dup := h.Items[name]; dup {
		return ErrAssetExists(name)
	}

	if err := font.Alloc(); err != nil {
		return err
	}

	h.Items[name] = font.ID()

	return nil
}

// Get gets an asset by name.
func (h *FontHandler) Get(name string) (*Font, error) {
	a, err := h.GetAsset(name)
	if err != nil {
		return nil, err
	}

	a2, ok := a.(*Font)
	if !ok {
		return nil, ErrAssetType(name)
	}

	return a2, nil
}

// MustGet is like GetAsset, but panics if an error occurs.
func (h *FontHandler) MustGet(name string) *Font {
	a, err := h.Get(name)
	if err != nil {
		panic(err)
	}

	return a
}

func (h *FontHandler) Name() string {
	return AssetNameFont
}

func NewFontHandler() *FontHandler {
	h := &FontHandler{}
	h.Items = make(map[string]uint32)
	h.Mu = &sync.RWMutex{}

	return h
}
