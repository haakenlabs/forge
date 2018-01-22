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

type Textbox struct {
	BaseComponent

	value string

	onChangeFunc func(string)

	background *Graphic
	text       *Text
}

func (w *Textbox) UIDraw() {
	w.background.Draw()
	w.text.Draw()
}

func NewTextbox() *Textbox {
	w := &Textbox{
		value: "Text",
	}

	w.SetName("UITextbox")
	engine.GetInstance().MustAssign(w)

	return w
}

func CreateTextbox(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	textbox := NewTextbox()

	textbox.background = NewGraphic()
	textbox.text = NewText()

	object.AddComponent(textbox)
	object.AddComponent(textbox.background)
	object.AddComponent(textbox.text)

	return object
}
