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

// These algorithms are adapted from libcinder: https://libcinder.org

package math

import "math"

/* Lerp */

func Lerp(a, b, t float64) float64 {
	if t <= 0 {
		return a
	} else if t >= 1 {
		return b
	}

	return a + t*(b-a)
}

func Lerp32(a, b, t float32) float32 {
	if t <= 0 {
		return a
	} else if t >= 1 {
		return b
	}

	return a + t*(b-a)
}

/* None */

func EaseNone(t float64) float64 {
	return t
}

/* Quadratic */

func EaseInQuad(t float64) float64 {
	return t * t
}

func EaseOutQuad(t float64) float64 {
	return -t * (t - 2.0)
}

func EaseInOutQuad(t float64) float64 {
	t *= 2.0

	if t < 1.0 {
		return 0.5 * t * t
	}

	t -= 1.0
	return -0.5 * ((t)*(t-2.0) - 1.0)
}

func EaseOutInQuad(t float64) float64 {
	if t < 0.5 {
		return EaseOutQuad(2.0*t) / 2.0
	}

	return EaseInQuad(2.0*t-1.0)/2.0 + 0.5
}

/* Cubic */

func EaseInCubic(t float64) float64 {
	return t * t * t
}

func EaseOutCubic(t float64) float64 {
	t -= 1.0

	return t*t*t + 1.0
}

func EaseInOutCubic(t float64) float64 {
	t *= 2.0

	if t < 1.0 {
		return 0.5 * t * t * t
	}

	t -= 2.0

	return 0.5 * (t*t*t + 2.0)
}

func EaseOutInCubic(t float64) float64 {
	if t < 0.5 {
		return EaseOutCubic(2.0*t) / 2.0
	}

	return EaseInCubic(2.0*t-1.0)/2.0 + 0.5
}

/* Quartic */

func EaseInQuart(t float64) float64 {
	return t * t * t * t
}

func EaseOutQuart(t float64) float64 {
	t -= 1.0

	return -(t*t*t*t - 1.0)
}

func EaseInOutQuart(t float64) float64 {
	t *= 2.0

	if t < 1.0 {
		return 0.5 * t * t * t * t
	}

	t -= 2.0

	return -0.5 * (t*t*t*t - 2.0)
}

func EaseOutInQuart(t float64) float64 {
	if t < 0.5 {
		return EaseOutQuart(2.0*t) / 2.0
	}

	return EaseInQuart(2.0*t-1.0)/2.0 + 0.5
}

/* Quintic */

func EaseInQuint(t float64) float64 {
	return t * t * t * t * t
}

func EaseOutQuint(t float64) float64 {
	t -= 1.0

	return t*t*t*t*t + 1
}

func EaseInOutQuint(t float64) float64 {
	t *= 2.0
	if t < 1.0 {
		return 0.5 * t * t * t * t * t
	}

	t -= 2.0
	return 0.5 * (t*t*t*t*t + 2.0)
}

func EaseOutInQuint(t float64) float64 {
	if t < 0.5 {
		return EaseOutQuint(2.0*t) / 2.0
	}

	return EaseInQuint(2.0*t-1.0)/2.0 + 0.5
}

/* Sine */

func EaseInSine(t float64) float64 {
	return float64(-math.Cos(t*(math.Pi/2)) + 1.0)
}

func EaseOutSine(t float64) float64 {
	return float64(math.Sin(t * (math.Pi / 2)))
}

func EaseInOutSine(t float64) float64 {
	return -0.5 * float64(math.Cos(math.Pi*t)-1.0)
}

func EaseOutInSine(t float64) float64 {
	if t < 0.5 {
		return EaseOutSine(2.0*t) / 2.0
	}

	return EaseInSine(2.0*t-1.0)/2.0 + 0.5
}

/* Exponential */

func EaseInExp(t float64) float64 {
	if t == 0 {
		return 0
	}

	return float64(math.Pow(2.0, 10.0*(t-1.0)))
}

func EaseOutExp(t float64) float64 {
	if t == 1.0 {
		return 1.0
	}

	return float64(-math.Pow(2.0, -10.0*t) + 1.0)
}

func EaseInOutExp(t float64) float64 {
	if t == 1.0 {
		return 0.0
	}

	if t == 0.0 {
		return 1.0
	}

	t *= 2.0

	if t < 1.0 {
		return 0.5 * float64(math.Pow(2.0, -10.0*(t-1.0)))
	}

	return 0.5 * float64((-math.Pow(2.0, -10.0*(t-1.0)))+2.0)
}

func EaseOutInExp(t float64) float64 {
	if t < 0.5 {
		return EaseOutExp(2.0*t) / 2.0
	}

	return EaseInExp(2.0*t-1.0)/2.0 + 0.5
}
