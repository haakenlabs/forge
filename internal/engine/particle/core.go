package particle

import "github.com/haakenlabs/forge/internal/engine"

type ModuleCore struct {
	StartColor    engine.Color
	Duration      float64
	StartDelay    float64
	StartLifetime float64
	StartSpeed    float64
	StartSize     float64
	MaxParticles  uint32
	Looping       bool
}

func NewModuleCore() *ModuleCore {
	m := &ModuleCore{
		StartColor:   engine.ColorWhite(),
		StartSpeed:   10.0,
		StartSize:    1.0,
		MaxParticles: 10000,
		Looping:      true,
	}

	return m
}
