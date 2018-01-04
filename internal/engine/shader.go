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

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

type ShaderComponent uint32

const (
	ShaderComponentVertex         ShaderComponent = gl.VERTEX_SHADER
	ShaderComponentGeometry                       = gl.GEOMETRY_SHADER
	ShaderComponentFragment                       = gl.FRAGMENT_SHADER
	ShaderComponentCompute                        = gl.COMPUTE_SHADER
	ShaderComponentTessControl                    = gl.TESS_CONTROL_SHADER
	ShaderComponentTessEvaluation                 = gl.TESS_EVALUATION_SHADER
)

var _ Object = &Shader{}

type Shader struct {
	BaseObject

	programId       uint32
	components      map[ShaderComponent]uint32
	data            []byte
	deferredCapable bool
}

func (s *Shader) Alloc() error {
	return s.Build()
}

// Dealloc releases builtin for this shader.
func (s *Shader) Dealloc() {
	if s.programId != 0 {
		for k := range s.components {
			destroyComponent(s.components[k], s.programId)
			delete(s.components, k)
		}

		gl.DeleteProgram(s.programId)

		s.programId = 0
	}
}

func (s *Shader) AddData(newData []byte) {
	s.data = append(s.data, newData...)
}

func (s *Shader) Build() error {
	// Create Program ID
	s.programId = gl.CreateProgram()

	if containsShaderType(ShaderComponentVertex, s.data) {
		componentId, err := loadComponent(s.programId, ShaderComponentVertex, s.data)
		if err != nil {
			return err
		}
		s.components[ShaderComponentVertex] = componentId
	}
	if containsShaderType(ShaderComponentGeometry, s.data) {
		componentId, err := loadComponent(s.programId, ShaderComponentGeometry, s.data)
		if err != nil {
			return err
		}
		s.components[ShaderComponentGeometry] = componentId
	}
	if containsShaderType(ShaderComponentFragment, s.data) {
		componentId, err := loadComponent(s.programId, ShaderComponentFragment, s.data)
		if err != nil {
			return err
		}
		s.components[ShaderComponentFragment] = componentId
	}
	if containsShaderType(ShaderComponentCompute, s.data) {
		componentId, err := loadComponent(s.programId, ShaderComponentCompute, s.data)
		if err != nil {
			return err
		}
		s.components[ShaderComponentCompute] = componentId
	}
	if containsShaderType(ShaderComponentTessControl, s.data) {
		componentId, err := loadComponent(s.programId, ShaderComponentTessControl, s.data)
		if err != nil {
			return err
		}
		s.components[ShaderComponentTessControl] = componentId
	}
	if containsShaderType(ShaderComponentTessEvaluation, s.data) {
		componentId, err := loadComponent(s.programId, ShaderComponentTessEvaluation, s.data)
		if err != nil {
			return err
		}
		s.components[ShaderComponentTessEvaluation] = componentId
	}

	// Set transform feedback varyings
	// TODO: Implement this

	// Validate and link
	err := Link(s.programId)

	return err
}

func (s *Shader) ProgramId() uint32 {
	return s.programId
}

func (s *Shader) Reference() uint32 {
	return s.programId
}

func (s *Shader) Bind() {
	BindShader(s.programId)
}

func (s *Shader) Unbind() {
	UnbindShader()
}

func (s *Shader) SetSubroutine(componentType ShaderComponent, subroutineName string) {
	idx := gl.GetSubroutineIndex(s.programId, uint32(componentType), gl.Str(subroutineName+"\x00"))
	gl.UniformSubroutinesuiv(uint32(componentType), 1, &idx)
}

func (s *Shader) SetUniform(uniformName string, value interface{}) {
	switch v := value.(type) {
	case bool:
		var val int32
		if v {
			val = 1
		}
		gl.Uniform1i(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), val)
	case int32:
		gl.Uniform1i(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), v)
	case float32:
		gl.Uniform1f(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), v)
	case uint32:
		gl.Uniform1ui(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), v)
	case mgl32.Vec2:
		gl.Uniform2fv(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), 1, &v[0])
	case mgl32.Vec3:
		gl.Uniform3fv(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), 1, &v[0])
	case mgl32.Vec4:
		gl.Uniform4fv(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), 1, &v[0])
	case mgl32.Mat2:
		gl.UniformMatrix2fv(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), 1, false, &v[0])
	case mgl32.Mat3:
		gl.UniformMatrix3fv(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), 1, false, &v[0])
	case mgl32.Mat4:
		gl.UniformMatrix4fv(gl.GetUniformLocation(s.programId, gl.Str(uniformName+"\x00")), 1, false, &v[0])
	}
}

