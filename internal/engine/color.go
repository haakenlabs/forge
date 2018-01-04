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

	"github.com/go-gl/mathgl/mgl32"
)

const (
	ErrColorParse = Error("color parse error")
)

type Color struct {
	R, G, B, A float32
}

func NewColorRGBA(v mgl32.Vec4) Color {
	c := Color{
		R: v.X(),
		G: v.Y(),
		B: v.Z(),
		A: v.W(),
	}

	return c
}

func NewColorRGB(v mgl32.Vec3) Color {
	c := Color{
		R: v.X(),
		G: v.Y(),
		B: v.Z(),
		A: 1.0,
	}

	return c
}

func NewColorRGBAHex(value string) (Color, error) {
	c := Color{}

	return c, ErrNotImplemented
}

func NewColorRGBHex(value string) (Color, error) {
	c := Color{}

	return c, ErrNotImplemented
}

func (c Color) RGBAHex() string {
	rgba32 := int32(c.R*255.0) << 24
	rgba32 += int32(c.G*255.0) << 16
	rgba32 += int32(c.B*255.0) << 8
	rgba32 += int32(c.A * 255.0)

	return fmt.Sprintf("%8X", rgba32)
}

func (c Color) RGBHex() string {
	return c.RGBAHex()[:6]
}

func (c Color) Vec3() mgl32.Vec3 {
	return mgl32.Vec3{c.R, c.G, c.B}
}

func (c Color) Vec4() mgl32.Vec4 {
	return mgl32.Vec4{c.R, c.G, c.B, c.A}
}

func (c Color) Elem() (float32, float32, float32, float32) {
	return c.R, c.G, c.B, c.A
}

func ColorBlack() Color {
	return NewColorRGB(mgl32.Vec3{0, 0, 0})
}

func ColorBlue() Color {
	return NewColorRGB(mgl32.Vec3{0, 0, 1})
}

func ColorClear() Color {
	return NewColorRGBA(mgl32.Vec4{0, 0, 0, 0})
}

func ColorCyan() Color {
	return NewColorRGB(mgl32.Vec3{0, 1, 1})
}

func ColorGray() Color {
	return NewColorRGB(mgl32.Vec3{0.5, 0.5, 0.5})
}

func ColorGreen() Color {
	return NewColorRGB(mgl32.Vec3{0, 1, 0})
}

func ColorMagenta() Color {
	return NewColorRGB(mgl32.Vec3{1, 0, 1})
}

func ColorRed() Color {
	return NewColorRGB(mgl32.Vec3{1, 0, 0})
}

func ColorWhite() Color {
	return NewColorRGB(mgl32.Vec3{1, 1, 1})
}

func ColorYellow() Color {
	return NewColorRGB(mgl32.Vec3{1, 0.92, 0.016})
}

func ColorOrange() Color {
	return NewColorRGB(mgl32.Vec3{1, 0.5, 0})
}

func ColorPurple() Color {
	return NewColorRGB(mgl32.Vec3{0.5, 0, 0.5})
}

func ColorIron() Color {
	return NewColorRGB(mgl32.Vec3{0.56, 0.57, 0.58})
}

func ColorCopper() Color {
	return NewColorRGB(mgl32.Vec3{0.95, 0.64, 0.54})
}

func ColorGold() Color {
	return NewColorRGB(mgl32.Vec3{1.00, 0.71, 0.29})
}

func ColorAluminium() Color {
	return NewColorRGB(mgl32.Vec3{0.91, 0.92, 0.92})
}

func ColorSilver() Color {
	return NewColorRGB(mgl32.Vec3{0.95, 0.93, 0.88})
}
