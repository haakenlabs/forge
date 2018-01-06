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
	//"fmt"
	"math"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"

	"fmt"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/haakenlabs/forge/internal/engine"
	"github.com/haakenlabs/forge/internal/engine/system/input"
	"github.com/haakenlabs/forge/internal/engine/system/instance"
)

var _ engine.Renderer = &System{}

const (
	workgroupSize = uint32(128)
)

type System struct {
	engine.BaseScriptComponent

	Core     *ModuleCore
	Emission *ModuleEmission
	Force    *ModuleForce
	Noise    *ModuleNoise
	Renderer *ModuleRenderer
	Shape    *ModuleShape
	Velocity *ModuleVelocity

	timeStart  float64
	inOffset   uint32
	outOffset  uint32
	playing    bool
	paused     bool
	bufferFlip bool
}

func (s *System) SupportsDeferred() bool {
	return false
}

func (s *System) RenderShader(shader *engine.Shader, camera *engine.Camera) {
	s.Render(camera)
	shader.Bind()
}

func (s *System) Simulate() {

	deltaTime := float32(engine.GetTime().DeltaTime()) * s.Core.PlaybackSpeed

	// Lifecycle Phase
	s.Core.lifecycleShader.Bind()
	s.Core.particleBuffer.Bind()

	s.Core.lifecycleShader.SetUniform("u_offset_in", s.inOffset)
	s.Core.lifecycleShader.SetUniform("u_offset_out", s.outOffset)
	s.Core.lifecycleShader.SetUniform("u_max_particles", s.Core.maxParticles)
	s.Core.lifecycleShader.SetUniform("u_delta_time", deltaTime)
	s.Core.lifecycleShader.SetUniform("u_start_color", s.Core.StartColor.Vec4())
	s.Core.lifecycleShader.SetUniform("u_angular_velocity_3d", mgl32.Vec3{})
	s.Core.lifecycleShader.SetUniform("u_rotation", mgl32.Vec3{})
	s.Core.lifecycleShader.SetUniform("u_position", mgl32.Vec3{})
	s.Core.lifecycleShader.SetUniform("u_random_seed", uint32(engine.GetTime().Frame())^s.Core.RandomSeed)
	s.Core.lifecycleShader.SetUniform("u_angular_velocity", 0.0)
	s.Core.lifecycleShader.SetUniform("u_start_lifetime", s.Core.StartLifetime)
	s.Core.lifecycleShader.SetUniform("u_start_size", s.Core.StartSize)
	s.Core.lifecycleShader.SetUniform("u_velocity", s.Core.StartSpeed)

	s.Core.particleBuffer.ResetIndex(s.Core.alive, s.Core.dead)

	if s.Core.alive > 0 {
		//fmt.Println("task_lifetime")
		// Compute lifetime
		s.Core.lifecycleShader.SetSubroutine(engine.ShaderComponentCompute, "task_lifetime")
		s.Core.lifecycleShader.SetUniform("u_invocations", s.Core.alive)

		gl.DispatchCompute((s.Core.alive/workgroupSize)+1, 1, 1)
		gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
	}

	// Emit new particles
	emitNow := uint32(math.Ceil(float64(s.Emission.Rate) * float64(deltaTime)))

	if emitNow > 0 && s.Core.dead > 0 {
		if emitNow > s.Core.dead {
			emitNow = s.Core.dead
		}

		s.Core.lifecycleShader.SetSubroutine(engine.ShaderComponentCompute, "task_emit")
		s.Core.lifecycleShader.SetUniform("u_invocations", emitNow)

		gl.DispatchCompute((emitNow/workgroupSize)+1, 1, 1)
		gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
	}

	s.Core.syncCounts()

	s.Core.lifecycleShader.Unbind()
	s.Core.simulateShader.Bind()

	// Simulate phase
	if s.Core.alive > 0 {
		s.Core.simulateShader.SetUniform("u_invocations", s.Core.alive)
		s.Core.simulateShader.SetUniform("u_offset_out", s.outOffset)
		s.Core.simulateShader.SetUniform("u_delta_time", deltaTime)
		s.Core.simulateShader.SetUniform("u_attractors", s.Force.EnableAttractors)

		gl.DispatchCompute((s.Core.alive/workgroupSize)+1, 1, 1)
		gl.MemoryBarrier(gl.SHADER_STORAGE_BARRIER_BIT)
	}

	s.Core.simulateShader.Unbind()

	s.swapBuffers()
}

func (s *System) Update() {
	if input.KeyDown(glfw.Key9) {
		s.Emission.Rate += 500.0
	} else if input.KeyDown(glfw.Key8) {
		s.Emission.Rate -= 500.0
		if s.Emission.Rate < 0 {
			s.Emission.Rate = 0
		}
	} else if input.KeyDown(glfw.Key0) {
		s.Emission.Rate = 0.0
	} else if input.KeyDown(glfw.KeySpace) {
		s.Emission.Rate = 10000.0
	}

	if input.KeyUp(glfw.KeySpace) {
		s.Emission.Rate = 0.0
	}

	if input.KeyDown(glfw.Key6) {
		s.Core.StartLifetime -= 1.0

		if s.Core.StartLifetime <= 0 {
			s.Core.StartLifetime = 1
		}
	} else if input.KeyDown(glfw.Key7) {
		s.Core.StartLifetime += 1.0
	}

	if input.KeyDown(glfw.Key5) {
		s.Force.EnableAttractors = !s.Force.EnableAttractors
		fmt.Printf("EnableAttractors: %v\n", s.Force.EnableAttractors)
	}

	if input.KeyDown(glfw.Key4) {
		s.Core.PlaybackSpeed += 0.1
	} else if input.KeyDown(glfw.Key3) {
		s.Core.PlaybackSpeed -= 0.1
		if s.Core.PlaybackSpeed < 0 {
			s.Core.PlaybackSpeed = 0
		}
	}
}

func (s *System) Render(camera *engine.Camera) {
	if s.Renderer != nil {
		s.Renderer.Render(camera)
	}
}

func (s *System) Stop() {
	s.playing = false
	s.paused = false
}

func (s *System) Play() {
	s.playing = true
	s.paused = false
}

func (s *System) Pause() {
	s.playing = false
	s.paused = true
}

func (s *System) swapBuffers() {
	s.bufferFlip = !s.bufferFlip

	s.inOffset = uint32(0)
	s.outOffset = s.Core.maxParticles

	if s.bufferFlip {
		s.inOffset = s.Core.maxParticles
		s.outOffset = 0
	}
}

func NewParticleSystem(maxParticles uint32) *System {
	s := &System{
		inOffset:  0,
		outOffset: maxParticles,
	}

	// Make the modules
	s.Core = NewModuleCore(maxParticles)
	s.Renderer = NewModuleRenderer(s)
	s.Emission = NewModuleEmission()
	s.Shape = NewModuleShape()
	s.Velocity = NewModuleVelocity()
	s.Noise = NewModuleNoise()
	s.Force = NewModuleForce()

	s.SetName("ParticleSystem")
	instance.MustAssign(s)

	s.Alloc()

	return s
}
