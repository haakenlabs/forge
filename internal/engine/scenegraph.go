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
	"github.com/sirupsen/logrus"

	"git.dbnservers.net/haakenlabs/forge/internal/sg"
)

type SceneGraphListener interface {
	// OnSceneGraphUpdate is called when the SceneGraph has been updated.
	OnSceneGraphUpdate()
}

type SceneGraph struct {
	root           *GameObject
	graph          *sg.Graph
	active         []*GameObject
	componentCache []Component
	scene          *Scene
	dirty          bool
}

func NewSceneGraph(scene *Scene) *SceneGraph {
	s := &SceneGraph{
		graph:          sg.NewGraph(),
		scene:          scene,
		dirty:          true,
		active:         []*GameObject{},
		componentCache: []Component{},
	}

	s.root = NewGameObject("__rootNode__")

	return s
}

func (s *SceneGraph) Update() {
	dfs := s.graph.DepthFirstSearch(1, false)
	s.active = s.active[:0]
	s.componentCache = s.componentCache[:0]

	for _, v := range dfs {
		o := s.graph.GetObjectAtVertex(v).(*GameObject)
		s.active = append(s.active, o)
		s.componentCache = append(s.componentCache, o.Components()...)
	}

	s.dirty = false
	s.notifyListeners()

	logrus.Debugf("SceneGraph updated. activeObjects: %d componentCache: %d", len(s.active), len(s.componentCache))
}

// Dirty returns the state of the graph. If true, the graph needs an update.
func (s *SceneGraph) Dirty() bool {
	return s.dirty
}

// SetDirty sets a flag which indicates the graph needs to be updated.
func (s *SceneGraph) SetDirty() {
	s.dirty = true
}

func (s *SceneGraph) AddGameObject(object, parent *GameObject) error {
	var err error
	var u sg.VertexDescriptor // parent
	var v sg.VertexDescriptor // object

	// Get parent VertexDescriptor.
	if parent != nil {
		u, err = s.graph.GetVertexByObject(parent)
	} else {
		u, err = s.graph.GetVertexByObject(s.root)
	}
	if err != nil {
		return err
	}

	// Get object VertexDescriptor.
	v, err = s.graph.GetVertexByObject(object)
	if err != nil {
		return err
	}
	// Add the edge.
	if err := s.graph.AddEdge(u, v); err != nil {
		return err
	}

	// Add children, if any.
	children := object.Children()
	for i := range children {
		if err := s.AddGameObject(children[i], object); err != nil {
			return err
		}
	}

	// Build associations.
	object.SetScene(s.scene)
	s.updateReferences(object)
	object.SendMessage(MessageActivate)

	// Notify graph of update.
	s.Update()

	return nil
}

func (s *SceneGraph) SendMessage(message Message) {
	for i := range s.active {
		s.active[i].SendMessage(message)
	}
}

func (s *SceneGraph) Parent(e *GameObject) *GameObject {
	v, err := s.graph.GetVertexByObject(e)
	if err != nil {
		return nil
	}

	u, err := s.graph.Parent(v)
	if err != nil {
		return nil
	}

	return s.objectAt(u)
}

// Children returns the direct descendants of the requested entity.
func (s *SceneGraph) Children(e *GameObject) []*GameObject {
	children := []*GameObject{}

	objDesc, err := s.graph.GetVertexByObject(e)
	if err != nil {
		return children
	}

	for _, v := range s.graph.ChildrenOf(objDesc) {
		children = append(children, s.objectAt(v))
	}

	return children
}

// Descendants returns the descendant entities of requested entity.
func (s *SceneGraph) Descendants(e *GameObject) []*GameObject {
	d := []*GameObject{}

	u, err := s.graph.GetVertexByObject(e)
	if err != nil {
		logrus.Error(err)
		return d
	}

	dfs := s.graph.DepthFirstSearch(u, true)

	for _, v := range dfs {
		obj := s.objectAt(v)
		if obj.ID() == e.ID() {
			continue
		}
		d = append(d, obj)
	}

	return d
}

// Components returns all active components in the SceneGraph.
func (s *SceneGraph) Components() []Component {
	return s.componentCache
}

func (s *SceneGraph) objectAt(u sg.VertexDescriptor) *GameObject {
	obj := s.graph.GetObjectAtVertex(u)
	if obj == nil {
		return nil
	}

	return obj.(*GameObject)
}

func (s *SceneGraph) notifyListeners() {
	s.scene.OnSceneGraphUpdate()

	s.SendMessage(MessageSGUpdate)
}

func (s *SceneGraph) updateReferences(g *GameObject) {
	parent := s.Parent(g)
	if parent != nil {
		parent.AddChild(g)
		g.SetParent(parent)
	}

	childObjects := s.Children(g)

	for idx := range childObjects {
		childObjects[idx].SetParent(g)
	}
}
