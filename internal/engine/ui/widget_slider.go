/*
Copyright (c) 2018 HaakenLabs

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

package ui

import "github.com/haakenlabs/forge/internal/engine"

type Slider struct {
	BaseComponent

	value float64
	min   float64
	max   float64

	backgroundColor engine.Color
	tint            engine.Color

	onChangeFunc func(float64)

	background  *Graphic
	activeTrack *Graphic
	thumb       *Graphic
}

func (w *Slider) UIDraw() {
	w.background.Draw()
	w.activeTrack.Draw()
	w.thumb.Draw()
}

func NewSlider() *Slider {
	w := &Slider{
		value: 0.5,
		min:   0.0,
		max:   1.0,
	}

	w.SetName("UISlider")
	engine.GetInstance().MustAssign(w)

	return w
}

func CreateSlider(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	slider := NewSlider()

	slider.background = NewGraphic()
	slider.activeTrack = NewGraphic()
	slider.thumb = NewGraphic()

	object.AddComponent(slider)
	object.AddComponent(slider.background)
	object.AddComponent(slider.activeTrack)
	object.AddComponent(slider.thumb)

	return object
}
