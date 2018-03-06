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

type Label struct {
	BaseComponent

	value     string
	textColor engine.Color

	text *Text
}

func NewLabel() *Label {
	w := &Label{
		value:     "Label",
		textColor: Styles.PrimaryTextColor,
	}

	w.SetName("UILabel")
	engine.GetInstance().MustAssign(w)

	return w
}

func (w *Label) UIDraw() {
	w.text.Draw()
}

func (w *Label) SetValue(value string) {
	w.value = value

	w.text.SetValue(value)
}

func (w *Label) SetTextColor(color engine.Color) {
	w.textColor = color
}

func (w *Label) Value() string {
	return w.value
}

func (w *Label) TextColor() engine.Color {
	return w.textColor
}

func (w *Label) SetFontSize(size int32) {
	w.text.SetFontSize(size)
}

func (w *Label) FontSize(size int32) int32 {
	return w.text.fontSize
}

func LabelComponent(g *engine.GameObject) *Label {
	c := g.Components()
	for i := range c {
		if ct, ok := c[i].(*Label); ok {
			return ct
		}
	}

	return nil
}

func CreateLabel(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	label := NewLabel()

	label.text = NewText()
	label.text.SetFontSize(12)

	object.AddComponent(label)
	object.AddComponent(label.text)

	return object
}
