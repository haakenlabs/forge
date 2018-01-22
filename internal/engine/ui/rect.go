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

package ui

import "github.com/go-gl/mathgl/mgl32"

type Rect struct {
	origin mgl32.Vec2
	size   mgl32.Vec2
}

// Rect Methods

func NewRect() Rect {
	return Rect{}
}

func NewRectFrom(origin mgl32.Vec2, size mgl32.Vec2) Rect {
	r := Rect{
		origin: origin,
		size:   size,
	}

	return r
}

func (r *Rect) SetOrigin(origin mgl32.Vec2) {
	r.origin = origin
}

func (r *Rect) SetSize(size mgl32.Vec2) {
	r.size = size
}

func (r *Rect) SetX(x float32) {
	r.origin[0] = x
}

func (r *Rect) SetY(y float32) {
	r.origin[1] = y
}

func (r *Rect) SetWidth(width float32) {
	r.size[0] = width
}

func (r *Rect) SetHeight(height float32) {
	r.size[1] = height
}

func (r *Rect) Left() float32 {
	return r.origin.X()
}

func (r *Rect) Top() float32 {
	return r.origin.Y()
}

func (r *Rect) Right() float32 {
	return r.origin.X() + r.size.X()
}

func (r *Rect) Bottom() float32 {
	return r.origin.Y() + r.size.Y()
}

func (r *Rect) Width() float32 {
	return r.size.X()
}

func (r *Rect) Height() float32 {
	return r.size.Y()
}

func (r *Rect) Origin() mgl32.Vec2 {
	return r.origin
}

func (r *Rect) OriginElem() (float32, float32) {
	return r.origin.Elem()
}

func (r *Rect) Center() mgl32.Vec2 {
	return r.origin.Add(r.size.Mul(0.5))
}

func (r *Rect) CenterElem() (float32, float32) {
	return r.Center().Elem()
}

func (r *Rect) Size() mgl32.Vec2 {
	return r.size
}

func (r *Rect) SizeElem() (float32, float32) {
	return r.size.Elem()
}

func (r *Rect) Extent() mgl32.Vec2 {
	return r.origin.Add(r.size)
}

func (r *Rect) ExtentElem() (float32, float32) {
	return r.Extent().Elem()
}

func (r *Rect) Min() mgl32.Vec2 {
	return r.Origin()
}

func (r *Rect) MinElem() (float32, float32) {
	return r.Min().Elem()
}

func (r *Rect) Max() mgl32.Vec2 {
	return r.Extent()
}

func (r *Rect) MaxElem() (float32, float32) {
	return r.Max().Elem()
}

func (r *Rect) Contains(point mgl32.Vec2) bool {
	return point.X() >= r.Left() && point.X() < r.Right() && point.Y() >= r.Top() && point.Y() < r.Bottom()
}

func (r *Rect) Intersects(rect Rect) bool {
	if r.Right() < rect.Left() {
		return false
	}
	if r.Left() > rect.Right() {
		return false
	}
	if r.Top() > rect.Bottom() {
		return false
	}
	if r.Bottom() < rect.Top() {
		return false
	}

	return true
}

func (r *Rect) Distance(point mgl32.Vec2) float32 {
	return 0
}

func (r *Rect) Distance2(point mgl32.Vec2) float32 {
	return 0
}

func (r *Rect) ExpandToContain(point mgl32.Vec2) {
	if r.Contains(point) {
		return
	}

	if point.X() < r.Left() {
		r.origin[0] = point.X()
		r.size[0] += r.origin.X() - point.X()
	}
	if point.Y() < r.Top() {
		r.origin[1] = point.Y()
		r.size[1] += r.origin.Y() - point.Y()
	}
	if point.X() > r.Right() {
		r.size[0] += point.X() - r.Right()
	}
	if point.Y() > r.Bottom() {
		r.size[1] += point.Y() - r.Bottom()
	}
}

func (r *Rect) ExpandToContainRect(rect Rect) {
	r.ExpandToContain(rect.Origin())
	r.ExpandToContain(rect.Extent())
}
