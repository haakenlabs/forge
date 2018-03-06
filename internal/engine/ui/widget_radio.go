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

type RadioState int

const (
	RadioStateOff RadioState = iota
	RadioStateMixed
	RadioStateOn
)

type RadioGroup struct {
	BaseComponent

	radios []*Radio
}

type Radio struct {
	BaseComponent

	checked RadioState

	backgroundColor engine.Color
	tint            engine.Color

	onChangeFunc func(RadioState)

	background *Graphic
	check      *Graphic
}

func (w *Radio) UIDraw() {
	w.background.Draw()
	w.check.Draw()
}

func (w *RadioGroup) AddRadio(radio ...*Radio) {
	w.radios = append(w.radios, radio...)
}

func NewRadio() *Radio {
	w := &Radio{}

	w.SetName("UIRadio")
	engine.GetInstance().MustAssign(w)

	return w
}

func NewRadioGroup() *RadioGroup {
	w := &RadioGroup{}

	w.SetName("UIRadioGroup")
	engine.GetInstance().MustAssign(w)

	return w
}

func CreateRadio(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	radio := NewRadio()

	radio.background = NewGraphic()
	radio.check = NewGraphic()

	object.AddComponent(radio)
	object.AddComponent(radio.background)
	object.AddComponent(radio.check)

	return object
}

func CreateRadioGroup(name string, radios ...*Radio) *engine.GameObject {
	object := CreateGenericObject(name)

	group := NewRadioGroup()
	group.AddRadio(radios...)

	object.AddComponent(group)

	return object
}
