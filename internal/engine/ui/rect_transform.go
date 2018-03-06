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

type AnchorPreset uint8
type PivotPreset uint8

const (
	AnchorTopLeft AnchorPreset = iota
	AnchorTopCenter
	AnchorTopRight
	AnchorMiddleLeft
	AnchorMiddleCenter
	AnchorMiddleRight
	AnchorBottomLeft
	AnchorBottomCenter
	AnchorBottomRight
	StretchAnchorLeft
	StretchAnchorCenter
	StretchAnchorRight
	StretchAnchorTop
	StretchAnchorMiddle
	StretchAnchorBottom
	StretchAnchorAll
)

const (
	PivotTopLeft PivotPreset = iota
	PivotTopCenter
	PivotTopRight
	PivotMiddleLeft
	PivotMiddleCenter
	PivotMiddleRight
	PivotBottomLeft
	PivotBottomCenter
	PivotBottomRight
)

type RectTransform struct {
	engine.BaseTransform

	rect      Rect
	anchorMax mgl32.Vec2
	anchorMin mgl32.Vec2
	offsetMax mgl32.Vec2
	offsetMin mgl32.Vec2
	pivot     mgl32.Vec2
}

func NewRectTransform() *RectTransform {
	t := &RectTransform{}

	t.SetRotationN(mgl32.QuatIdent())
	t.SetScaleN(mgl32.Vec3{1.0, 1.0, 1.0})

	t.SetName("RectTransform")
	engine.GetInstance().MustAssign(t)

	return t
}

func RectTransformComponent(g *engine.GameObject) *RectTransform {
	c := g.Components()
	for i := range c {
		if ct, ok := c[i].(*RectTransform); ok {
			return ct
		}
	}

	return nil
}

func (t *RectTransform) Rect() Rect {
	return t.rect
}

func (t *RectTransform) AnchorMax() mgl32.Vec2 {
	return t.anchorMax
}

func (t *RectTransform) AnchorMin() mgl32.Vec2 {
	return t.anchorMin
}

func (t *RectTransform) OffsetMax() mgl32.Vec2 {
	return t.offsetMax
}

func (t *RectTransform) OffsetMin() mgl32.Vec2 {
	return t.offsetMin
}

func (t *RectTransform) Pivot() mgl32.Vec2 {
	return t.pivot
}

func (t *RectTransform) Size() mgl32.Vec2 {
	return t.rect.Size()
}

func (t *RectTransform) SetRect(rect Rect) {
	t.rect = rect
	t.ComputeOffsets()
	t.Recompute(true)
}

func (t *RectTransform) SetPosition2D(position mgl32.Vec2) {
	t.rect.SetOrigin(position)
	t.ComputeOffsets()
	t.Recompute(true)
}

func (t *RectTransform) SetSize(size mgl32.Vec2) {
	t.rect.SetSize(size)
	t.ComputeOffsets()

	t.Recompute(true)
}

func (t *RectTransform) SetAnchorMax(anchor mgl32.Vec2) {
	t.anchorMax = anchor
	t.ComputeOffsets()
}

func (t *RectTransform) SetAnchorMin(anchor mgl32.Vec2) {
	t.anchorMin = anchor
	t.ComputeOffsets()
}

func (t *RectTransform) SetPivot(pivot mgl32.Vec2) {
	t.pivot = pivot
	t.ComputeOffsets()
}

func (t *RectTransform) SetAnchorPreset(preset AnchorPreset) {
	switch preset {
	case AnchorTopLeft:
		t.anchorMin = mgl32.Vec2{}
		t.anchorMax = mgl32.Vec2{}
		break
	case AnchorTopCenter:
		t.anchorMin = mgl32.Vec2{0.5, 0}
		t.anchorMax = mgl32.Vec2{0.5, 0}
		break
	case AnchorTopRight:
		t.anchorMin = mgl32.Vec2{1, 0}
		t.anchorMax = mgl32.Vec2{1, 0}
		break
	case AnchorMiddleLeft:
		t.anchorMin = mgl32.Vec2{0, 0.5}
		t.anchorMax = mgl32.Vec2{0, 0.5}
		break
	case AnchorMiddleCenter:
		t.anchorMin = mgl32.Vec2{0.5, 0.5}
		t.anchorMax = mgl32.Vec2{0.5, 0.5}
		break
	case AnchorMiddleRight:
		t.anchorMin = mgl32.Vec2{1, 0.5}
		t.anchorMax = mgl32.Vec2{1, 0.5}
		break
	case AnchorBottomLeft:
		t.anchorMin = mgl32.Vec2{0, 1}
		t.anchorMax = mgl32.Vec2{0, 1}
		break
	case AnchorBottomCenter:
		t.anchorMin = mgl32.Vec2{0.5, 1}
		t.anchorMax = mgl32.Vec2{0.5, 1}
		break
	case AnchorBottomRight:
		t.anchorMin = mgl32.Vec2{1, 1}
		t.anchorMax = mgl32.Vec2{1, 1}
		break
	case StretchAnchorLeft:
		t.anchorMin = mgl32.Vec2{0, 0}
		t.anchorMax = mgl32.Vec2{0, 1}
		break
	case StretchAnchorCenter:
		t.anchorMin = mgl32.Vec2{0.5, 0}
		t.anchorMax = mgl32.Vec2{0.5, 1}
		break
	case StretchAnchorRight:
		t.anchorMin = mgl32.Vec2{1, 0}
		t.anchorMax = mgl32.Vec2{1, 1}
		break
	case StretchAnchorTop:
		t.anchorMin = mgl32.Vec2{0, 0}
		t.anchorMax = mgl32.Vec2{1, 0}
		break
	case StretchAnchorMiddle:
		t.anchorMin = mgl32.Vec2{0, 0.5}
		t.anchorMax = mgl32.Vec2{1, 0.5}
		break
	case StretchAnchorBottom:
		t.anchorMin = mgl32.Vec2{0, 1}
		t.anchorMax = mgl32.Vec2{1, 1}
		break
	case StretchAnchorAll:
		t.anchorMin = mgl32.Vec2{0, 0}
		t.anchorMax = mgl32.Vec2{1, 1}
		break
	default:
		break
	}

	t.ComputeOffsets()
}

