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

type CheckState int

const (
	CheckStateOff CheckState = iota
	CheckStateMixed
	CheckStateOn
)

type CheckboxGroup struct {
	BaseComponent

	checkboxes []*Checkbox
}

type Checkbox struct {
	BaseComponent

	state CheckState

	backgroundColor engine.Color
	tint            engine.Color

	onChangeFunc func(CheckState)

	background *Graphic
	check      *Graphic
}

func (w *Checkbox) UIDraw() {
	w.background.Draw()
	w.check.Draw()
}

func (w *CheckboxGroup) AddCheckbox(checkbox ...*Checkbox) {
	w.checkboxes = append(w.checkboxes, checkbox...)
}

func NewCheckbox() *Checkbox {
	w := &Checkbox{}

	w.SetName("UICheckbox")
	engine.GetInstance().MustAssign(w)

	return w
}

func NewCheckboxGroup() *CheckboxGroup {
	w := &CheckboxGroup{}

	w.SetName("UICheckboxGroup")
	engine.GetInstance().MustAssign(w)

	return w
}

func CreateCheckbox(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	slider := NewCheckbox()

	slider.background = NewGraphic()
	slider.check = NewGraphic()

	object.AddComponent(slider)
	object.AddComponent(slider.background)
	object.AddComponent(slider.check)

	return object
}

func CreateCheckboxGroup(name string, checkboxes ...*Checkbox) *engine.GameObject {
	object := CreateGenericObject(name)

	group := NewCheckboxGroup()
	group.AddCheckbox(checkboxes...)

	object.AddComponent(group)

	return object
}
