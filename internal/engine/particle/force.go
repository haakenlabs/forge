package particle

import "github.com/go-gl/mathgl/mgl32"

type AttractorMode uint32

const (
	AttractorNormal AttractorMode = iota
	AttractorRepell
	AttractorBlackhole
	AttractorGlobal
)

type ModuleForce struct {
	system     *System
	attractors []*Attractor

	EnableAttractors bool
}

type Attractor struct {
	Position  mgl32.Vec4
	Direction mgl32.Vec4
	Mode      AttractorMode
	Force     float32
	Range     float32
	Unused    float32
}

func (m *ModuleForce) AddAttractor(a *Attractor) {
	m.attractors = append(m.attractors, a)
}

func NewModuleForce() *ModuleForce {
	m := &ModuleForce{}

	return m
}

func NewAttractor() *Attractor {
	return &Attractor{}
}
