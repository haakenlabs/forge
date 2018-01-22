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
	"math"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

type ParticleSystem struct {
	BaseScriptComponent

	buffer          *ParticleBuffer
	genShader       *Shader
	startColor      Color
	startRotation3D mgl32.Vec3
	particleCount   uint32
	maxParticles    uint32
	duration        float32
	playbackSpeed   float32
	startDelay      float32
	startLifetime   float32
	startRotation   float32
	startSpeed      float32
	startSize       float32
	time            float32
	looping         bool
	paused          bool
	playing         bool
	autoplay        bool
}

var _ Renderer = &ParticleRenderer{}

type ParticleRenderer struct {
	BaseComponent

	particleShader *Shader
	renderShader   *Shader
	system         *ParticleSystem
	sprite         *Texture2D
	clock          float64
}

type Particle struct {
	StartColor      Color
	AngularVelocity mgl32.Vec4
	Rotation        mgl32.Vec4
	Position        mgl32.Vec4
	Lifecycle       mgl32.Vec4
}

type ParticleAttractor struct {
	Position mgl32.Vec3
	Force    float32
}

type ParticleBuffer struct {
	BaseObject

	vao  uint32
	ssbo uint32
	size uint32
}

const (
	sizeOfParticle = 80
)

func (p *ParticleSystem) Clear() {
	p.buffer.Upload(nil)
}

func (p *ParticleSystem) CreateParticles() {
	if p.maxParticles <= 0 || p.maxParticles > math.MaxUint32 {
		logrus.Errorf("ParticleSystem has invalid maxParticle count: %d", p.maxParticles)
		return
	}

	p.genShader.Bind()
	p.buffer.Bind()

	p.genShader.SetUniform("c_max_particles", p.maxParticles)
	p.genShader.SetUniform("c_start_color", p.startColor.Vec4())
	p.genShader.SetUniform("c_angular_velocity_3d", mgl32.Vec3{})
	p.genShader.SetUniform("c_rotation", mgl32.Vec3{})
	p.genShader.SetUniform("c_position", mgl32.Vec3{})
	p.genShader.SetUniform("c_random_seed", 0)
	p.genShader.SetUniform("c_angular_velocity", 0.0)
	p.genShader.SetUniform("c_lifetime", p.startLifetime)
	p.genShader.SetUniform("c_start_lifetime", p.startLifetime)
	p.genShader.SetUniform("c_start_size", p.startSize)
	p.genShader.SetUniform("c_velocity", 0.0)

	dispatchCount := uint32(100)

	if p.maxParticles > 50000 {
		dispatchCount = 1000
	}

	gl.DispatchCompute(dispatchCount, 1, 1)
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)

	p.buffer.Unbind()
	p.genShader.Unbind()

	p.particleCount = p.maxParticles

	logrus.Debugf("%s created %d particles", p, p.particleCount)
}

func (p *ParticleSystem) Emit(count uint32) {

}

func (p *ParticleSystem) Pause() {
	p.paused = true
}

func (p *ParticleSystem) Play() {
	p.paused = false
	p.playing = true
}

func (p *ParticleSystem) Stop() {
	p.paused = false
	p.playing = false
}

func (p *ParticleSystem) GetData() []Particle {
	return p.buffer.GetData()
}

func (p *ParticleSystem) ParticleCount() uint32 {
	return p.particleCount
}

func (p *ParticleSystem) MaxParticles() uint32 {
	return p.maxParticles
}

func (p *ParticleSystem) StartColor() Color {
	return p.startColor
}

func (p *ParticleSystem) StartRotation3D() mgl32.Vec3 {
	return p.startRotation3D
}

func (p *ParticleSystem) Duration() float32 {
	return p.duration
}

func (p *ParticleSystem) PlaybackSpeed() float32 {
	return p.playbackSpeed
}

func (p *ParticleSystem) StartDelay() float32 {
	return p.startDelay
}

func (p *ParticleSystem) StartLifetime() float32 {
	return p.startLifetime
}

func (p *ParticleSystem) StartRotation() float32 {
	return p.startRotation
}

func (p *ParticleSystem) StartSize() float32 {
	return p.startSize
}

func (p *ParticleSystem) Time() float32 {
	return p.time
}

func (p *ParticleSystem) Autoplay() bool {
	return p.autoplay
}

func (p *ParticleSystem) Looping() bool {
	return p.looping
}

func (p *ParticleSystem) Paused() bool {
	return p.paused
}

func (p *ParticleSystem) Playing() bool {
	return p.playing
}

func (p *ParticleSystem) SetMaxParticles(value uint32) {
	p.maxParticles = value
}

func (p *ParticleSystem) SetStartColor(value Color) {
	p.startColor = value
}

func (p *ParticleSystem) SetDuration(value float32) {
	p.duration = value
}

func (p *ParticleSystem) SetPlaybackSpeed(value float32) {
	p.playbackSpeed = value
}

func (p *ParticleSystem) SetStartDelay(value float32) {
	p.startDelay = value
}

func (p *ParticleSystem) SetStartLifetime(value float32) {
	p.startLifetime = value
}

func (p *ParticleSystem) SetStartRotation(value float32) {
	p.startRotation = value
}

