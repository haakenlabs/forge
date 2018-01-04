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

package sg

import (
	"fmt"
	"reflect"
	"testing"
)

type Object struct {
	active     bool
	instanceId uint32
}

func newObject(id uint32) *Object {
	o := &Object{
		active:     true,
		instanceId: id,
	}

	return o
}

func (o *Object) Active() bool {
	return o.active
}

func (o *Object) ID() uint32 {
	return o.instanceId
}

func TestGraph_AddVertex(t *testing.T) {
	g := NewGraph()

	object1 := newObject(1)
	object2 := newObject(2)
	object3 := newObject(3)

	if _, err := g.AddVertex(object1); err != nil {
		t.Error(err)
	}
	if _, err := g.AddVertex(object2); err != nil {
		t.Error(err)
	}
	if _, err := g.AddVertex(object3); err != nil {
		t.Error(err)
	}

	// Validate vertexList length.
	if c := len(g.vertexList); c != 3 {
		t.Error("len(g.vertexList) expected 3, got:", c)
	}
}

func TestGraph_RemoveVertex(t *testing.T) {
	g := NewGraph()

	// Create some objects.
	object1 := newObject(1)
	object2 := newObject(2)
	object3 := newObject(3)

	var err error
	var obj1Desc VertexDescriptor
	var obj2Desc VertexDescriptor
	var obj3Desc VertexDescriptor

	// Add some vertices.
	obj1Desc, err = g.AddVertex(object1)
	if err != nil {
		t.Error(err)
	}
	obj2Desc, err = g.AddVertex(object2)
	if err != nil {
		t.Error(err)
	}
	obj3Desc, err = g.AddVertex(object3)
	if err != nil {
		t.Error(err)
	}

	// Remove a vertex.
	err = g.RemoveVertex(obj1Desc)
	if err != nil {
		t.Error(err)
	}

	// Validate vertexList length.
	if c := len(g.vertexList); c != 2 {
		t.Error("len(g.vertexList) expected 2, got:", c)
	}

	// Remove the other vertices.
	err = g.RemoveVertex(obj2Desc)
	if err != nil {
		t.Error(err)
	}
	err = g.RemoveVertex(obj3Desc)
	if err != nil {
		t.Error(err)
	}

	// Validate vertexList length.
	if c := len(g.vertexList); c != 0 {
		t.Error("len(g.vertexList) expected 0, got:", c)
	}
}

func TestGraph_GetVertex(t *testing.T) {
	g := NewGraph()

	// Create some objects.
	object1 := newObject(1)
	object2 := newObject(2)

	var err error
	var obj1Desc VertexDescriptor
	var obj2Desc VertexDescriptor

	// Add some vertices.
	obj1Desc, err = g.AddVertex(object1)
	if err != nil {
		t.Error(err)
	}
	obj2Desc, err = g.AddVertex(object2)
	if err != nil {
		t.Error(err)
	}

	v, vErr := g.GetVertexById(object1.instanceId)
	if vErr != nil {
		t.Error(vErr)
	}
	if v != obj1Desc {
		t.Errorf("g.GetVertexById(object1.instanceId) mismatch: %d != %d", obj1Desc, v)
	}

	v, vErr = g.GetVertexByObject(object2)
	if vErr != nil {
		t.Error(vErr)
	}
	if v != obj2Desc {
		t.Errorf("g.GetVertexByObject(object2) mismatch: %d != %d", obj2Desc, v)
	}

	obj := g.GetObjectAtVertex(obj2Desc)
	if obj == nil {
		t.Error("g.GetObjectAtVertex(obj2Desc) is nil")
	}
	if obj.ID() != object2.instanceId {
		t.Errorf("g.GetObjectAtVertex(obj2Desc) mismatch: %d != %d", obj2Desc, v)
	}
}

func TestGraph_AddEdge(t *testing.T) {
	g := NewGraph()

	var err error
	var obj1Desc VertexDescriptor
	var obj2Desc VertexDescriptor
	var obj3Desc VertexDescriptor
	var obj4Desc VertexDescriptor

	// Create some objects.
	object1 := newObject(1)
	object2 := newObject(2)
	object3 := newObject(3)
	object4 := newObject(4)

	// Add vertices.
	obj1Desc, err = g.AddVertex(object1)
	if err != nil {
		t.Error(err)
	}
	obj2Desc, err = g.AddVertex(object2)
	if err != nil {
		t.Error(err)
	}
	obj3Desc, err = g.AddVertex(object3)
	if err != nil {
		t.Error(err)
	}
	obj4Desc, err = g.AddVertex(object4)
	if err != nil {
		t.Error(err)
	}

	// Add some edges.
	if err = g.AddEdge(obj1Desc, obj2Desc); err != nil {
		t.Error(err)
	}
	if err = g.AddEdge(obj1Desc, obj3Desc); err != nil {
		t.Error(err)
	}
	if err = g.AddEdge(obj3Desc, obj4Desc); err != nil {
		t.Error(err)
	}

	// Validate vertexList length.
	if c := len(g.vertexList); c != 4 {
		t.Error("len(g.vertexList) expected 4, got:", c)
	}

	// Validate each outEdgeList for expected length
	if c := len(g.vertexList[obj1Desc].edges); c != 2 {
		fmt.Println(g.vertexList[obj1Desc].edges)
		t.Error("len(g.outEdgeList[obj1Desc].edges) expected 2, got: ", c)
	}
	if c := len(g.vertexList[obj2Desc].edges); c != 0 {
		fmt.Println(g.vertexList[obj2Desc].edges)
		t.Error("len(g.outEdgeList[obj2Desc].edges) expected 0, got: ", c)
	}
	if c := len(g.vertexList[obj3Desc].edges); c != 1 {
		fmt.Println(g.vertexList[obj3Desc].edges)
		t.Error("len(g.outEdgeList[obj3Desc].edges) expected 1, got: ", c)
	}
	if c := len(g.vertexList[obj4Desc].edges); c != 0 {
		fmt.Println(g.vertexList[obj4Desc].edges)
		t.Error("len(g.outEdgeList[obj4Desc].edges) expected 0, got: ", c)
	}
}