func (s *Shader) DeferredCapable() bool {
	return s.deferredCapable
}

// Common Functions

func Link(programId uint32) error {
	gl.LinkProgram(programId)
	return ValidateProgram(programId)
}

func ValidateComponent(componentId uint32) error {
	var status int32
	gl.GetShaderiv(componentId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(componentId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(componentId, logLength, nil, gl.Str(log))

		return fmt.Errorf("shader %d compilation failed: %v", componentId, log)
	}

	return nil
}

func ValidateProgram(programId uint32) error {
	var status int32
	gl.GetProgramiv(programId, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programId, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programId, logLength, nil, gl.Str(log))

		return fmt.Errorf("program %d link failed: %v", programId, log)
	}

	return nil
}

func BindShader(programId uint32) {
	gl.UseProgram(programId)
}

func UnbindShader() {
	gl.UseProgram(0)
}

func destroyComponent(componentId uint32, programId uint32) {
	gl.DetachShader(programId, componentId)
	gl.DeleteShader(componentId)
}

func containsShaderType(shaderType ShaderComponent, data []byte) bool {
	switch shaderType {
	case ShaderComponentVertex:
		return bytes.Contains(data, []byte("#ifdef _VERTEX_"))
	case ShaderComponentGeometry:
		return bytes.Contains(data, []byte("#ifdef _GEOMETRY_"))
	case ShaderComponentFragment:
		return bytes.Contains(data, []byte("#ifdef _FRAGMENT_"))
	case ShaderComponentCompute:
		return bytes.Contains(data, []byte("#ifdef _COMPUTE_"))
	case ShaderComponentTessControl:
		return bytes.Contains(data, []byte("#ifdef _TESSCONTROL_"))
	case ShaderComponentTessEvaluation:
		return bytes.Contains(data, []byte("#ifdef _TESSEVAL_"))
	}

	return false
}

func loadComponent(programId uint32, componentType ShaderComponent, data []byte) (uint32, error) {
	header := []byte("#version 430\n")

	switch componentType {
	case ShaderComponentVertex:
		header = append(header, []byte("#define _VERTEX_\n")...)
	case ShaderComponentGeometry:
		header = append(header, []byte("#define _GEOMETRY_\n")...)
	case ShaderComponentFragment:
		header = append(header, []byte("#define _FRAGMENT_\n")...)
	case ShaderComponentCompute:
		header = append(header, []byte("#define _COMPUTE_\n")...)
	case ShaderComponentTessControl:
		header = append(header, []byte("#define _TESSCONTROL_\n")...)
	case ShaderComponentTessEvaluation:
		header = append(header, []byte("#define _TESSEVAL_\n")...)
	default:
		return 0, fmt.Errorf("loadComponent failed: unknown component type: %d", componentType)
	}

	data = append(header, data...)

	componentId := gl.CreateShader(uint32(componentType))

	csrc, free := gl.Strs(string(data))
	srcLength := int32(len(data))
	gl.ShaderSource(componentId, 1, csrc, &srcLength)
	free()
	gl.CompileShader(componentId)

	err := ValidateComponent(componentId)
	if err != nil {
		fmt.Println(string(data))
		return 0, err
	}

	gl.AttachShader(programId, componentId)

	logrus.Debugf("Loaded component(%s) %d for program %d", ShaderComponentToString(componentType), componentId, programId)

	return componentId, nil
}

// ShaderComponentToString returns the string representation of a core.ShaderComponent.
func ShaderComponentToString(component ShaderComponent) string {
	switch component {
	case ShaderComponentVertex:
		return "VERTEX"
	case ShaderComponentGeometry:
		return "GEOMETRY"
	case ShaderComponentFragment:
		return "FRAGMENT"
	case ShaderComponentCompute:
		return "COMPUTE"
	case ShaderComponentTessControl:
		return "TESSCONTROL"
	case ShaderComponentTessEvaluation:
		return "TESSEVAL"
	}

	return "INVALID"
}

func NewShader() *Shader {
	s := &Shader{
		components: make(map[ShaderComponent]uint32),
	}

	s.SetName("Shader")
	GetInstance().MustAssign(s)

	return s
}

func NewShaderUtilsCopy() *Shader {
	return GetAsset().MustGet(AssetNameShader, "utils/copy").(*Shader)
}

func NewShaderUtilsSkybox() *Shader {
	return GetAsset().MustGet(AssetNameShader, "utils/skybox").(*Shader)
}
