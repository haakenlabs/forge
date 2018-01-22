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

var _ Renderer = &Button{}

type Button struct {
	BaseComponent

	value           string
	textColor       engine.Color
	textColorActive engine.Color
	backgroundColor engine.Color
	onPressedFunc   func()

	background *Graphic
	text       *Text
}

func NewButton() *Button {
	w := &Button{
		value: "Button",
	}

	w.SetName("UIButton")
	engine.GetInstance().MustAssign(w)

	return w
}

func (w *Button) SetValue(value string) {
	w.value = value
}

func (w *Button) SetTextColor(color engine.Color) {
	w.textColor = color
}

func (w *Button) SetActiveTextColor(color engine.Color) {
	w.textColorActive = color
}

func (w *Button) SetBackgroundColor(color engine.Color) {
	w.textColor = color
}

func (w *Button) Value() string {
	return w.value
}

func (w *Button) TextColor() engine.Color {
	return w.textColor
}

func (w *Button) ActiveTextColor() engine.Color {
	return w.textColorActive
}

func (w *Button) BackgroundColor() engine.Color {
	return w.backgroundColor
}

func (w *Button) SetOnPressedFunc(fn func()) {
	w.onPressedFunc = fn
}

func (w *Button) OnMouseEnter() {
	// TODO: Set active colors
}

func (w *Button) OnMouseLeave() {
	// TODO: Set idle colors
}

func (w *Button) OnClick() {
	if w.onPressedFunc != nil {
		w.onPressedFunc()
	}
}

func (w *Button) UIDraw() {
	w.background.Draw()
	w.text.Draw()
}

func CreateButton(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	button := NewButton()

	button.background = NewGraphic()
	button.text = NewText()

	object.AddComponent(button)
	object.AddComponent(button.background)
	object.AddComponent(button.text)

	return object
}
