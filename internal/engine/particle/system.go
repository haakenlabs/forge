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
	"git.dbnservers.net/haakenlabs/forge/internal/engine"
	"git.dbnservers.net/haakenlabs/forge/internal/engine/system/instance"
)

var _ engine.Renderer = &System{}

type System struct {
	engine.BaseScriptComponent

	Core     *ModuleCore
	Emission *ModuleEmission
	Force    *ModuleForce
	Noise    *ModuleNoise
	Renderer *ModuleRenderer
	Shape    *ModuleShape
	Velocity *ModuleVelocity

	pbuffer   *buffer
	timeStart float64
	playing   bool
	paused    bool
}

func (s *System) SupportsDeferred() bool {
	return false
}

func (s *System) RenderShader(shader *engine.Shader, camera *engine.Camera) {
	s.Render(camera)
	shader.Bind()
}

func (s *System) Update() {

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

func NewParticleSystem() *System {
	s := &System{}

	// Make the modules
	s.Core = NewModuleCore()
	s.Renderer = NewModuleRenderer()
	s.Emission = NewModuleEmission()
	s.Shape = NewModuleShape()
	s.Velocity = NewModuleVelocity()
	s.Noise = NewModuleNoise()
	s.Force = NewModuleForce()

	s.SetName("ParticleSystem")
	instance.MustAssign(s)

	s.pbuffer = newBuffer(s.Core.MaxParticles)
	s.Alloc()

	return s
}
