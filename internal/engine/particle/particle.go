package particle

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"

	"math"

	"github.com/haakenlabs/forge/internal/engine"
	"github.com/haakenlabs/forge/internal/engine/system/instance"
)

const (
	sizeOfParticle = 80
)

type particle struct {
	Color           engine.Color
	AngularVelocity mgl32.Vec4
	Rotation        mgl32.Vec4
	Position        mgl32.Vec4
	Lifecycle       mgl32.Vec4
}

type buffer struct {
	engine.BaseObject

	vao  uint32
	ssbo uint32
	size uint32
}

func (b *buffer) Bind() {
	gl.BindVertexArray(b.ssbo)
}

func (b *buffer) Unbind() {
	gl.BindVertexArray(0)
}

func (b *buffer) Alloc() error {
	gl.GenVertexArrays(1, &b.vao)
	gl.GenBuffers(1, &b.ssbo)

	gl.BindVertexArray(b.vao)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, b.ssbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.ssbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(b.size)*sizeOfParticle, nil, gl.DYNAMIC_DRAW)

	// start_color
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, sizeOfParticle, gl.PtrOffset(0))
	// angular_velocity
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, sizeOfParticle, gl.PtrOffset(16))
	// rotation
	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 4, gl.FLOAT, false, sizeOfParticle, gl.PtrOffset(32))
	// position
	gl.EnableVertexAttribArray(3)
	gl.VertexAttribPointer(3, 4, gl.FLOAT, false, sizeOfParticle, gl.PtrOffset(48))
	// lifecycle
	gl.EnableVertexAttribArray(4)
	gl.VertexAttribPointer(4, 4, gl.FLOAT, false, sizeOfParticle, gl.PtrOffset(64))

	gl.BindVertexArray(0)

	logrus.Debugf("%s allocated particle buffer", b)

	return nil
}

func (b *buffer) Dealloc() {
	gl.DeleteBuffers(1, &b.ssbo)
	gl.DeleteVertexArrays(1, &b.vao)
}

func (b *buffer) Upload(data []particle) {
	if len(data) > math.MaxUint32 {
		return
	}

	b.size = uint32(len(data))

	b.Bind()
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(b.size)*sizeOfParticle, gl.Ptr(data), gl.DYNAMIC_DRAW)
	b.Unbind()
}

func newBuffer(size uint32) *buffer {
	b := &buffer{
		size: size,
	}

	b.SetName("ParticleBuffer")
	instance.MustAssign(b)

	return b
}
