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

package input

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/haakenlabs/forge/internal/engine"
)

func KeyDown(key glfw.Key) bool {
	return engine.GetWindow().KeyDown(key)
}

func KeyUp(key glfw.Key) bool {
	return engine.GetWindow().KeyUp(key)
}

func KeyPressed() bool {
	return engine.GetWindow().KeyPressed()
}

func MouseDown(button glfw.MouseButton) bool {
	return engine.GetWindow().MouseDown(button)
}

func MouseUp(button glfw.MouseButton) bool {
	return engine.GetWindow().MouseUp(button)
}

func MouseWheelX() float64 {
	return engine.GetWindow().MouseWheelX()
}

func MouseWheelY() float64 {
	return engine.GetWindow().MouseWheelY()
}

func MouseWheel() bool {
	return engine.GetWindow().MouseWheel()
}

func MouseMoved() bool {
	return engine.GetWindow().MouseMoved()
}

func MousePressed() bool {
	return engine.GetWindow().MousePressed()
}

func MousePosition() mgl32.Vec2 {
	return engine.GetWindow().MousePosition()
}

func WindowResized() bool {
	return engine.GetWindow().WindowResized()
}

func ShouldClose() bool {
	return engine.GetWindow().ShouldClose()
}

func HandleEvents() {
	engine.GetWindow().HandleEvents()
}

func HasEvents() bool {
	return engine.GetWindow().HasEvents()
}
