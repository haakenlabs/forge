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

package engine

import "github.com/go-gl/mathgl/mgl32"

// Transform is a component which handles scale, rotation, and
// position transformations.
type Transform interface {
	Component

	ModelMatrix() mgl32.Mat4
	ActiveMatrix() mgl32.Mat4
	Rotation() mgl32.Quat
	Position() mgl32.Vec3
	Scale() mgl32.Vec3
	SetRotation(mgl32.Quat)
	SetPosition(mgl32.Vec3)
	SetScale(mgl32.Vec3)
	Recompute(bool)
}

// Transform is a component which handles scale, rotation, and
// position transformations.
type BaseTransform struct {
	BaseComponent

	modelMatrix  mgl32.Mat4
	activeMatrix mgl32.Mat4
	rotation     mgl32.Quat
	position     mgl32.Vec3
	scale        mgl32.Vec3
}

func (t *BaseTransform) ModelMatrix() mgl32.Mat4 {
	return t.modelMatrix
}

func (t *BaseTransform) ActiveMatrix() mgl32.Mat4 {
	return t.activeMatrix
}

func (t *BaseTransform) Rotation() mgl32.Quat {
	return t.rotation
}

func (t *BaseTransform) Position() mgl32.Vec3 {
	return t.position
}

func (t *BaseTransform) Scale() mgl32.Vec3 {
	return t.scale
}

func (t *BaseTransform) SetRotation(rotation mgl32.Quat) {
	t.rotation = rotation
	t.Recompute(true)
}

func (t *BaseTransform) SetPosition(position mgl32.Vec3) {
	t.position = position
	t.Recompute(true)
}

func (t *BaseTransform) SetScale(scale mgl32.Vec3) {
	t.scale = scale
	t.Recompute(true)
}

func (t *BaseTransform) SetRotationN(rotation mgl32.Quat) {
	t.rotation = rotation
}

func (t *BaseTransform) SetPositionN(position mgl32.Vec3) {
	t.position = position
}

func (t *BaseTransform) SetScaleN(scale mgl32.Vec3) {
	t.scale = scale
}

func (t *BaseTransform) Recompute(updateChildren bool) {
	tp := mgl32.Translate3D(t.position.X(), t.position.Y(), t.position.Z())
	tr := t.rotation.Mat4()
	ts := mgl32.Scale3D(t.scale.X(), t.scale.Y(), t.scale.Z())

	t.modelMatrix = tp.Mul4(tr.Mul4(ts))
	t.activeMatrix = t.modelMatrix

	if t.GameObject() != nil {
		if parent := t.GameObject().Parent(); parent != nil {
			t.activeMatrix = parent.Transform().ActiveMatrix().Mul4(t.modelMatrix)
		}

		if updateChildren {
			childComponents := t.GameObject().ComponentsInChildren()
			for idx := range childComponents {
				if child, ok := childComponents[idx].(Transform); ok {
					child.Recompute(false)
				}
			}
		}
	}
}

func NewTransform() *BaseTransform {
	t := &BaseTransform{
		rotation: mgl32.QuatIdent(),
		scale:    mgl32.Vec3{1.0, 1.0, 1.0},
	}

	t.SetName("Transform")
	GetInstance().MustAssign(t)

	t.Recompute(false)

	return t
}
