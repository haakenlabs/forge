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

type Component interface {
	Object

	// GameObject returns the GameObject for this component.
	GameObject() *GameObject

	// SetGameObject sets the GameObject for this component.
	SetGameObject(*GameObject)

	// GetTransform returns the Transform for this component.
	GetTransform() Transform

	// Validate the state of the component. By default, this does nothing.
	Validate() error
}

type ScriptComponent interface {
	Component

	// Active returns the active state of this component.
	Active() bool

	// SetActive sets the active state of this component.
	SetActive(bool)

	// OnActivate is called when the component transitions to the active state.
	OnActivate()

	// OnDeactivate is called when the component transitions to the inactive state.
	OnDeactivate()

	// Awake is called when the Entity is loaded by the scene graph. Note: it is
	// not guaranteed that all parent and child associations are made. If such conditions
	// are required, use Start() instead.
	Awake()

	// FixedUpdate is called at a fixed interval, but not every frame.
	FixedUpdate()

	// LateUpdate is called after all Update calls.
	LateUpdate()

	// Update is called every frame, before the Render call.
	Update()

	// Start is called when the component is first initialized. This occurs when the
	// Entity is added to the scene graph and after associations are made between the
	// component and Entity.
	Start()

	// GUIRender is called during the GUI drawing phase of the rendering.
	GUIRender()
}

var _ Component = &BaseComponent{}
var _ ScriptComponent = &BaseScriptComponent{}

type BaseComponent struct {
	BaseObject

	gameobject *GameObject
}

type BaseScriptComponent struct {
	BaseComponent

	active bool
}

// GameObject returns the GameObject for this component.
func (c *BaseComponent) GameObject() *GameObject {
	return c.gameobject
}

// GetTransform returns the Transform for this component.
func (c *BaseComponent) GetTransform() Transform {
	if c.gameobject != nil {
		return c.gameobject.Transform()
	}

	return nil
}

// SetGameObject sets the GameObject for this component.
func (c *BaseComponent) SetGameObject(gameobject *GameObject) {
	c.gameobject = gameobject
}

// Validate the state of the component. By default, this does nothing.
func (c *BaseComponent) Validate() error {
	return nil
}

// Active returns the active state of this component.
func (c *BaseScriptComponent) Active() bool {
	return c.active
}

// SetActive sets the active state of this component.
func (c *BaseScriptComponent) SetActive(active bool) {
	if c.active != active {
		c.active = active

		if c.active {
			c.OnActivate()
		} else {
			c.OnDeactivate()
		}
	}
}

// OnActivate is called when the component transitions to the active state.
func (c *BaseScriptComponent) OnActivate() {}

// OnDeactivate is called when the component transitions to the inactive state.
func (c *BaseScriptComponent) OnDeactivate() {}

// Awake is called when the Entity is loaded by the scene graph. Note: it is
// not guaranteed that all parent and child associations are made. If such conditions
// are required, use Start() instead.
func (c *BaseScriptComponent) Awake() {}

// FixedUpdate is called at a regular interval, but not every frame.
func (c *BaseScriptComponent) FixedUpdate() {}

// LateUpdate is called after all Update calls.
func (c *BaseScriptComponent) LateUpdate() {}

// Update is called every frame, before the Render call.
func (c *BaseScriptComponent) Update() {}

// Start is called when the component is first initialized. This occurs when the
// Entity is added to the scene graph and after associations are made between the
// component and Entity.
func (c *BaseScriptComponent) Start() {}

// GUIRender is called during the GUI drawing phase of the rendering.
func (c *BaseScriptComponent) GUIRender() {}
