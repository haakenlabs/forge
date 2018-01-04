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
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"path/filepath"
	"sync"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"

	"github.com/haakenlabs/forge/internal/image/hdr"
	"github.com/haakenlabs/forge/internal/math"

	_ "image/jpeg"
	_ "image/png"
)

const (
	AssetNameSkybox = "skybox"
)

var rotMatrices = [6]mgl32.Mat4{
	// X
	mgl32.LookAtV(mgl32.Vec3{}, mgl32.Vec3{1, 0, 0}, mgl32.Vec3{0, 1, 0}),
	mgl32.LookAtV(mgl32.Vec3{}, mgl32.Vec3{-1, 0, 0}, mgl32.Vec3{0, 1, 0}),
	// Y
	mgl32.LookAtV(mgl32.Vec3{}, mgl32.Vec3{0, 1, 0}, mgl32.Vec3{0, 0, -1}),
	mgl32.LookAtV(mgl32.Vec3{}, mgl32.Vec3{0, -1, 0}, mgl32.Vec3{0, 0, 1}),
	// Z
	mgl32.LookAtV(mgl32.Vec3{}, mgl32.Vec3{0, 0, 1}, mgl32.Vec3{0, 1, 0}),
	mgl32.LookAtV(mgl32.Vec3{}, mgl32.Vec3{0, 0, -1}, mgl32.Vec3{0, 1, 0}),
}

var _ AssetHandler = &SkyboxHandler{}

type SkyboxMetadata struct {
	Name       string `json:"name"`
	Radiance   string `json:"radiance"`
	Specular   string `json:"specular"`
	Irradiance string `json:"irradiance"`
}

type SkyboxHandler struct {
	BaseAssetHandler
}

func NewSkyboxHandler() *SkyboxHandler {
	h := &SkyboxHandler{}
	h.Items = make(map[string]uint32)
	h.Mu = &sync.RWMutex{}

	return h
}

func (h *SkyboxHandler) Name() string {
	return AssetNameSkybox
}

func (h *SkyboxHandler) Load(r *Resource) error {
	m := &SkyboxMetadata{}

	if err := json.Unmarshal(r.Bytes(), m); err != nil {
		return err
	}

	if _, dup := h.Items[m.Name]; dup {
		return ErrAssetExists(m.Name)
	}

	skybox, err := h.loadMap(m, r.DirPrefix())
	if err != nil {
		return err
	}

	h.Items[m.Name] = skybox.ID()

	return nil
}

func (h *SkyboxHandler) loadMap(m *SkyboxMetadata, dir string) (skybox *Skybox, err error) {
	skybox = NewSkybox(nil, nil, nil)

	var specR, irrdR *Resource
	var specTex, irrdTex *Texture2D

	genSpecular := len(m.Specular) == 0
	genIrradiance := len(m.Irradiance) == 0

	radiR, err := NewResource(filepath.Join(dir, m.Radiance))
	if err != nil {
		return nil, err
	}
	if err := GetAsset().ReadResource(radiR); err != nil {
		return nil, err
	}

	if !genSpecular {
		specR, err = NewResource(filepath.Join(dir, m.Specular))
		if err != nil {
			return nil, err
		}
		if err := GetAsset().ReadResource(specR); err != nil {
			return nil, err
		}
	}
	if !genIrradiance {
		irrdR, err = NewResource(filepath.Join(dir, m.Irradiance))
		if err != nil {
			return nil, err
		}
		if err := GetAsset().ReadResource(irrdR); err != nil {
			return nil, err
		}
	}

	rimg := loadImage(radiR)
	simg := loadImage(specR)
	iimg := loadImage(irrdR)

	radiTex, err := loadTexture(rimg)
	if err != nil {
		return nil, err
	}
	if !genSpecular {
		specTex, err = loadTexture(simg)
		if err != nil {
			return nil, err
		}
	}
	if !genIrradiance {
		irrdTex, err = loadTexture(iimg)
		if err != nil {
			return nil, err
		}
	}

	fbo := NewFramebufferRaw()
	defer fbo.Dealloc()

	skybox.radiance, err = makeCubemap(radiTex, fbo, radiTex.Size().Y()/2)
	if err != nil {
		return nil, err
	}

	if genSpecular {
		skybox.specular, err = generateSpecular(skybox.radiance, fbo)
	} else {
		skybox.specular, err = makeCubemap(specTex, fbo, specTex.Size().Y()/2)
	}
	if err != nil {
		return nil, err
	}

	if genIrradiance {
		skybox.irradiance, err = generateIrradiance(skybox.radiance, fbo)
	} else {
		skybox.irradiance, err = makeCubemap(irrdTex, fbo, irrdTex.Size().Y()/2)
	}
	if err != nil {
		return nil, err
	}

	return skybox, nil
}

