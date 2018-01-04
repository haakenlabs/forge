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

package window

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/haakenlabs/forge/internal/engine"
	"github.com/haakenlabs/forge/internal/math"
)

func AspectRatio() float32 {
	return engine.GetWindow().AspectRatio()
}

func CenterWindow() {
	engine.GetWindow().CenterWindow()
}

func ClearBuffers() {
	engine.GetWindow().ClearBuffers()
}

func EnableVsync(enable bool) {
	engine.GetWindow().EnableVsync(enable)
}

func Resolution() math.IVec2 {
	return engine.GetWindow().Resolution()
}

func SetSize(size math.IVec2) {
	engine.GetWindow().SetSize(size)
}

func SwapBuffers() {
	engine.GetWindow().SwapBuffers()
}

func Vsync() bool {
	return engine.GetWindow().Vsync()
}

func GLFWWindow() *glfw.Window {
	return engine.GetWindow().GLFWWindow()
}

func OrthoMatrix() mgl32.Mat4 {
	return engine.GetWindow().OrthoMatrix()
}
