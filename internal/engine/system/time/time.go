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

package time

import "github.com/haakenlabs/forge/internal/engine"

func FrameTime() float64 {
	return engine.GetTime().FixedTime()
}

func DeltaTime() float64 {
	return engine.GetTime().DeltaTime()
}

func FixedTime() float64 {
	return engine.GetTime().FixedTime()
}

func Delta() float64 {
	return engine.GetTime().Delta()
}

func Now() float64 {
	return engine.GetTime().Now()
}

func FrameStart() {
	engine.GetTime().FrameStart()
}

func FrameEnd() {
	engine.GetTime().FrameEnd()
}

func LogicTick() {
	engine.GetTime().LogicTick()
}

func LogicUpdate() bool {
	return engine.GetTime().LogicUpdate()
}
