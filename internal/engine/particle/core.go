package particle

import (
	"encoding/binary"
	"unsafe"

	"github.com/go-gl/gl/v4.3-core/gl"

	"github.com/haakenlabs/forge/internal/engine"
)

type ModuleCore struct {
	StartColor    engine.Color
	Duration      float32
	StartDelay    float32
	StartLifetime float32
	StartSpeed    float32
	StartSize     float32
	PlaybackSpeed float32
	RandomSeed    uint32
	Looping       bool

	alive           uint32
	dead            uint32
	emit            uint32
	maxParticles    uint32
	particleBuffer  *buffer
	lifecycleShader *engine.Shader
	simulateShader  *engine.Shader
}

const (
	sizeOfParticle = 80
)

func (m *ModuleCore) SetMaxParticles(value uint32) {
	m.maxParticles = value
	m.dead = m.maxParticles
	m.alive = 0

	m.particleBuffer.SetSize(m.maxParticles)
}

func (m *ModuleCore) syncCounts() {
	m.particleBuffer.Bind()

	gl.BindBuffer(gl.SHADER_STORAGE_BUFFER, m.particleBuffer.idIndex)
	p := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, gl.READ_ONLY)

	out := make([]byte, 12)

	for i := range out {
		out[i] = *((*byte)(unsafe.Pointer(uintptr(p) + uintptr(i))))
	}

	m.alive = binary.LittleEndian.Uint32(out[0:])
	m.dead = binary.LittleEndian.Uint32(out[4:])
	m.emit = binary.LittleEndian.Uint32(out[8:])

	gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	m.particleBuffer.Unbind()
}

func (m *ModuleCore) MaxParticles() uint32 {
	return m.maxParticles
}

func (m *ModuleCore) ParticleCount() uint32 {
	return m.alive
}

func NewModuleCore(maxParticles uint32) *ModuleCore {
	m := &ModuleCore{
		StartColor:    engine.ColorWhite,
		StartSpeed:    10.0,
		StartSize:     1.0,
		PlaybackSpeed: 1.0,
		maxParticles:  maxParticles,
		Looping:       true,
	}

	m.lifecycleShader = engine.GetAsset().MustGet(engine.AssetNameShader, "particle/lifecycle").(*engine.Shader)
	m.simulateShader = engine.GetAsset().MustGet(engine.AssetNameShader, "particle/simulate").(*engine.Shader)

	m.particleBuffer = newBuffer(m.maxParticles)
	m.particleBuffer.Alloc()

	m.SetMaxParticles(m.maxParticles)

	return m
}
