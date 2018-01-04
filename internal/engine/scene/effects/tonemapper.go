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

package effects

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"git.dbnservers.net/haakenlabs/forge/internal/engine"
	"git.dbnservers.net/haakenlabs/forge/internal/engine/system/asset/shader"
	forgemath "git.dbnservers.net/haakenlabs/forge/internal/math"
)

type Tonemapper struct {
	shader *engine.Shader

	exposure  float32
	exposureC float32
	exposureL float32
}

func NewTonemapper() *Tonemapper {
	e := &Tonemapper{
		exposure:  0.35,
		exposureL: 0.1,
	}

	e.exposureC = e.exposure

	e.shader = shader.MustGet("effect/tonemapper")

	return e
}

func (e *Tonemapper) Render(w engine.EffectWriter) {
	if e.exposure != e.exposureC {
		e.exposureC = forgemath.Lerp32(e.exposureC, e.exposure, e.exposureL)

		if math.Abs(float64(e.exposure-e.exposureC)) < 0.001 {
			e.exposureC = e.exposure
		}
	}

	e.shader.Bind()
	e.shader.SetSubroutine(engine.ShaderComponentFragment, "pass_basic")
	e.shader.SetUniform("f_exposure", e.exposureC)

	w.EffectPass()

	e.shader.Unbind()
}

func (e *Tonemapper) Type() engine.EffectType {
	return engine.EffectTypeTonemapper
}

func (e *Tonemapper) Exposure() float32 {
	return e.exposure
}

func (e *Tonemapper) SetExposure(exp float32) {
	e.exposure = mgl32.Clamp(exp, 0.1, 2.0)
}
