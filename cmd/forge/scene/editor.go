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
	"github.com/haakenlabs/forge/internal/engine/scene"
	"github.com/haakenlabs/forge/internal/engine/scene/effects"
)

const NameEditor = "editor"

func NewEditorScene() *engine.Scene {
	s := engine.NewScene(NameEditor)
	s.SetLoadFunc(func() error {
		testObject := engine.NewGameObject("testObject")
		camera := scene.CreateCamera("camera", true, engine.RenderPathDeferred)
		camera.AddComponent(scene.NewControlOrbit())
		tonemapper := effects.NewTonemapper()

		cameraC := engine.CameraComponent(camera)
		cameraC.AddEffect(tonemapper)

		toneControl := scene.NewControlExposure()
		toneControl.SetTonemapper(tonemapper)
		camera.AddComponent(toneControl)

		test := scene.CreateOrb("orb")

		scene.ControlOrbitComponent(camera).Target = test.Transform()

		if err := s.Graph().AddGameObject(testObject, nil); err != nil {
			return err
		}
		if err := s.Graph().AddGameObject(camera, nil); err != nil {
			return err
		}
		if err := s.Graph().AddGameObject(test, nil); err != nil {
			return err
		}

		return nil
	})

	return s
}
