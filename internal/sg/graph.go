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
	"errors"
	"fmt"
	"math"
	"sync"
)

type SceneGraphListener interface {
	OnSceneGraphUpdate()
}

type Graph struct {
	vertexList     map[VertexDescriptor]*Vertex
	nextDescriptor VertexDescriptor
	mutex          *sync.Mutex
}

// Common functions

func NewGraph() *Graph {
	g := &Graph{
		vertexList: make(map[VertexDescriptor]*Vertex),
		mutex:      &sync.Mutex{},
	}

	return g
}

// Vertex Operations

func (g *Graph) VertexExistsWithDescriptor(vert VertexDescriptor) bool {
	_, ok := g.vertexList[vert]

	return ok
}

func (g *Graph) VertexExistsWithId(id uint32) bool {
	for idx := range g.vertexList {
		if g.vertexList[idx].data.ID() == id {
			return true
		}
	}

	return false
}

func (g *Graph) RemoveVertex(u VertexDescriptor) error {
	if !g.VertexExistsWithDescriptor(u) {
		return fmt.Errorf("remove vertex: descriptor %d does not exist", u)
	}

	// If this vertex has a parent, remove the reference.
	parent, err := g.Parent(u)
	if err == nil {
		g.vertexList[parent].edges = removeVertexDescriptorElement(g.vertexList[parent].edges, u)
	}

	// Get list of descendant descriptors.
	descendants := g.DepthFirstSearch(u, true)

	// Remove descendants references
	for idx := range descendants {
		delete(g.vertexList, descendants[idx])
	}

	// Remove this vertex
	delete(g.vertexList, u)

	return nil
}

func (g *Graph) MoveVertex(vert, parent VertexDescriptor) error {
	if !g.VertexExistsWithDescriptor(vert) {
		return fmt.Errorf("move vertex: target descriptor %d does not exist", vert)
	}

	if !g.VertexExistsWithDescriptor(parent) {
		return fmt.Errorf("move vertex: parent descriptor %d does not exist", vert)
	}

	if !g.DescendantOf(parent, vert) {
		return fmt.Errorf("move vertex: parent descriptor %d is a descendant of %d", parent, vert)
	}

	oldParent, err := g.Parent(vert)
	if err != nil {
		return fmt.Errorf("move vertex: invalid move, descriptor %d is orphaned or is root node", vert)
	}

	// Remove existing edge.
	// TODO: Check this for correctness.
	for idx := range g.vertexList[oldParent].edges {
		if g.vertexList[VertexDescriptor(idx)].descriptor == vert {
			g.vertexList[oldParent].edges = removeVertexDescriptorElement(g.vertexList[oldParent].edges, vert)
			break
		}
	}

	// Add new edge.
	g.AddEdge(parent, vert)

	return nil
}

func (g *Graph) AddVertex(object VertexNode) (VertexDescriptor, error) {
	v := &Vertex{
		edges: []VertexDescriptor{},
		data:  object,
	}
	var err error
	v.descriptor, err = g.NextDescriptor()
	if err != nil {
		return 0, err
	}

	g.vertexList[v.descriptor] = v

	return v.descriptor, nil
}

func (g *Graph) GetVertexByObject(object VertexNode) (VertexDescriptor, error) {
	if !g.VertexExistsWithId(object.ID()) {
		return g.AddVertex(object)
	}

	return g.GetVertexById(object.ID())
}

func (g *Graph) GetVertexById(id uint32) (VertexDescriptor, error) {
	for idx := range g.vertexList {
		if g.vertexList[idx].data.ID() == id {
			return g.vertexList[idx].descriptor, nil
		}
	}

	return 0, fmt.Errorf("get vertex: no such vertex with id: %d", id)
}

func (g *Graph) GetObjectAtVertex(v VertexDescriptor) VertexNode {
	if _, ok := g.vertexList[v]; ok {
		return g.vertexList[v].data
	}

	return nil
}

func (g *Graph) NextDescriptor() (VertexDescriptor, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if len(g.vertexList) >= math.MaxUint32 {
		return 0, errors.New("next descriptor: exceeded maximum number of vertex descriptors")
	}

	id := g.nextDescriptor + 1
	_, ok := g.vertexList[id]

	for ok {
		id = g.nextDescriptor + 1
		_, ok = g.vertexList[id]
	}

	g.nextDescriptor = id

	return g.nextDescriptor, nil
}

// Edge Operations

func (g *Graph) AddEdge(u, v VertexDescriptor) error {
	if g.DescendantOf(u, v) {
		return fmt.Errorf("add edge: %d is a descendant of %d", u, v)
	}

	g.vertexList[u].edges = append(g.vertexList[u].edges, v)

	return nil
}

func (g *Graph) EdgeExists(e Edge) bool {
	return g.EdgeExistsUV(e[0], e[1])
}

