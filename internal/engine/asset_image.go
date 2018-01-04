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
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"sync"

	"git.dbnservers.net/haakenlabs/forge/internal/math"

	_ "git.dbnservers.net/haakenlabs/forge/internal/image/hdr"
	_ "image/jpeg"
	_ "image/png"
)

const (
	AssetNameImage = "image"
)

var _ AssetHandler = &ImageHandler{}

type ImageHandler struct {
	BaseAssetHandler
}

// Load will load data from the reader.
func (h *ImageHandler) Load(r *Resource) error {
	var texture *Texture2D
	var img image.Image

	name := r.Base()

	if _, dup := h.Items[r.Base()]; dup {
		return ErrAssetExists(name)
	}

	img, _, err := image.Decode(r.Reader())
	if err != nil {
		return err
	}

	x := int32(img.Bounds().Dx())
	y := int32(img.Bounds().Dy())

	texture = NewTexture2D(math.IVec2{x, y}, TextureFormatDefaultColor)

	switch img.ColorModel() {
	// 4 channels, 16 bits per channel
	case color.RGBA64Model:
		rgba := image.NewRGBA64(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		texture.SetTexFormat(TextureFormatRGBA16)
		texture.SetData(rgba.Pix)
		// 4 channels, 8 bits per channel
	case color.RGBAModel:
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		texture.SetTexFormat(TextureFormatRGBA8)
		texture.SetData(rgba.Pix)
		// 2 channels, 16 bits per channel
	case color.Alpha16Model:
		alpha := image.NewAlpha16(img.Bounds())
		draw.Draw(alpha, alpha.Bounds(), img, image.Point{}, draw.Src)
		texture.SetTexFormat(TextureFormatRG16)
		texture.SetData(alpha.Pix)
		// 2 channels, 8 bits per channel
	case color.AlphaModel:
		alpha := image.NewAlpha(img.Bounds())
		draw.Draw(alpha, alpha.Bounds(), img, image.Point{}, draw.Src)
		texture.SetTexFormat(TextureFormatRG8)
		texture.SetData(alpha.Pix)
		// 1 channel, 16 bits per channel
	case color.Gray16Model:
		gray := image.NewGray16(img.Bounds())
		draw.Draw(gray, gray.Bounds(), img, image.Point{}, draw.Src)
		texture.SetTexFormat(TextureFormatR16)
		texture.SetData(gray.Pix)
		// 1 channel, 16 bits per channel
	case color.GrayModel:
		gray := image.NewGray(img.Bounds())
		draw.Draw(gray, gray.Bounds(), img, image.Point{}, draw.Src)
		texture.SetTexFormat(TextureFormatR8)
		texture.SetData(gray.Pix)
	case color.NRGBA64Model:
		rgba := image.NewNRGBA64(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		texture.SetTexFormat(TextureFormatRGBA16)
		texture.SetData(rgba.Pix)
	case color.NRGBAModel:
		rgba := image.NewNRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		texture.SetTexFormat(TextureFormatRGBA8)
		texture.SetData(rgba.Pix)
	default:
		return fmt.Errorf("invalid color format: %v", img.ColorModel())
	}

	return h.Add(name, texture)
}

func (h *ImageHandler) Add(name string, texture *Texture2D) error {
	if _, dup := h.Items[name]; dup {
		return ErrAssetExists(name)
	}

	if err := texture.Alloc(); err != nil {
		return err
	}

	h.Items[name] = texture.ID()

	return nil
}

// Get gets an asset by name.
func (h *ImageHandler) Get(name string) (*Texture2D, error) {
	a, err := h.GetAsset(name)
	if err != nil {
		return nil, err
	}

	a2, ok := a.(*Texture2D)
	if !ok {
		return nil, ErrAssetType(name)
	}

	return a2, nil
}

// MustGet is like GetAsset, but panics if an error occurs.
func (h *ImageHandler) MustGet(name string) *Texture2D {
	a, err := h.Get(name)
	if err != nil {
		panic(err)
	}

	return a
}

func (h *ImageHandler) Name() string {
	return AssetNameImage
}

func NewImageHandler() *ImageHandler {
	h := &ImageHandler{}
	h.Items = make(map[string]uint32)
	h.Mu = &sync.RWMutex{}

	return h
}
