package particle

import (
	"git.dbnservers.net/haakenlabs/forge/internal/engine"
)

type ModuleRenderer struct {
}

func (m *ModuleRenderer) Render(camera *engine.Camera) {
}

func NewModuleRenderer() *ModuleRenderer {
	m := &ModuleRenderer{}

	return m
}
