package particle

import (
	"github.com/go-gl/gl/v4.3-core/gl"

	"github.com/haakenlabs/forge/internal/engine"
)

type ModuleRenderer struct {
	system       *System
	renderShader *engine.Shader
	sprite       *engine.Texture2D
}

func (m *ModuleRenderer) Render(camera *engine.Camera) {
	m.system.Simulate()

	m.renderShader.Bind()
	m.system.Core.particleBuffer.Bind()

	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE)

	m.renderShader.SetUniform("v_model_matrix", m.system.GetTransform().ActiveMatrix())
	m.renderShader.SetUniform("v_view_matrix", camera.ViewMatrix())
	m.renderShader.SetUniform("v_projection_matrix", camera.ProjectionMatrix())
	m.renderShader.SetUniform("v_offset", m.system.inOffset)

	m.sprite.ActivateTexture(gl.TEXTURE0)
	gl.DrawArrays(gl.POINTS, 0, int32(m.system.Core.alive))

	m.system.Core.particleBuffer.Unbind()
	m.renderShader.Unbind()

	gl.Disable(gl.BLEND)
	gl.Enable(gl.DEPTH_TEST)
}

func NewModuleRenderer(system *System) *ModuleRenderer {
	m := &ModuleRenderer{
		system: system,
	}

	m.renderShader = engine.GetAsset().MustGet(engine.AssetNameShader, "particle/render").(*engine.Shader)
	m.sprite = engine.GetAsset().MustGet(engine.AssetNameImage, "particle.png").(*engine.Texture2D)

	return m
}
