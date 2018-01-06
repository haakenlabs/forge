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
	"github.com/go-gl/glfw/v3.2/glfw"
)

var _ System = &Time{}

const SysNameTime = "time"

// Time implements a time system.
type Time struct {
	frameTime     float64
	deltaTime     float64
	nextLogicTick float64
	frame         uint64
}

// Setup sets up the System.
func (t *Time) Setup() error {
	return nil
}

// Teardown tears down the System.
func (t *Time) Teardown() {

}

// Name returns the name of the System.
func (t *Time) Name() string {
	return SysNameTime
}

func (t *Time) FrameTime() float64 {
	return t.frameTime
}

func (t *Time) DeltaTime() float64 {
	return t.deltaTime
}

func (t *Time) FixedTime() float64 {
	return fixedTime
}

func (t *Time) Delta() float64 {
	return t.deltaTime
}

func (t *Time) Now() float64 {
	return glfw.GetTime()
}

func (t *Time) FrameStart() {
	t.frameTime = t.Now()
}

func (t *Time) FrameEnd() {
	t.deltaTime = t.Now() - t.frameTime
	t.frame++
}

func (t *Time) Frame() uint64 {
	return t.frame
}

func (t *Time) LogicTick() {
	t.nextLogicTick += fixedTime
}

func (t *Time) LogicUpdate() bool {
	return t.Now() > t.nextLogicTick
}

// NewTime creates a new time system.
func NewTime() *Time {
	return &Time{}
}

// GetTime gets the time system from the current app.
func GetTime() *Time {
	return CurrentApp().MustSystem(SysNameTime).(*Time)
}
