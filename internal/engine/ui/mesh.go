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
	"github.com/go-gl/mathgl/mgl32"
	"github.com/haakenlabs/forge/internal/engine"
)

type Mesh struct {
	engine.BaseObject

	size int32
	vao  uint32
	vbo  uint32
}

func (m *Mesh) Alloc() error {
	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 32, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 32, gl.PtrOffset(12))
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 32, gl.PtrOffset(24))

	gl.BufferData(gl.ARRAY_BUFFER, 32, nil, gl.DYNAMIC_DRAW)

	m.Unbind()

	return nil
}

func (m *Mesh) Dealloc() {
	gl.DeleteBuffers(1, &m.vbo)
	gl.DeleteVertexArrays(1, &m.vao)
}

func (m *Mesh) Bind() {
	gl.BindVertexArray(m.vao)
}

func (m *Mesh) Unbind() {
	gl.BindVertexArray(0)
}

func (m *Mesh) Upload(vertices []engine.Vertex) {
	m.size = int32(len(vertices))

	m.Bind()
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	if m.size == 0 {
		gl.BufferData(gl.ARRAY_BUFFER, 0, nil, gl.DYNAMIC_DRAW)
	} else {
		gl.BufferData(gl.ARRAY_BUFFER, int(m.size*32), gl.Ptr(vertices), gl.DYNAMIC_DRAW)
	}

	m.Unbind()
}

func (m *Mesh) Draw() {
	if m.size <= 0 {
		return
	}

	gl.DrawArrays(gl.TRIANGLES, 0, m.size)
}

func NewMesh() *Mesh {
	m := &Mesh{}

	m.SetName("UIMesh")
	engine.GetInstance().MustAssign(m)

	return m
}

func MakeQuad(w, h float32) []engine.Vertex {
	ul := engine.Vertex{V: mgl32.Vec3{}, U: mgl32.Vec2{0, 1}}
	ur := engine.Vertex{V: mgl32.Vec3{w, 0, 0}, U: mgl32.Vec2{1, 1}}
	lr := engine.Vertex{V: mgl32.Vec3{w, h, 0}, U: mgl32.Vec2{1, 0}}
	ll := engine.Vertex{V: mgl32.Vec3{0, h, 0}, U: mgl32.Vec2{0, 0}}

	return []engine.Vertex{ul, lr, ur, ul, ll, lr}
}
