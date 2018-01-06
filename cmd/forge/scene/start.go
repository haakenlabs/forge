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

package scene

import (
	"github.com/haakenlabs/forge/internal/engine"
	"github.com/haakenlabs/forge/internal/engine/particle"
	"github.com/haakenlabs/forge/internal/engine/scene"
	"github.com/haakenlabs/forge/internal/engine/scene/effects"
)

const NameStart = "start"

func NewStartScene() *engine.Scene {
	s := engine.NewScene(NameStart)
	s.SetLoadFunc(func() error {
		camera := scene.CreateCamera("camera", true, engine.RenderPathForward)
		camera.AddComponent(scene.NewControlOrbit())
		tonemapper := effects.NewTonemapper()

		cameraC := engine.CameraComponent(camera)
		cameraC.AddEffect(tonemapper)
		cameraC.SetClearMode(engine.ClearModeColor)

		toneControl := scene.NewControlExposure()
		toneControl.SetTonemapper(tonemapper)
		camera.AddComponent(toneControl)

		target := engine.NewGameObject("target")

		//particleSys := engine.NewParticleSystem()
		//
		//particleRen := engine.NewParticleRenderer()
		//
		//particleRen.SetParticleSystem(particleSys)
		//particleRen.SetSprite(image.MustGet("particle.png"))
		//particleRen.SetRenderShader(shader.MustGet("particles/render-particle"))
		//particleRen.SetParticleShader(shader.MustGet("particles/compute-particle"))
		//
		//particleSys.SetStartColor(engine.ColorOrange())
		//particleSys.SetStartLifetime(30.0)
		//
		//particleSys.CreateParticles()
		//
		//target.AddComponent(particleSys)
		//target.AddComponent(particleRen)

		psys := particle.NewParticleSystem(1000000)
		psys.Emission.Rate = 1000
		psys.Core.StartLifetime = 5
		psys.Core.PlaybackSpeed = 1.0

		target.AddComponent(psys)

		scene.ControlOrbitComponent(camera).Target = target.Transform()

		if err := s.Graph().AddGameObject(target, nil); err != nil {
			return err
		}
		if err := s.Graph().AddGameObject(camera, nil); err != nil {
			return err
		}

		return nil
	})

	return s
}