func (p *ParticleSystem) SetStartSize(value float32) {
	p.startSize = value
}

func (p *ParticleSystem) SetTime(value float32) {
	p.time = value
}

func (p *ParticleSystem) SetLooping(value bool) {
	p.looping = value
}

func (p *ParticleSystem) SetAutoplay(value bool) {
	p.autoplay = value
}

func (p *ParticleRenderer) SetParticleShader(shader *Shader) {
	p.particleShader = shader
}

func (p *ParticleRenderer) SetRenderShader(shader *Shader) {
	p.renderShader = shader
}

func (p *ParticleRenderer) SetParticleSystem(system *ParticleSystem) {
	p.system = system
}

func (p *ParticleRenderer) SetSprite(sprite *Texture2D) {
	p.sprite = sprite
}

func (p *ParticleRenderer) SupportsDeferred() bool {
	return false
}

func (p *ParticleRenderer) RenderShader(shader *Shader, camera *Camera) {
	p.Render(camera)
	shader.Bind()
}

func (p *ParticleRenderer) Render(camera *Camera) {
	if p.system == nil || p.renderShader == nil || p.particleShader == nil || p.sprite == nil {
		return
	}

	p.particleShader.Bind()
	p.system.buffer.Bind()
	p.clock += GetTime().DeltaTime()
	fmt.Println(p.clock)
	p.particleShader.SetUniform("c_delta_time", float32(GetTime().DeltaTime()))
	p.particleShader.SetUniform("c_max_particles", p.system.MaxParticles())

	dispatchCount := uint32(100)

	if p.system.MaxParticles() > 50000 {
		dispatchCount = 1000
	}

	gl.DispatchCompute(dispatchCount, 1, 1)
	gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)

	p.renderShader.Bind()
	p.system.buffer.Bind()

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	p.renderShader.SetUniform("v_model_matrix", p.GetTransform().ActiveMatrix())
	p.renderShader.SetUniform("v_view_matrix", camera.ViewMatrix())
	p.renderShader.SetUniform("v_projection_matrix", camera.ProjectionMatrix())

	gl.PointSize(1.0)
	p.sprite.ActivateTexture(gl.TEXTURE0)
	p.renderShader.SetUniform("f_time", GetTime().Now())
	gl.DrawArrays(gl.POINTS, 0, int32(p.system.ParticleCount()))

	p.renderShader.Unbind()

	gl.Disable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
}

func (p *ParticleBuffer) Bind() {
	gl.BindVertexArray(p.vao)
}

func (p *ParticleBuffer) Unbind() {
	gl.BindVertexArray(0)
}

func (p *ParticleBuffer) GetData() []Particle {
	if p.size <= 0 {
		return nil
	}

	data := make([]Particle, p.size)

	p.Bind()
	ptr := gl.MapBuffer(gl.SHADER_STORAGE_BUFFER, gl.READ_ONLY)
	ptrArr := *((*[]Particle)(ptr))

	copy(data, ptrArr)

	gl.UnmapBuffer(gl.SHADER_STORAGE_BUFFER)
	p.Unbind()

	return data
}

func (p *ParticleBuffer) Size() uint32 {
	return p.size
}

func (p *ParticleBuffer) Upload(data []Particle) {
	p.size = uint32(len(data))

	p.Bind()
	gl.BufferData(gl.SHADER_STORAGE_BUFFER, int(p.size)*sizeOfParticle, gl.Ptr(data), gl.DYNAMIC_DRAW)
	p.Unbind()
}

func (p *ParticleBuffer) Alloc() error {
	gl.GenVertexArrays(1, &p.vao)
	gl.GenBuffers(1, &p.ssbo)

	gl.BindVertexArray(p.vao)
	gl.BindBufferBase(gl.SHADER_STORAGE_BUFFER, 0, p.ssbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, p.ssbo)
	gl.BufferData(gl.ARRAY_BUFFER, int(p.size)*sizeOfParticle, nil, gl.DYNAMIC_DRAW)

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

	logrus.Debugf("%s allocated particle buffer", p)

	return nil
}

func (p *ParticleBuffer) Dealloc() {
	gl.DeleteBuffers(1, &p.ssbo)
	gl.DeleteVertexArrays(1, &p.vao)
}

func NewParticleSystem() *ParticleSystem {
	p := &ParticleSystem{
		maxParticles:  1000000,
		startColor:    ColorWhite,
		playbackSpeed: 1.0,
		startSpeed:    10.0,
		startSize:     1.0,
		looping:       true,
	}

	p.SetName("ParticleSystem")
	GetInstance().MustAssign(p)

	p.buffer = NewParticleBuffer(p.maxParticles)
	p.buffer.Alloc()

	p.genShader = GetAsset().MustGet(AssetNameShader, "particles/create-particle").(*Shader)

	return p
}

func NewParticleRenderer() *ParticleRenderer {
	p := &ParticleRenderer{}

	p.SetName("ParticleRenderer")
	GetInstance().MustAssign(p)

	return p
}

func NewParticleBuffer(size uint32) *ParticleBuffer {
	p := &ParticleBuffer{
		size: size,
	}

	p.SetName("ParticleBuffer")
	GetInstance().MustAssign(p)

	return p
}
