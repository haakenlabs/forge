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

package hdr

import (
	"image"
	"image/color"
	"math"
)

// RGB96 is an in-memory image whose At method returns RGB96 values.
type RGB96 struct {
	// Pix holds the image's pixels, in R, G, B, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8

	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int

	// Rect is the image's bounds.
	Rect image.Rectangle
}

var _ image.Image = &RGB96{}

func (p *RGB96) ColorModel() color.Model { return RGB96Model }

func (p *RGB96) Bounds() image.Rectangle { return p.Rect }

func (p *RGB96) At(x, y int) color.Color {
	return p.RGB96At(x, y)
}

func (p *RGB96) RGB96At(x, y int) RGB96Color {
	if !(image.Point{x, y}.In(p.Rect)) {
		return RGB96Color{}
	}

	i := p.PixOffset(x, y)

	r := uint32(p.Pix[i+0])<<24 | uint32(p.Pix[i+1])<<16 | uint32(p.Pix[i+2])<<8 | uint32(p.Pix[i+3])
	g := uint32(p.Pix[i+4])<<24 | uint32(p.Pix[i+5])<<16 | uint32(p.Pix[i+6])<<8 | uint32(p.Pix[i+7])
	b := uint32(p.Pix[i+8])<<24 | uint32(p.Pix[i+9])<<16 | uint32(p.Pix[i+10])<<8 | uint32(p.Pix[i+11])

	return RGB96Color{
		math.Float32frombits(r),
		math.Float32frombits(g),
		math.Float32frombits(b),
	}
}

func (p *RGB96) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.PixOffset(x, y)
	c1 := RGB96Model.Convert(c).(RGB96Color)

	r := math.Float32bits(c1.R)
	g := math.Float32bits(c1.G)
	b := math.Float32bits(c1.B)

	p.Pix[i+0] = uint8(r >> 24)
	p.Pix[i+1] = uint8(r >> 16)
	p.Pix[i+2] = uint8(r >> 8)
	p.Pix[i+3] = uint8(r)
	p.Pix[i+4] = uint8(g >> 24)
	p.Pix[i+5] = uint8(g >> 16)
	p.Pix[i+6] = uint8(g >> 8)
	p.Pix[i+7] = uint8(g)
	p.Pix[i+8] = uint8(b >> 24)
	p.Pix[i+9] = uint8(b >> 16)
	p.Pix[i+10] = uint8(b >> 8)
	p.Pix[i+11] = uint8(b)
}

func (p *RGB96) SetRGB96(x, y int, c RGB96Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}

	i := p.PixOffset(x, y)
	r := math.Float32bits(c.R)
	g := math.Float32bits(c.G)
	b := math.Float32bits(c.B)

	p.Pix[i+0] = uint8(r >> 24)
	p.Pix[i+1] = uint8(r >> 16)
	p.Pix[i+2] = uint8(r >> 8)
	p.Pix[i+3] = uint8(r)
	p.Pix[i+4] = uint8(g >> 24)
	p.Pix[i+5] = uint8(g >> 16)
	p.Pix[i+6] = uint8(g >> 8)
	p.Pix[i+7] = uint8(g)
	p.Pix[i+8] = uint8(b >> 24)
	p.Pix[i+9] = uint8(b >> 16)
	p.Pix[i+10] = uint8(b >> 8)
	p.Pix[i+11] = uint8(b)
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *RGB96) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*12
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *RGB96) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to be inside
	// either r1 or r2 if the intersection is empty. Without explicitly checking for
	// this, the Pix[i:] expression below can panic.
	if r.Empty() {
		return &RGB96{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &RGB96{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque. Since
// RGB96 does not have an alpha channel, it is always opaque.
func (p *RGB96) Opaque() bool {
	return true
}

// NewRGB96 returns a new RGB96 image with the given bounds.
func NewRGB96(r image.Rectangle) *RGB96 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 12*w*h)
	return &RGB96{pix, 12 * w, r}
}

// RGB96Color represents a 96-bit HDR color, having 32 bits for each of red, green, and blue.
type RGB96Color struct {
	R, G, B float32
}

var _ color.Color = &RGB96Color{}

func (c RGB96Color) RGBA() (r, g, b, a uint32) {
	r = math.Float32bits(c.R)
	g = math.Float32bits(c.G)
	b = math.Float32bits(c.B)

	return
}

var RGB96Model color.Model = color.ModelFunc(rgb96Model)

func rgb96Model(c color.Color) color.Color {
	if _, ok := c.(RGB96Color); ok {
		return c
	}

	r, g, b, _ := c.RGBA()

	return RGB96Color{math.Float32frombits(r), math.Float32frombits(g), math.Float32frombits(b)}
}
