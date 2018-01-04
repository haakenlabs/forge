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

type Scene struct {
	environment      *Environment
	graph            *SceneGraph
	cameras          []*Camera
	loadFunc         func() error
	onActivateFunc   func()
	onDeactivateFunc func()
	name             string
	loaded           bool
}

func (s *Scene) Setup() {
	s.graph = NewSceneGraph(s)
}

// Name returns the name of this scene.
func (s *Scene) Name() string {
	return s.name
}

// OnActivate is called when the scene transitions to the active state.
func (s *Scene) OnActivate() {
	if s.onActivateFunc != nil {
		s.onActivateFunc()
	}
}

// OnDeactivate is called when the scene transitions to the inactive state.
func (s *Scene) OnDeactivate() {
	if s.onDeactivateFunc != nil {
		s.onDeactivateFunc()
	}
}

// Load is called when the scene is being initialized.
func (s *Scene) Load() error {
	if s.loaded {
		return nil
	}

	s.graph = NewSceneGraph(s)
	s.environment = NewEnvironment()

	if s.loadFunc != nil {
		s.loadFunc()
	}

	s.loaded = true

	return nil
}

// Loaded reports if the scene has been loaded.
func (s *Scene) Loaded() bool {
	return s.loaded
}

// Graph gets the SceneGraph for this Scene.
func (s *Scene) Graph() *SceneGraph {
	return s.graph
}

func (s *Scene) OnSceneGraphUpdate() {
	s.cameras = s.cameras[:0]

	// Update renderer cache.
	components := s.graph.Components()
	for i := range components {
		if c, ok := components[i].(*Camera); ok {
			s.cameras = append(s.cameras, c)
		}
	}
}

// Graph gets the SceneGraph for this Scene.
func (s *Scene) Environment() *Environment {
	return s.environment
}

func (s *Scene) SetLoadFunc(fn func() error) {
	s.loadFunc = fn
}

func (s *Scene) SetOnActivateFunc(fn func()) {
	s.graph.notifyListeners()

	s.onActivateFunc = fn
}

func (s *Scene) SetOnDeactivateFunc(fn func()) {
	s.onDeactivateFunc = fn
}

func NewScene(name string) *Scene {
	s := &Scene{
		name:    name,
		cameras: []*Camera{},
	}

	return s
}
