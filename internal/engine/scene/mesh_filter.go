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

package scene

import (
	"github.com/haakenlabs/forge/internal/engine"
	"github.com/haakenlabs/forge/internal/engine/system/instance"
)

// MeshFilter is a component that allows a Mesh to be associated with an entity.
type MeshFilter struct {
	engine.BaseComponent

	mesh *engine.Mesh
}


// NewMeshFilter creates a new MeshFilter component.
func NewMeshFilter(mesh *engine.Mesh) *MeshFilter {
	m := &MeshFilter{
		mesh: mesh,
	}

	m.SetName("MeshFilter")
	instance.MustAssign(m)

	return m
}

// MeshFilterComponent gets the first occurrence of MeshFilter from the entity.
func MeshFilterComponent(g *engine.GameObject) *MeshFilter {
	c := g.Components()
	for i := range c {
		if ct, ok := c[i].(*MeshFilter); ok {
			return ct
		}
	}

	return nil
}

// Mesh gets the Mesh associated with this MeshFilter.
func (m *MeshFilter) Mesh() *engine.Mesh {
	return m.mesh
}

// SetMesh sets the Mesh for this MeshFilter.
func (m *MeshFilter) SetMesh(mesh *engine.Mesh) {
	m.mesh = mesh
}