func (g *Graph) EdgeExistsUV(u, v VertexDescriptor) bool {
	return g.vertexDescriptorInOutEdgeList(u, v)
}

func (g *Graph) RemoveEdge(edge Edge) error {
	if _, ok := g.vertexList[edge.U()]; !ok {
		return fmt.Errorf("remove edge: descriptor %d not found in outEdgeList", edge.U())
	}

	for idx := range g.vertexList[edge.U()].edges {
		if g.vertexList[edge.U()].edges[idx] == edge.V() {
			g.vertexList[edge.U()].edges = deleteVertexDescriptorElement(g.vertexList[edge.U()].edges, idx)
			return nil
		}
	}

	return fmt.Errorf("remove edge: descriptor %v not found in outEdgeList[%d]", edge.V(), edge.U())
}

// Utility Functions

func (g *Graph) ParentOf(parent, descendant VertexDescriptor) bool {
	if g.VertexExistsWithDescriptor(parent) {
		return false
	}
	if g.VertexExistsWithDescriptor(descendant) {
		return false
	}

	return g.vertexDescriptorInOutEdgeList(parent, descendant)
}

func (g *Graph) DescendantOf(descendant, parent VertexDescriptor) bool {
	descendants := g.DepthFirstSearch(parent, true)

	for idx := range descendants {
		if descendants[idx] == descendant {
			return true
		}
	}

	return false
}

func (g *Graph) Parent(vertex VertexDescriptor) (VertexDescriptor, error) {
	if vertex != 0 {
		for edge := range g.vertexList {
			if edge == vertex {
				continue
			}

			if g.vertexDescriptorInOutEdgeList(edge, vertex) {
				return edge, nil
			}
		}
	}

	v := g.getParent(vertex)
	if v == 0 {
		return 0, fmt.Errorf("parent: vertex %d has no parent", vertex)

	}

	return v, nil
}

// Search Functions

func (g *Graph) DepthFirstSearch(u VertexDescriptor, includeDisabled bool) []VertexDescriptor {
	nodeList := make([]VertexDescriptor, 0)

	if !g.VertexExistsWithDescriptor(u) {
		return nodeList
	}

	if !includeDisabled {
		if !g.GetObjectAtVertex(u).Active() {
			return nodeList
		}
	}

	nodeList = append(nodeList, u)

	for idx := range g.vertexList[u].edges {
		newList := g.DepthFirstSearch(g.vertexList[u].edges[idx], includeDisabled)
		nodeList = append(nodeList, newList...)
	}

	return nodeList
}

func (g *Graph) BreadthFirstSearch(u VertexDescriptor, includeDisabled bool) []VertexDescriptor {
	nodeList := []VertexDescriptor{u}

	return append(nodeList, g.bFS(u, includeDisabled)...)
}

func (g *Graph) bFS(u VertexDescriptor, includeDisabled bool) []VertexDescriptor {
	nodeList := []VertexDescriptor{}

	if !g.VertexExistsWithDescriptor(u) {
		return nodeList
	}

	if !includeDisabled {
		if !g.GetObjectAtVertex(u).Active() {
			return nodeList
		}
	}

	nodeList = append(nodeList, g.vertexList[u].edges...)

	for idx := range g.vertexList[u].edges {
		newList := g.bFS(g.vertexList[u].edges[idx], includeDisabled)
		nodeList = append(nodeList, newList...)
	}

	return nodeList
}

func (g *Graph) ChildrenOf(u VertexDescriptor) []VertexDescriptor {
	if _, ok := g.vertexList[u]; !ok {
		return make([]VertexDescriptor, 0)
	}

	return g.vertexList[u].edges
}

// Utility Functions

func removeVertexDescriptorElement(s []VertexDescriptor, desc VertexDescriptor) []VertexDescriptor {
	for idx := range s {
		if s[idx] == desc {
			return deleteVertexDescriptorElement(s, idx)
		}
	}

	return s
}

func deleteVertexDescriptorElement(s []VertexDescriptor, idx int) []VertexDescriptor {
	if idx < 0 || idx >= len(s) {
		return s
	}

	if idx != len(s)-1 {
		s[idx] = s[len(s)-1]
	}

	return s[:len(s)-1]
}

func (g *Graph) vertexDescriptorInOutEdgeList(u, v VertexDescriptor) bool {
	if _, ok := g.vertexList[u]; !ok {
		return false
	}

	for idx := range g.vertexList[u].edges {
		if g.vertexList[u].edges[idx] == v {
			return true
		}
	}

	return false
}

func (g *Graph) getParent(u VertexDescriptor) VertexDescriptor {
	if _, ok := g.vertexList[u]; !ok {
		return 0
	}

	for i := range g.vertexList {
		if i == u {
			continue
		}

		for j := range g.vertexList[i].edges {
			if g.vertexList[i].edges[j] == u {
				return g.vertexList[i].descriptor
			}
		}
	}

	return 0
}
