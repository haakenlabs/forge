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
	"github.com/go-gl/gl/v4.3-core/gl"

	"github.com/haakenlabs/forge/internal/math"
)

type RenderBuffer struct {
	BaseObject

	size           math.IVec2
	reference      uint32
	internalFormat uint32
}

func NewRenderBuffer(size math.IVec2, format TextureFormat) *RenderBuffer {
	return NewRenderBufferIntFmt(size, uint32(TextureFormatToInternal(format)))
}

func NewRenderBufferIntFmt(size math.IVec2, internalFormat uint32) *RenderBuffer {
	r := &RenderBuffer{
		size:           size,
		internalFormat: internalFormat,
	}

	r.SetName("RenderBuffer")
	GetInstance().MustAssign(r)

	gl.GenRenderbuffers(1, &r.reference)

	r.Allocate()

	return r
}

func (r *RenderBuffer) Release() {
	if r.reference != 0 {
		gl.DeleteRenderbuffers(1, &r.reference)
		r.reference = 0
	}
}

func (r *RenderBuffer) Reference() uint32 {
	return r.reference
}

func (r *RenderBuffer) Allocate() {
	gl.BindRenderbuffer(gl.RENDERBUFFER, r.reference)
	gl.RenderbufferStorage(gl.RENDERBUFFER, r.internalFormat, r.size.X(), r.size.Y())
}

func (r *RenderBuffer) Attach(location uint32) {
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, location, gl.RENDERBUFFER, r.reference)
}

func (r *RenderBuffer) SetSize(size math.IVec2) {
	r.size = size
	r.Allocate()
}

func (r *RenderBuffer) Size() math.IVec2 {
	return r.size
}