func TestGraph_DepthFirstSearch(t *testing.T) {
	g := NewGraph()

	expectedValue := []VertexDescriptor{1, 2, 4, 5, 7, 6, 3}

	var err error
	var obj1Desc VertexDescriptor
	var obj2Desc VertexDescriptor
	var obj3Desc VertexDescriptor
	var obj4Desc VertexDescriptor
	var obj5Desc VertexDescriptor
	var obj6Desc VertexDescriptor
	var obj7Desc VertexDescriptor

	// Create some objects and add them to the graph.
	if obj1Desc, err = g.AddVertex(newObject(1)); err != nil {
		t.Error(err)
	}
	if obj2Desc, err = g.AddVertex(newObject(2)); err != nil {
		t.Error(err)
	}
	if obj3Desc, err = g.AddVertex(newObject(3)); err != nil {
		t.Error(err)
	}
	if obj4Desc, err = g.AddVertex(newObject(4)); err != nil {
		t.Error(err)
	}
	if obj5Desc, err = g.AddVertex(newObject(5)); err != nil {
		t.Error(err)
	}
	if obj6Desc, err = g.AddVertex(newObject(6)); err != nil {
		t.Error(err)
	}
	if obj7Desc, err = g.AddVertex(newObject(7)); err != nil {
		t.Error(err)
	}

	// Add some edges.
	if err := g.AddEdge(obj1Desc, obj2Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj1Desc, obj3Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj2Desc, obj4Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj2Desc, obj5Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj2Desc, obj6Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj5Desc, obj7Desc); err != nil {
		t.Error(err)
	}

	dfs := g.DepthFirstSearch(obj1Desc, false)

	if !reflect.DeepEqual(dfs, expectedValue) {
		for _, val := range dfs {
			fmt.Printf("%d ", val)
		}
		t.Error("DepthFirstSearch result does not equal expected result")
	}
}

func TestGraph_BreadthFirstSearch(t *testing.T) {
	g := NewGraph()

	expectedValue := []VertexDescriptor{1, 2, 3, 4, 5, 6, 7}

	var err error
	var obj1Desc VertexDescriptor
	var obj2Desc VertexDescriptor
	var obj3Desc VertexDescriptor
	var obj4Desc VertexDescriptor
	var obj5Desc VertexDescriptor
	var obj6Desc VertexDescriptor
	var obj7Desc VertexDescriptor

	// Create some objects and add them to the graph.
	if obj1Desc, err = g.AddVertex(newObject(1)); err != nil {
		t.Error(err)
	}
	if obj2Desc, err = g.AddVertex(newObject(2)); err != nil {
		t.Error(err)
	}
	if obj3Desc, err = g.AddVertex(newObject(3)); err != nil {
		t.Error(err)
	}
	if obj4Desc, err = g.AddVertex(newObject(4)); err != nil {
		t.Error(err)
	}
	if obj5Desc, err = g.AddVertex(newObject(5)); err != nil {
		t.Error(err)
	}
	if obj6Desc, err = g.AddVertex(newObject(6)); err != nil {
		t.Error(err)
	}
	if obj7Desc, err = g.AddVertex(newObject(7)); err != nil {
		t.Error(err)
	}

	// Add some edges.
	if err := g.AddEdge(obj1Desc, obj2Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj1Desc, obj3Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj2Desc, obj4Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj2Desc, obj5Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj2Desc, obj6Desc); err != nil {
		t.Error(err)
	}
	if err := g.AddEdge(obj5Desc, obj7Desc); err != nil {
		t.Error(err)
	}

	bfs := g.BreadthFirstSearch(obj1Desc, true)

	if !reflect.DeepEqual(bfs, expectedValue) {
		for _, val := range bfs {
			fmt.Printf("%d ", val)
		}
		t.Error("DepthFirstSearch result does not equal expected result")
	}
}
