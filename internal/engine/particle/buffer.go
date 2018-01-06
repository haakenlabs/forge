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

package particle

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/sirupsen/logrus"

	"github.com/haakenlabs/forge/internal/engine"
	"github.com/haakenlabs/forge/internal/engine/system/instance"
)

type buffer struct {
	engine.BaseObject

	size         uint32
	vao          uint32
	idParticle   uint32
	idAlive      uint32
	idDead       uint32
	idIndex      uint32
	idAttractors uint32
}

func (b *buffer) Bind() {
	gl.BindVertexArray(b.vao)
}

func (b *buffer) Unbind() {
	gl.BindVertexArray(0)
}

func (b *buffer) Dealloc() {
	gl.DeleteBuffers(1, &b.idParticle)
	gl.DeleteBuffers(1, &b.idAlive)
	gl.DeleteBuffers(1, &b.idDead)
	gl.DeleteBuffers(1, &b.idIndex)
	gl.DeleteVertexArrays(1, &b.vao)
}

func (b *buffer) Alloc() error {
	gl.GenVertexArrays(1, &b.vao)
	gl.GenBuffers(1, &b.idParticle)
	gl.GenBuffers(1, &b.idAlive)
	gl.GenBuffers(1, &b.idDead)
	gl.GenBuffers(1, &b.idIndex)
	gl.GenBuffers(1, &b.idAttractors)

	b.Bind()

	// Particle Pool
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, b.idParticle)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idParticle)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(b.size)*sizeOfParticle, nil, gl.DYNAMIC_DRAW)

	// Alive List
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 1, b.idAlive)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idAlive)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(b.size*2)*4, nil, gl.DYNAMIC_DRAW)

	// Dead List
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 2, b.idDead)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idDead)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(b.size)*4, nil, gl.DYNAMIC_DRAW)

	// Indices
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 3, b.idIndex)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idIndex)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 12, nil, gl.DYNAMIC_DRAW)

	// Attractors
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 4, b.idAttractors)
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idAttractors)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 48, nil, gl.DYNAMIC_DRAW)

	b.Unbind()

	logrus.Debugf("%s allocated particle buffer", b)

	return nil
}

func (b *buffer) Reserve() {
	b.Bind()

	// Particle Pool
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idParticle)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(b.size)*sizeOfParticle, nil, gl.DYNAMIC_DRAW)

	// Alive List
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idAlive)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(b.size*2)*4, nil, gl.DYNAMIC_DRAW)

	// Dead List
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idDead)
	data := make([]uint32, b.size)
	for i := range data {
		data[i] = b.size - uint32(i) - 1
	}
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(b.size)*4, gl.Ptr(data), gl.DYNAMIC_DRAW)

	// Indices
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idIndex)
	data = []uint32{0, b.size, 0}
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 12, gl.Ptr(data), gl.DYNAMIC_DRAW)

	b.Unbind()
}

func (b *buffer) SetSize(size uint32) {
	b.size = size

	b.Reserve()
}

func (b *buffer) ResetIndex(alive, dead uint32) {
	data := []uint32{alive, dead, 0}

	b.Bind()
	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, b.idIndex)
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, 12, gl.Ptr(data), gl.DYNAMIC_DRAW)
	b.Unbind()
}

func newBuffer(size uint32) *buffer {
	b := &buffer{}
	b.size = size

	b.SetName("ParticleBuffer")
	instance.MustAssign(b)

	return b
}