func (t *RectTransform) SetPivotPreset(preset PivotPreset) {
	switch preset {
	case PivotTopLeft:
		t.pivot = mgl32.Vec2{0, 0}
		break
	case PivotTopCenter:
		t.pivot = mgl32.Vec2{0.5, 0}
		break
	case PivotTopRight:
		t.pivot = mgl32.Vec2{1, 0}
		break
	case PivotMiddleLeft:
		t.pivot = mgl32.Vec2{0, 0.5}
		break
	case PivotMiddleCenter:
		t.pivot = mgl32.Vec2{0.5, 0.5}
		break
	case PivotMiddleRight:
		t.pivot = mgl32.Vec2{1, 0.5}
		break
	case PivotBottomLeft:
		t.pivot = mgl32.Vec2{0, 1}
		break
	case PivotBottomCenter:
		t.pivot = mgl32.Vec2{0.5, 1}
		break
	case PivotBottomRight:
		t.pivot = mgl32.Vec2{1, 1}
		break
	default:
		break
	}

	t.ComputeOffsets()
}

func (t *RectTransform) SetPresets(anchor AnchorPreset, pivot PivotPreset) {
	t.SetAnchorPreset(anchor)
	t.SetPivotPreset(pivot)
}

func (t *RectTransform) Recompute(updateChildren bool) {
	var (
		anchorMinActual mgl32.Vec2
		anchorMaxActual mgl32.Vec2
		apparentSize    mgl32.Vec2
	)

	if parent := t.ParentTransform(); parent != nil {
		anchorMinActual = mgl32.Vec2{t.rect.Size().X() * t.anchorMin.X(), t.rect.Size().Y() * t.anchorMin.Y()}
		anchorMaxActual = mgl32.Vec2{t.rect.Size().X() * t.anchorMax.X(), t.rect.Size().Y() * t.anchorMax.Y()}
		apparentSize = anchorMaxActual.Add(t.offsetMax).Sub(anchorMinActual.Add(t.offsetMin))

		t.rect.SetSize(apparentSize)
	}

	t.SetPosition(anchorMinActual.Add(t.offsetMin).Vec3(0))

	t.BaseTransform.Recompute(updateChildren)

	//if t.GameObject() != nil {
	//	x := t.GameObject().Components()
	//	for i := range x {
	//		if o, ok := x[i].(Component); ok {
	//			o.TransformChanged()
	//		}
	//	}
	//}
}

func (t *RectTransform) ComputeOffsets() {
	parent := t.ParentTransform()
	if parent == nil {
		t.offsetMin = t.rect.Min()
		t.offsetMax = t.rect.Max()
		return
	}

	pivotSkew := mgl32.Vec2{t.rect.Size().X() * t.pivot.X(), t.rect.Size().Y() * t.pivot.Y()}

	offsetMinX := (parent.Size().X()*t.anchorMin.X() - pivotSkew.X() + t.rect.Min().X()) - parent.Size().X()*t.anchorMin.X()
	offsetMinY := (parent.Size().Y()*t.anchorMin.Y() - pivotSkew.Y() + t.rect.Min().Y()) - parent.Size().Y()*t.anchorMin.Y()
	offsetMaxX := (parent.Size().X()*t.anchorMax.X() - pivotSkew.X() + t.rect.Max().X()) - parent.Size().X()*t.anchorMax.X()
	offsetMaxY := (parent.Size().Y()*t.anchorMax.Y() - pivotSkew.Y() + t.rect.Max().Y()) - parent.Size().Y()*t.anchorMax.Y()

	t.offsetMin = mgl32.Vec2{offsetMinX, offsetMinY}
	t.offsetMax = mgl32.Vec2{offsetMaxX, offsetMaxY}
}

func (t *RectTransform) WorldPosition() mgl32.Vec2 {
	return t.ActiveMatrix().Col(3).Vec2()
}

func (t *RectTransform) ParentTransform() *RectTransform {
	if t.GameObject() != nil {
		if parent := t.GameObject().Parent(); parent != nil {
			if obj, ok := parent.Transform().(*RectTransform); ok {
				return obj
			}
		}
	} else {

	}

	return nil
}
