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
	"math"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/haakenlabs/forge/internal/engine"
	"github.com/haakenlabs/forge/internal/engine/system/instance"
	"github.com/haakenlabs/forge/internal/engine/system/input"
	"github.com/haakenlabs/forge/internal/engine/system/time"

	forgemath "github.com/haakenlabs/forge/internal/math"
)

type ControlOrbit struct {
	engine.BaseScriptComponent

	Target engine.Transform

	radial  float64 // Radial desired amount
	radialL float64 // Radial lerp factor
	radialC float64 // Radial current amount

	phi  float64
	phiL float64
	phiC float64

	theta  float64
	thetaL float64
	thetaC float64

	mouseDrag  bool
	mouseDown  bool
	mouseStart mgl32.Vec2
	mouseLast  mgl32.Vec2
	mouseDelta mgl32.Vec2
}

func NewControlOrbit() *ControlOrbit {
	c := &ControlOrbit{
		radial:  4,
		radialL: 0.1,
		phi:     math.Pi / 2.0,
		phiL:    0.1,
		thetaL:  0.1,
	}

	c.SetName("ControlOrbit")
	instance.MustAssign(c)

	c.radialC = c.radial
	c.phiC = c.phi
	c.thetaC = c.theta

	return c
}

func ControlOrbitComponent(g *engine.GameObject) *ControlOrbit {
	c := g.Components()
	for i := range c {
		if ct, ok := c[i].(*ControlOrbit); ok {
			return ct
		}
	}

	return nil
}

func (c *ControlOrbit) move() {
	if c.Target == nil {
		return
	}

	c.GetTransform().SetPosition(
		sphericalToCartesian(c.radialC, c.thetaC, c.phiC))
	c.GetTransform().SetRotation(
		lookAt(c.GetTransform().Position(), c.Target.Position(), mgl32.Vec3{0, 1, 0}))
	engine.CameraComponent(c.GameObject()).SetViewMatrix(
		mgl32.LookAtV(c.GetTransform().Position(), c.Target.Position(), mgl32.Vec3{0, 1, 0}))
}

func (c *ControlOrbit) Awake() {
	c.move()
}

func (c *ControlOrbit) LateUpdate() {
	var changed bool
	var dX float64
	var dY float64

	if c.Target == nil {
		return
	}

	if input.KeyDown(glfw.KeyR) {
		c.radial = 4
		c.phi = math.Pi / 2.0
		c.theta = 0

		c.move()
		return
	}

	if input.MouseWheel() {
		c.radial -= input.MouseWheelY() * .25
		if c.radial < 0.1 {
			c.radial = 0.1
		}
	}

	if input.MouseDown(glfw.MouseButtonLeft) {
		c.mouseDown = true
	}

	if input.MouseMoved() && c.mouseDown {
		if !c.mouseDrag {
			c.mouseDrag = true
			c.mouseStart = input.MousePosition()
			c.mouseLast = input.MousePosition()
		}
	}

	if input.MouseUp(glfw.MouseButtonLeft) {
		c.mouseDown = false
		c.mouseDrag = false
	}

	// Dragging, do movement stuff here.
	if c.mouseDrag && input.MouseMoved() {
		c.mouseDelta = input.MousePosition().Sub(c.mouseLast)
		c.mouseLast = input.MousePosition()

		dX = float64(c.mouseDelta.X()*.25) * time.Delta()
		dY = float64(-c.mouseDelta.Y()*.25) * time.Delta()

		tmpPhi := c.phi

		c.theta -= dX
		tmpPhi += dY

		if tmpPhi > 0 && tmpPhi < math.Pi {
			c.phi = tmpPhi
		}

		changed = true
	}

	if c.radial != c.radialC {
		c.radialC = forgemath.Lerp(c.radialC, c.radial, c.radialL)

		if math.Abs(c.radial-c.radialC) < 0.001 {
			c.radialC = c.radial
		}

		changed = true
	}

	if c.phiC != c.phi {
		c.phiC = forgemath.Lerp(c.phiC, c.phi, c.phiL)

		if math.Abs(c.phi-c.phiC) < 0.001 {
			c.phiC = c.phi
		}

		changed = true
	}

	if c.thetaC != c.theta {
		c.thetaC = forgemath.Lerp(c.thetaC, c.theta, c.thetaL)

		if math.Abs(c.theta-c.thetaC) < 0.001 {
			if c.theta > math.Pi*2 {
				for c.theta > math.Pi*2 {
					c.theta -= math.Pi * 2
				}
			} else if c.theta < -math.Pi*2 {
				for c.theta < -math.Pi*2 {
					c.theta += math.Pi * 2
				}
			}

			c.thetaC = c.theta
		}

		changed = true
	}

	if changed {
		c.move()
	}
}

func sphericalToCartesian(radial, theta, phi float64) mgl32.Vec3 {
	st, ct := math.Sincos(theta)
	sp, cp := math.Sincos(phi)
	r := float32(radial)

	return mgl32.Vec3{r * float32(sp*st), r * float32(cp), r * float32(sp*ct)}
}

func lookAt(eye, center, up mgl32.Vec3) mgl32.Quat {
	direction := center.Sub(eye).Normalize()

	rotDir := mgl32.QuatBetweenVectors(mgl32.Vec3{0, 0, -1}, direction)

	right := direction.Cross(mgl32.Vec3{0, 1, 0})
	up = right.Cross(direction)

	upCur := rotDir.Rotate(mgl32.Vec3{0, 1, 0})
	rotUp := mgl32.QuatBetweenVectors(upCur, up)

	rotTarget := rotUp.Mul(rotDir)
	return rotTarget.Inverse()
}
