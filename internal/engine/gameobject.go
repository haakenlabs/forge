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

type Message uint8

const (
	MessageActivate Message = iota
	MessageStart
	MessageAwake
	MessageUpdate
	MessageLateUpdate
	MessageFixedUpdate
	MessageGUIRender
	MessageSGUpdate
)

type GameObject struct {
	BaseObject

	components []Component
	children   []*GameObject
	parent     *GameObject
	scene      *Scene
	active     bool
}

// Active returns the active state of this game object.
func (g *GameObject) Active() bool {
	return g.active
}

// SetActive sets the active state of this game object.
func (g *GameObject) SetActive(active bool) {
	if g.active != active {
		g.active = active

		// TODO: Notify scene graph
		if g.Scene() != nil && g.Scene().Graph() != nil {
			g.Scene().Graph().SetDirty()
		}
	}
}

// Transform returns the transform for this game object.
func (g *GameObject) Transform() Transform {
	return g.components[0].(Transform)
}

// SetTransform sets the transform for this game object.
func (g *GameObject) SetTransform(transform Transform) {
	if transform == nil {
		return
	}
	if g.Transform().ID() == transform.ID() {
		return
	}

	g.components[0] = transform

	// TODO: Notify transform of transition
}

// SendMessage calls the function associated with the given message.
func (g *GameObject) SendMessage(msg Message) {
	if !g.active {
		return
	}

	if msg == MessageActivate {
		g.activate()
		g.SendMessage(MessageAwake)
		return
	}

	for i := range g.components {
		switch msg {
		case MessageStart:
			if c, ok := g.components[i].(ScriptComponent); ok {
				c.Start()
			}
		case MessageAwake:
			if c, ok := g.components[i].(ScriptComponent); ok {
				c.Awake()
			}
		case MessageUpdate:
			if c, ok := g.components[i].(ScriptComponent); ok {
				c.Update()
			}
		case MessageLateUpdate:
			if c, ok := g.components[i].(ScriptComponent); ok {
				c.LateUpdate()
			}
		case MessageFixedUpdate:
			if c, ok := g.components[i].(ScriptComponent); ok {
				c.FixedUpdate()
			}
		case MessageGUIRender:
			if c, ok := g.components[i].(ScriptComponent); ok {
				c.GUIRender()
			}
		case MessageSGUpdate:
			if c, ok := g.components[i].(SceneGraphListener); ok {
				c.OnSceneGraphUpdate()
			}
		}
	}
}

// SetParent sets the parent of this game object.
func (g *GameObject) SetParent(object *GameObject) {
	if object == nil {
		return
	}

	g.parent = object
}

// Parent returns the parent of this game object.
func (g *GameObject) Parent() *GameObject {
	return g.parent
}

// Children returns the child game objects of this game object.
func (g *GameObject) Children() []*GameObject {
	return g.children
}

// Components returns the components attached to this game object.
func (g *GameObject) Components() []Component {
	return g.components
}

// ComponentsInChildren returns the components in any of the child objects.
func (g *GameObject) ComponentsInChildren() []Component {
	var components []Component

	if len(g.children) == 0 {
		return components
	}

	if g.Scene() == nil {
		for i := range g.children {
			components = append(components, g.children[i].Components()...)
			components = append(components, g.children[i].ComponentsInChildren()...)
		}
		return components
	}

	sg := g.Scene().Graph()
	descendants := sg.Descendants(g)
	for i := range descendants {
		components = append(components, descendants[i].Components()...)
	}

	return components
}

// ComponentsInChildren returns the components in any of the parent objects.
func (g *GameObject) ComponentsInParent() []Component {
	ancestors := g.Ancestors()
	var components []Component

	if len(ancestors) == 0 {
		return components
	}

	for i := range ancestors {
		components = append(components, ancestors[i].Components()...)
	}

	return components
}

// AddComponent attaches a component to this game object.
func (g *GameObject) AddComponent(component Component) {
	if component == nil {
		return
	}

	g.components = append(g.components, component)
	component.SetGameObject(g)
}

// AddChild adds a child game object to this game object.
func (g *GameObject) AddChild(child *GameObject) {
	for i := range g.children {
		if g.children[i].ID() == child.ID() {
			return
		}
	}

	g.children = append(g.children, child)
}

// RemoveChild removes a child game object from this game object by ID.
func (g *GameObject) RemoveChild(id uint32) {
	// TODO: break all underlying associations
	for i := range g.children {
		if g.children[i].ID() == id {
			g.children[i] = g.children[len(g.children)-1]
			g.children = g.children[:len(g.children)-1]
		}
	}
}

// AddChild removes all child objects from this game object.
func (g *GameObject) RemoveAllChildren() {
	// TODO: Add me
	panic(ErrNotImplemented)
}

// activate is called when the game object is initialized and needs to build
// associations between itself and other objects and components. This typically
// occurs when an object is added to a scene graph.
func (g *GameObject) activate() {
	// Update component references.
	for i := range g.components {
		g.components[i].SetGameObject(g)
	}
}

// SetScene sets the scene for this game object. Once set, the scene cannot be
// changed or unset.
func (g *GameObject) SetScene(scene *Scene) {
	if g.scene == nil && scene != nil {
		g.scene = scene
	}
}

// Scene returns the scene for this game object.
func (g *GameObject) Scene() *Scene {
	return g.scene
}

// Ancestors lists all ancestor objects of this game object.
func (g *GameObject) Ancestors() []*GameObject {
	var ancestors []*GameObject

	if g.Parent() != nil {
		ancestors = append(ancestors, g.Parent())
		ancestors = append(ancestors, g.Parent().Ancestors()...)
	}

	return ancestors
}

// OnParentChanged is called when the parent of this gameobject has changed.
func (g *GameObject) OnParentChanged() {
	g.Transform().Recompute(true)
}

// NewGameObject creates a new GameObject.
func NewGameObject(name string) *GameObject {
	g := &GameObject{
		active:     true,
		components: make([]Component, 1),
		children:   []*GameObject{},
	}

	g.SetName(name)
	GetInstance().MustAssign(g)

	g.components[0] = NewTransform()

	return g
}
