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

type GBuffer struct {
	Framebuffer

	hdr bool
}

func NewGBuffer(size math.IVec2, depth *AttachmentTexture2D, hdr bool) *GBuffer {
	g := &GBuffer{
		hdr: hdr,
	}

	g.size = size
	g.attachments = make(map[uint32]Attachment)
	g.drawBuffers = []uint32{}

	g.SetName("GBuffer")
	GetInstance().MustAssign(g)

	gl.GenFramebuffers(1, &g.reference)

	attachment0 := NewAttachmentTexture2D(g.size, TextureFormatRGBA32)
	attachment1 := NewAttachmentTexture2D(g.size, TextureFormatRGBA32UI)

	var attachment2 *AttachmentTexture2D
	if hdr {
		attachment2 = NewAttachmentTexture2D(g.size, TextureFormatRGBA16)
	} else {
		attachment2 = NewAttachmentTexture2D(g.size, TextureFormatDefaultColor)
	}

	attachment1.AttachmentObject().Bind()
	attachment1.AttachmentObject().SetFilter(gl.NEAREST, gl.NEAREST)

	g.SetAttachment(gl.COLOR_ATTACHMENT0, attachment0)
	g.SetAttachment(gl.COLOR_ATTACHMENT1, attachment1)
	g.SetAttachment(gl.COLOR_ATTACHMENT2, attachment2)
	g.SetAttachment(gl.DEPTH_ATTACHMENT, depth)

	g.SetDrawBuffers([]uint32{gl.COLOR_ATTACHMENT0, gl.COLOR_ATTACHMENT1, gl.COLOR_ATTACHMENT2})

	return g
}

func (g *GBuffer) HDR() bool {
	return g.hdr
}

func (g *GBuffer) SetHDR(enable bool) {
	if g.hdr != enable {
		g.RemoveAttachment(gl.COLOR_ATTACHMENT2)

		var attachment2 *AttachmentTexture2D

		if enable {
			attachment2 = NewAttachmentTexture2D(g.size, TextureFormatRGBA16)
		} else {
			attachment2 = NewAttachmentTexture2D(g.size, TextureFormatDefaultColor)
		}

		g.SetAttachment(gl.COLOR_ATTACHMENT2, attachment2)
	}

	// TODO: Handle error
	g.Alloc()

	g.hdr = enable
}

func (g *GBuffer) Attachment0() *Texture2D {
	if a, ok := g.GetAttachment(gl.COLOR_ATTACHMENT0).(*AttachmentTexture2D); ok {
		return a.AttachmentObject()
	}

	return nil
}

func (g *GBuffer) Attachment1() *Texture2D {
	if a, ok := g.GetAttachment(gl.COLOR_ATTACHMENT1).(*AttachmentTexture2D); ok {
		return a.AttachmentObject()
	}

	return nil
}

func (g *GBuffer) AttachmentDepth() *Texture2D {
	if a, ok := g.GetAttachment(gl.DEPTH_ATTACHMENT).(*AttachmentTexture2D); ok {
		return a.AttachmentObject()
	}

	return nil
}
