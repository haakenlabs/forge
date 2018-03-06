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

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/haakenlabs/forge/internal/engine"
)

type Image struct {
	BaseComponent

	graphic *Graphic
}

func (w *Image) UIDraw() {
	w.graphic.Draw()
}

func (w *Image) Color() engine.Color {
	return w.graphic.Color()
}

func (w *Image) Texture() *engine.Texture2D {
	return w.graphic.Texture()
}

func (w *Image) SetColor(color engine.Color) {
	w.graphic.SetColor(color)
}

func (w *Image) SetTexture(texture *engine.Texture2D) {
	w.graphic.SetTexture(texture)
}

func (w *Image) OnTransformChanged() {
	w.graphic.Refresh()
}

func (w *Image) Start() {
	w.graphic.Refresh()
}

func NewImage() *Image {
	w := &Image{}

	w.SetName("UIImage")
	engine.GetInstance().MustAssign(w)

	return w
}

func ImageComponent(g *engine.GameObject) *Image {
	c := g.Components()
	for i := range c {
		if ct, ok := c[i].(*Image); ok {
			return ct
		}
	}

	return nil
}

func CreateImage(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	image := NewImage()

	image.graphic = NewGraphic()

	object.AddComponent(image)
	object.AddComponent(image.graphic)

	return object
}

func CreatePanel(name string) *engine.GameObject {
	object := CreateGenericObject(name)

	rt := RectTransformComponent(object)
	rt.SetSize(mgl32.Vec2{480, 320})

	image := NewImage()
	image.graphic = NewGraphic()
	image.graphic.SetColor(Styles.BackgroundColor)

	object.AddComponent(image)
	object.AddComponent(image.graphic)

	return object
}
