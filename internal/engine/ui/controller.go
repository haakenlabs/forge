/*
Copyright (c) 2018 HaakenLabs

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

package ui

import (
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/haakenlabs/forge/internal/engine"
)

type Controller struct {
	engine.BaseScriptComponent

	renderers []Renderer
}

func (c *Controller) UpdateCache() {
	c.renderers = c.renderers[:0]

	components := c.GameObject().ComponentsInChildren()
	for i := range components {
		if r, ok := components[i].(Renderer); ok {
			c.renderers = append(c.renderers, r)
		}
	}
}

func (c *Controller) OnSceneGraphUpdate() {
	c.UpdateCache()
}

func (c *Controller) GUIRender() {
	if len(c.renderers) == 0 {
		return
	}

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE_MINUS_SRC_ALPHA)

	for i := range c.renderers {
		c.renderers[i].UIDraw()
	}

	gl.Disable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
}

func (c *Controller) Resize() {
	if c.GameObject() != nil {
		RectTransformComponent(c.GameObject()).SetSize(engine.GetWindow().Resolution().Vec2())
	}
}

func (c *Controller) Start() {
	c.Resize()
	c.UpdateCache()
}

func (c *Controller) Update() {
	if engine.GetWindow().WindowResized() {
		c.Resize()
	}
}

func NewController() *Controller {
	c := &Controller{}

	c.SetName("UIController")
	engine.GetInstance().MustAssign(c)

	return c
}

func CreateController(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	controller := NewController()

	object.AddComponent(controller)

	return object
}
