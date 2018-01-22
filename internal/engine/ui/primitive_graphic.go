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

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/haakenlabs/forge/internal/engine"
	"github.com/haakenlabs/forge/internal/engine/system/asset/shader"
)

var _ Primitive = &Graphic{}

type Graphic struct {
	BasePrimitive

	color       engine.Color
	maskLayer   uint8
	textureMode bool
}

func (g *Graphic) SetTexture(texture *engine.Texture2D) {
	g.material.SetTexture(0, texture)
}

func (g *Graphic) SetColor(color engine.Color) {
	g.color = color
}

func (g *Graphic) Texture() *engine.Texture2D {
	return g.material.Texture(0).(*engine.Texture2D)
}

func (g *Graphic) Color() engine.Color {
	return g.color
}

func (g *Graphic) Refresh() {
	size := g.RectTransform().Size()
	verts := MakeQuad(size.Elem())

	g.mesh.Upload(verts)
}

func (g *Graphic) Draw() {
	if g.material == nil || g.mesh.size == 0 {
		return
	}

	g.textureMode = g.material.Texture(0) != nil

	g.material.Bind()
	g.mesh.Bind()

	g.material.SetProperty("v_ortho_matrix", engine.GetWindow().OrthoMatrix())
	g.material.SetProperty("v_model_matrix", g.GetTransform().ActiveMatrix())
	g.material.SetProperty("f_texture_mode", g.textureMode)
	g.material.SetProperty("f_alpha", float32(1.0))
	g.material.SetProperty("f_color", g.color.Vec4())

	gl.StencilFunc(gl.ALWAYS, int32(g.maskLayer), 0xFF)
	gl.StencilMask(0)

	g.mesh.Draw()

	g.mesh.Unbind()
	g.material.Unbind()
}

func NewGraphic() *Graphic {
	g := &Graphic{
		color: engine.ColorWhite,
	}

	g.SetName("UIGraphic")
	engine.GetInstance().MustAssign(g)

	g.material = engine.NewMaterial()
	g.material.SetShader(shader.MustGet("ui/basic"))

	g.mesh = NewMesh()
	g.mesh.Alloc()

	return g
}
