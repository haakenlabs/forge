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
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/spf13/cast"
)

func ToVec2(i interface{}) mgl32.Vec2 {
	v, _ := ToVec2E(i)

	return v
}

func ToIVec2(i interface{}) IVec2 {
	v, e := ToIVec2E(i)

	if e != nil {
		panic(e)
	}

	return v
}

func ToVec2E(i interface{}) (mgl32.Vec2, error) {
	if i == nil {
		return mgl32.Vec2{}, fmt.Errorf("unable to cast %#v to mgl32.Vec2", i)
	}

	switch v := i.(type) {
	case mgl32.Vec2:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		if s.Len() != 2 {
			return mgl32.Vec2{}, fmt.Errorf("unable to cast %#v to mgl32.Vec2", i)
		}
		a := mgl32.Vec2{}
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToIntE(s.Index(j).Interface())
			if err != nil {
				return mgl32.Vec2{}, fmt.Errorf("unable to cast %#v to mgl32.Vec2", i)
			}
			a[j] = float32(val)
		}

		return a, nil

	default:
		return mgl32.Vec2{}, fmt.Errorf("unable to cast %#v to mgl32.Vec2", i)
	}
}

func ToIVec2E(i interface{}) (IVec2, error) {
	if i == nil {
		return IVec2{}, fmt.Errorf("unable to cast %#v to IVec2", i)
	}

	switch v := i.(type) {
	case IVec2:
		return v, nil
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		if s.Len() != 2 {
			return IVec2{}, fmt.Errorf("unable to cast %#v to IVec2", i)
		}
		a := IVec2{}
		for j := 0; j < s.Len(); j++ {
			val, err := cast.ToIntE(s.Index(j).Interface())
			if err != nil {
				return IVec2{}, fmt.Errorf("unable to cast %#v to IVec2", i)
			}
			a[j] = int32(val)
		}

		return a, nil

	default:
		return IVec2{}, fmt.Errorf("unable to cast %#v to IVec2", i)
	}
}