// Get gets an asset by name.
func (h *SkyboxHandler) Get(name string) (*Skybox, error) {
	a, err := h.GetAsset(name)
	if err != nil {
		return nil, err
	}

	a2, ok := a.(*Skybox)
	if !ok {
		return nil, ErrAssetType(name)
	}

	return a2, nil
}

// MustGet is like GetAsset, but panics if an error occurs.
func (h *SkyboxHandler) MustGet(name string) *Skybox {
	a, err := h.Get(name)
	if err != nil {
		panic(err)
	}

	return a
}

func loadImage(r *Resource) (img image.Image) {
	img, _, err := image.Decode(r.Reader())
	if err != nil {
		logrus.Error(err)
		return nil
	}

	return img
}

func loadTexture(img image.Image) (tex *Texture2D, err error) {
	x := int32(img.Bounds().Dx())
	y := int32(img.Bounds().Dy())

	tex = NewTexture2D(math.IVec2{x, y}, TextureFormatDefaultColor)

	switch img.ColorModel() {
	case color.CMYKModel:
		panic("color.CMYKModel unsupported")
	case color.YCbCrModel:
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		tex.SetTexFormat(TextureFormatRGBA8)
		tex.SetData(rgba.Pix)
	case color.RGBA64Model:
		rgba := image.NewRGBA64(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		tex.SetTexFormat(TextureFormatRGBA16)
		tex.SetData(rgba.Pix)
	case color.RGBAModel:
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		tex.SetTexFormat(TextureFormatRGBA8)
		tex.SetData(rgba.Pix)
	case color.NRGBA64Model:
		rgba := image.NewNRGBA64(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		tex.SetTexFormat(TextureFormatRGBA16)
		tex.SetData(rgba.Pix)
	case color.NRGBAModel:
		rgba := image.NewNRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		tex.SetTexFormat(TextureFormatRGBA8)
		tex.SetData(rgba.Pix)
	case hdr.RGB96Model:
		rgba := hdr.NewRGB96(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		tex.SetTexFormat(TextureFormatRGB32)

		var data []float32
		for y := 0; y < rgba.Rect.Dy(); y++ {
			for x := 0; x < rgba.Rect.Dx(); x++ {
				c := rgba.RGB96At(x, y)
				data = append(data, c.R)
				data = append(data, c.G)
				data = append(data, c.B)
			}
		}
		tex.SetHDRData(data)

	default:
		return nil, errors.New("unknown image type")
	}

	if err := tex.Alloc(); err != nil {
		return nil, err
	}

	return tex, err
}

func makeCubemap(tex *Texture2D, fbo *Framebuffer, faceSize int32) (cubemap *TextureCubemap, err error) {
	fbo.Bind()
	fbo.SetSize(math.IVec2{faceSize, faceSize})

	cubemap = NewTextureCubemap(math.IVec2{faceSize, faceSize}, tex.TexFormat())
	if err := cubemap.Alloc(); err != nil {
		return nil, err
	}

	mesh := NewMeshQuadBack()
	mesh.Bind()

	gl.Disable(gl.DEPTH_TEST)
	gl.DepthMask(false)

	shader := GetAsset().MustGet(AssetNameShader, "utils/cubeconv").(*Shader)
	shader.Bind()
	shader.SetUniform("v_projection_matrix", mgl32.Perspective(math.Pi32/2.0, 1.0, 0.1, 2.0))
	tex.ActivateTexture(gl.TEXTURE0)

	fbo.ClearBuffers()

	for i := uint32(0); i < 6; i++ {
		shader.SetUniform("v_view_matrix", rotMatrices[i])
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_CUBE_MAP_POSITIVE_X+i, cubemap.Reference(), 0)
		mesh.Draw()
	}

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, 0, 0)

	gl.DepthMask(true)
	gl.Enable(gl.DEPTH_TEST)

	shader.Unbind()
	mesh.Unbind()
	fbo.Unbind()

	return
}

func generateSpecular(radiance *TextureCubemap, fbo *Framebuffer) (spec *TextureCubemap, err error) {
	//return nil, ErrNotImplemented

	return nil, nil
}

func generateIrradiance(radiance *TextureCubemap, fbo *Framebuffer) (irrd *TextureCubemap, err error) {
	//return nil, ErrNotImplemented

	return nil, nil
}
