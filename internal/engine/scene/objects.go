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
	"github.com/haakenlabs/forge/internal/engine/system/asset/mesh"
	"github.com/haakenlabs/forge/internal/engine/system/asset/shader"
)

func CreateCube(name string) *engine.GameObject {
	// Create necessary objects
	object := engine.NewGameObject(name)
	meshRenderer := NewMeshRenderer()
	material := engine.NewMaterial()
	meshFilter := NewMeshFilter(mesh.MustGet("cube"))

	// Assign stuff
	material.SetShader(shader.MustGet("standard"))
	meshRenderer.SetMaterial(material)

	// Add components to object
	object.AddComponent(meshRenderer)
	object.AddComponent(meshFilter)

	return object
}

func CreateOrb(name string) *engine.GameObject {
	// Create necessary objects
	object := engine.NewGameObject(name)
	meshRenderer := NewMeshRenderer()
	material := engine.NewMaterial()
	meshFilter := NewMeshFilter(mesh.MustGet("orb"))

	// Assign stuff
	material.SetShader(shader.MustGet("standard"))
	meshRenderer.SetMaterial(material)

	// Add components to object
	object.AddComponent(meshRenderer)
	object.AddComponent(meshFilter)

	return object
}

func CreateSphere(name string) *engine.GameObject {
	// Create necessary objects
	object := engine.NewGameObject(name)
	meshRenderer := NewMeshRenderer()
	material := engine.NewMaterial()
	meshFilter := NewMeshFilter(mesh.MustGet("sphere"))

	// Assign stuff
	material.SetShader(shader.MustGet("standard"))
	meshRenderer.SetMaterial(material)

	// Add components to object
	object.AddComponent(meshRenderer)
	object.AddComponent(meshFilter)

	return object
}

func CreateCamera(name string, hdr bool, path engine.RenderPath) *engine.GameObject {
	object := engine.NewGameObject(name)
	cameraComponent := engine.NewCamera(path, hdr)

	object.AddComponent(cameraComponent)

	return object
}
