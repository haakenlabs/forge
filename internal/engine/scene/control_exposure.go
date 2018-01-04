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
	"github.com/go-gl/glfw/v3.2/glfw"

	"git.dbnservers.net/haakenlabs/forge/internal/engine"
	"git.dbnservers.net/haakenlabs/forge/internal/engine/scene/effects"
	"git.dbnservers.net/haakenlabs/forge/internal/engine/system/input"
	"git.dbnservers.net/haakenlabs/forge/internal/engine/system/instance"
)

type ControlExposure struct {
	engine.BaseScriptComponent

	tonemapper *effects.Tonemapper
}

func NewControlExposure() *ControlExposure {
	c := &ControlExposure{}

	c.SetName("ControlFly")
	instance.MustAssign(c)

	return c
}

func ControlExposureComponent(e *engine.GameObject) *ControlExposure {
	c := e.Components()
	for i := range c {
		if ct, ok := c[i].(*ControlExposure); ok {
			return ct
		}
	}

	return nil
}

func (c *ControlExposure) SetTonemapper(t *effects.Tonemapper) {
	c.tonemapper = t
}

func (c *ControlExposure) Update() {
	if c.tonemapper == nil {
		return
	}

	if input.KeyDown(glfw.KeyMinus) {
		c.tonemapper.SetExposure(c.tonemapper.Exposure() - 0.05)
	} else if input.KeyDown(glfw.KeyEqual) {
		c.tonemapper.SetExposure(c.tonemapper.Exposure() + 0.05)
	}
}
