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

package math

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
)

type DVec2 [2]float64

func (v DVec2) X() float64 {
	return v[0]
}

func (v DVec2) Y() float64 {
	return v[1]
}

func (v DVec2) String() string {
	return fmt.Sprintf("DVec2(%d, %d)", v.X(), v.Y())
}

type IVec2 [2]int32

func (v IVec2) X() int32 {
	return v[0]
}

func (v IVec2) Y() int32 {
	return v[1]
}

func (v IVec2) String() string {
	return fmt.Sprintf("IVec2(%d, %d)", v.X(), v.Y())
}

func (v IVec2) Vec2() mgl32.Vec2 {
	return mgl32.Vec2{float32(v[0]), float32(v[1])}
}
