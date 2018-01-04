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
	"fmt"

	"github.com/go-gl/gl/v4.3-core/gl"

	"git.dbnservers.net/haakenlabs/forge/internal/math"
	"unsafe"
)

type TextureFormat uint32

const (
	TextureFormatDefaultColor TextureFormat = iota
	TextureFormatDefaultHDRColor
	TextureFormatDefaultDepth
	TextureFormatR8
	TextureFormatRG8
	TextureFormatRGB8
	TextureFormatRGBA8
	TextureFormatR16
	TextureFormatRG16
	TextureFormatRGB16
	TextureFormatRGBA16
	TextureFormatRGBA16UI
	TextureFormatR32
	TextureFormatRG32
	TextureFormatRGB32
	TextureFormatRGBA32
	TextureFormatRGB32UI
	TextureFormatRGBA32UI
	TextureFormatDepth16
	TextureFormatDepth24
	TextureFormatDepth24Stencil8
	TextureFormatStencil8
)

type Texture interface {
	Object

	ActivateTexture(textureUnit uint32)
	Bind()
	FilterMag() int32
	FilterMin() int32
	GenerateMipmaps()
	GLFormat() uint32
	GLInternalFormat() int32
	GLStorageFormat() uint32
	GLTextureType() uint32
	Height() int32
	Layers() int32
	MipLevels() uint32
	Resizable() bool
	SetFilter(magFilter, minFilter int32)
	SetGLFormats(internalFormat int32, format uint32, storageFormat uint32)
	SetLayers(int32)
	SetMagFilter(magFilter int32)
	SetMinFilter(minFilter int32)
	SetResizable(resizable bool)
	SetSize(size math.IVec2) error
	SetTexFormat(format TextureFormat)
	SetWrapR(wrapR int32)
	SetWrapRST(wrapR, wrapS, wrapT int32)
	SetWrapS(wrapS int32)
	SetWrapST(wrapS, wrapT int32)
	SetWrapT(wrapT int32)
	Size() math.IVec2
	TexFormat() TextureFormat
	Reference() uint32
	Unbind()
	Upload()
	Width() int32
	WrapR() int32
	WrapS() int32
	WrapT() int32
}

type UploadFunc func()

type BaseTexture struct {
	BaseObject

	uploadFunc     UploadFunc
	internalFormat int32
	storageFormat  uint32
	glFormat       uint32
	filterMag      int32
	filterMin      int32
	wrapR          int32
	wrapS          int32
	wrapT          int32
	layers         int32
	reference      uint32
	textureFormat  TextureFormat
	size           math.IVec2
	resizable      bool
	textureType    uint32
}

func TextureFormatToInternal(format TextureFormat) int32 {
	switch format {
	case TextureFormatR8:
		return gl.R8
	case TextureFormatRG8:
		return gl.RG8
	case TextureFormatRGB8:
		return gl.RGB8
	case TextureFormatDefaultColor:
		fallthrough
	case TextureFormatRGBA8:
		return gl.RGBA8
	case TextureFormatR16:
		return gl.R16F
	case TextureFormatRG16:
		return gl.RG16F
	case TextureFormatRGB16:
		return gl.RGB16F
	case TextureFormatDefaultHDRColor:
		fallthrough
	case TextureFormatRGBA16:
		return gl.RGBA16F
	case TextureFormatR32:
		return gl.R32F
	case TextureFormatRG32:
		return gl.RG32F
	case TextureFormatRGB32:
		return gl.RGB32F
	case TextureFormatRGBA32:
		return gl.RGBA32F
	case TextureFormatRGB32UI:
		return gl.RGB32UI
	case TextureFormatRGBA32UI:
		return gl.RGBA32UI
	case TextureFormatDepth16:
		return gl.DEPTH_COMPONENT16
	case TextureFormatDefaultDepth:
		fallthrough
	case TextureFormatDepth24:
		return gl.DEPTH_COMPONENT24
	case TextureFormatDepth24Stencil8:
		return gl.DEPTH24_STENCIL8
	case TextureFormatStencil8:
		return gl.STENCIL_INDEX8
	case TextureFormatRGBA16UI:
		return gl.RGBA16UI
	}

	return 0
}

func TextureFormatToFormat(format TextureFormat) uint32 {
	switch format {
	case TextureFormatR8:
		fallthrough
	case TextureFormatR16:
		fallthrough
	case TextureFormatR32:
		return gl.RED
	case TextureFormatRG8:
		fallthrough
	case TextureFormatRG16:
		fallthrough
	case TextureFormatRG32:
		return gl.RG
	case TextureFormatRGB8:
		fallthrough
	case TextureFormatRGB16:
		fallthrough
	case TextureFormatRGB32:
		return gl.RGB
	case TextureFormatRGB32UI:
		return gl.RGB_INTEGER
	case TextureFormatDefaultColor:
		fallthrough
	case TextureFormatRGBA8:
		fallthrough
	case TextureFormatDefaultHDRColor:
		fallthrough
	case TextureFormatRGBA16:
		fallthrough
	case TextureFormatRGBA16UI:
		fallthrough
	case TextureFormatRGBA32:
		return gl.RGBA
	case TextureFormatRGBA32UI:
		return gl.RGBA_INTEGER
	case TextureFormatDefaultDepth:
		fallthrough
	case TextureFormatDepth16:
		fallthrough
	case TextureFormatDepth24:
		return gl.DEPTH_COMPONENT
	case TextureFormatDepth24Stencil8:
		fallthrough
	case TextureFormatStencil8:
		return 0
	}

	return 0
}

func TextureFormatToStorage(format TextureFormat) uint32 {
	switch format {
	case TextureFormatDefaultColor:
		fallthrough
	case TextureFormatR8:
		fallthrough
	case TextureFormatRG8:
		fallthrough
	case TextureFormatRGB8:
		fallthrough
	case TextureFormatRGBA8:
		fallthrough
	case TextureFormatStencil8:
		return gl.UNSIGNED_BYTE
	case TextureFormatR16:
		fallthrough
	case TextureFormatRG16:
		fallthrough
	case TextureFormatRGB16:
		fallthrough
	case TextureFormatDefaultHDRColor:
		fallthrough
	case TextureFormatRGBA16:
		return gl.HALF_FLOAT
	case TextureFormatRGBA16UI:
		return gl.UNSIGNED_SHORT
	case TextureFormatR32:
		fallthrough
	case TextureFormatRG32:
		fallthrough
	case TextureFormatRGB32:
		fallthrough
	case TextureFormatRGBA32:
		return gl.FLOAT
	case TextureFormatRGB32UI:
		fallthrough
	case TextureFormatRGBA32UI:
		return gl.UNSIGNED_INT
	case TextureFormatDefaultDepth:
		fallthrough
	case TextureFormatDepth16:
		fallthrough
	case TextureFormatDepth24:
		return gl.FLOAT
	case TextureFormatDepth24Stencil8:
		return gl.UNSIGNED_INT_24_8
	}

	return 0
}

func (t *BaseTexture) Alloc() error {
	if t.reference != 0 {
		return nil
	}

	gl.GenTextures(1, &t.reference)

	t.filterMag = gl.LINEAR
	t.filterMin = gl.LINEAR
	t.wrapR = gl.CLAMP_TO_EDGE
	t.wrapS = gl.CLAMP_TO_EDGE
	t.wrapT = gl.CLAMP_TO_EDGE
	t.resizable = true
	t.layers = 1

	t.uploadFunc()

	t.SetFilter(t.filterMag, t.filterMin)
	t.SetWrapRST(t.wrapR, t.wrapS, t.wrapT)

	return nil
}

// Release
func (t *BaseTexture) Dealloc() {
	if t.reference != 0 {
		gl.DeleteTextures(1, &t.reference)
		t.reference = 0
	}
}

func (t *BaseTexture) TextureType() uint32 {
	return t.textureType
}

// ActivateTexture
func (t *BaseTexture) ActivateTexture(textureUnit uint32) {
	gl.ActiveTexture(textureUnit)
	t.Bind()
}

// Bind
func (t *BaseTexture) Bind() {
	gl.BindTexture(t.textureType, t.reference)
}

// FilterMag
func (t *BaseTexture) FilterMag() int32 {
	return t.filterMag
}

// FilterMin
func (t *BaseTexture) FilterMin() int32 {
	return t.filterMin
}

// GenerateMipmaps
func (t *BaseTexture) GenerateMipmaps() {

}

// GLFormat
func (t *BaseTexture) GLFormat() uint32 {
	return t.glFormat
}

// GLInternalFormat
func (t *BaseTexture) GLInternalFormat() int32 {
	return t.internalFormat
}

// GLStorageFormat
func (t *BaseTexture) GLStorageFormat() uint32 {
	return t.storageFormat
}

// GLTextureType
func (t *BaseTexture) GLTextureType() uint32 {
	return t.textureType
}

// Height
func (t *BaseTexture) Height() int32 {
	return t.size.Y()
}

// Layers
func (t *BaseTexture) Layers() int32 {
	return t.layers
}

func (t *BaseTexture) SetLayers(layers int32) {
	t.layers = layers
}

// MipLevels
func (t *BaseTexture) MipLevels() uint32 {
	return 1
}

// Resizable
func (t *BaseTexture) Resizable() bool {
	return t.resizable
}

// SetFilter
func (t *BaseTexture) SetFilter(magFilter, minFilter int32) {
	t.SetMagFilter(magFilter)
	t.SetMinFilter(minFilter)
}

// SetGLFormats
func (t *BaseTexture) SetGLFormats(internalFormat int32, format uint32, storageFormat uint32) {
	t.internalFormat = internalFormat
	t.glFormat = format
	t.storageFormat = storageFormat

	//t.Upload()
	t.uploadFunc()
}

// SetMagFilter
func (t *BaseTexture) SetMagFilter(magFilter int32) {
	t.filterMag = magFilter
	gl.TexParameteri(t.textureType, gl.TEXTURE_MAG_FILTER, t.filterMag)
}

// SetMinFilter
func (t *BaseTexture) SetMinFilter(minFilter int32) {
	t.filterMin = minFilter
	gl.TexParameteri(t.textureType, gl.TEXTURE_MIN_FILTER, t.filterMin)
}

// SetResizable
func (t *BaseTexture) SetResizable(resizable bool) {
	t.resizable = resizable
}

// SetSize
func (t *BaseTexture) SetSize(size math.IVec2) error {
	if !t.resizable {
		return fmt.Errorf("texture setSize error: texture %d is not resizable", t.reference)
	}
	if size.X() <= 0 || size.Y() <= 0 {
		return fmt.Errorf("texture setSize error: invalid size: %s", size)
	}

	t.size = size
	t.uploadFunc()

	return nil
}

// SetTexFormat
func (t *BaseTexture) SetTexFormat(format TextureFormat) {
	t.SetGLFormats(TextureFormatToInternal(format), TextureFormatToFormat(format), TextureFormatToStorage(format))
}

// SetWrapR
func (t *BaseTexture) SetWrapR(wrapR int32) {
	t.wrapR = wrapR
	gl.TexParameteri(t.textureType, gl.TEXTURE_WRAP_R, t.wrapR)
	if t.wrapR == gl.CLAMP_TO_BORDER {
		color := [4]float32{}
		gl.TexParameterfv(t.textureType, gl.TEXTURE_BORDER_COLOR, &color[0])
	}
}

// SetWrapRST
func (t *BaseTexture) SetWrapRST(wrapR, wrapS, wrapT int32) {
	t.SetWrapR(wrapR)
	t.SetWrapS(wrapS)
	t.SetWrapT(wrapT)
}

// SetWrapS
func (t *BaseTexture) SetWrapS(wrapS int32) {
	t.wrapS = wrapS
	gl.TexParameteri(t.textureType, gl.TEXTURE_WRAP_S, t.wrapS)
	if t.wrapS == gl.CLAMP_TO_BORDER {
		color := [4]float32{}
		gl.TexParameterfv(t.textureType, gl.TEXTURE_BORDER_COLOR, &color[0])
	}
}

// SetWrapST
func (t *BaseTexture) SetWrapST(wrapS, wrapT int32) {
	t.SetWrapS(wrapS)
	t.SetWrapT(wrapT)
}

// SetWrapT
func (t *BaseTexture) SetWrapT(wrapT int32) {
	t.wrapT = wrapT
	gl.TexParameteri(t.textureType, gl.TEXTURE_WRAP_T, t.wrapT)
	if t.wrapT == gl.CLAMP_TO_BORDER {
		color := [4]float32{}
		gl.TexParameterfv(t.textureType, gl.TEXTURE_BORDER_COLOR, &color[0])
	}
}

// Size
func (t *BaseTexture) Size() math.IVec2 {
	return t.size
}

// TexFormat
func (t *BaseTexture) TexFormat() TextureFormat {
	return t.textureFormat
}

// Reference
func (t *BaseTexture) Reference() uint32 {
	return t.reference
}

// Upload
func (t *BaseTexture) Upload() {
	panic("Unimplemented!")
}

// Unbind
func (t *BaseTexture) Unbind() {
	gl.BindTexture(t.textureType, 0)
}

// Width
func (t *BaseTexture) Width() int32 {
	return t.size.X()
}

// WrapR
func (t *BaseTexture) WrapR() int32 {
	return t.wrapR
}

// WrapS
func (t *BaseTexture) WrapS() int32 {
	return t.wrapS
}

// WrapT
func (t *BaseTexture) WrapT() int32 {
	return t.wrapT
}

type Texture2D struct {
	BaseTexture

	data    []uint8
	hdrData []float32
}

// Texture2D Methods

func NewTexture2D(size math.IVec2, format TextureFormat) *Texture2D {
	t := &Texture2D{}

	t.textureType = gl.TEXTURE_2D

	t.SetName("Texture2D")
	GetInstance().MustAssign(t)

	t.size = size
	t.uploadFunc = t.Upload

	t.internalFormat = TextureFormatToInternal(format)
	t.glFormat = TextureFormatToFormat(format)
	t.storageFormat = TextureFormatToStorage(format)

	return t
}

func NewTexture2DFrom(texture Texture2D) *Texture2D {
	t := &Texture2D{}

	t.textureType = gl.TEXTURE_2D

	t.SetName("Texture2D")
	GetInstance().MustAssign(t)

	t.size = texture.Size()
	t.uploadFunc = t.Upload

	t.internalFormat = texture.GLInternalFormat()
	t.glFormat = texture.GLFormat()
	t.storageFormat = texture.GLStorageFormat()

	return t
}

func (t *Texture2D) Upload() {
	t.Bind()

	var ptr unsafe.Pointer

	if t.hdrData != nil && len(t.hdrData) > 0 {
		ptr = gl.Ptr(t.hdrData)
	} else if len(t.data) > 0 {
		ptr = gl.Ptr(t.data)
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, t.internalFormat, t.size.X(), t.size.Y(), 0, t.glFormat, t.storageFormat, ptr)
}

func (t *Texture2D) SetData(data []uint8) {
	t.data = data
}

func (t *Texture2D) SetHDRData(data []float32) {
	t.hdrData = data
}

type Texture3D struct {
	BaseTexture
}

func NewTexture3D(size math.IVec2, layers int32, format TextureFormat) *Texture3D {
	t := &Texture3D{}

	t.textureType = gl.TEXTURE_3D

	t.SetName("Texture3D")
	GetInstance().MustAssign(t)

	t.size = size
	t.uploadFunc = t.Upload

	t.internalFormat = TextureFormatToInternal(format)
	t.glFormat = TextureFormatToFormat(format)
	t.storageFormat = TextureFormatToStorage(format)

	return t
}

func NewTexture3DFrom(texture Texture3D) *Texture3D {
	t := &Texture3D{}

	t.textureType = gl.TEXTURE_3D

	t.SetName("Texture3D")
	GetInstance().MustAssign(t)

	t.size = texture.Size()
	t.uploadFunc = t.Upload

	t.internalFormat = texture.GLInternalFormat()
	t.glFormat = texture.GLFormat()
	t.storageFormat = texture.GLStorageFormat()

	return t
}

func (t *Texture3D) Upload() {
	t.Bind()

	gl.TexImage3D(t.textureType, 0, t.internalFormat, t.size.X(), t.size.Y(), t.layers, 0, t.glFormat, t.storageFormat, nil)
}

type TextureColor struct {
	BaseTexture

	color Color
}

func NewTextureColor(color Color) *TextureColor {
	t := &TextureColor{}

	t.textureType = gl.TEXTURE_3D

	t.SetName("TextureColor")
	GetInstance().MustAssign(t)

	t.uploadFunc = t.Upload

	return t
}

func (t *TextureColor) Upload() {
	t.Bind()
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA32F, t.size.X(), t.size.Y(), 0, gl.RGBA, gl.FLOAT, gl.Ptr(t.color))
}

func (t *TextureColor) Color() Color {
	return t.color
}

func (t *TextureColor) SetColor(color Color) {
	t.color = color
	t.uploadFunc()
}

type TextureCubemap struct {
	BaseTexture

	data    [6][]uint8
	hdrData [6][]float32
}

func NewTextureCubemap(size math.IVec2, format TextureFormat) *TextureCubemap {
	t := &TextureCubemap{}

	t.data = [6][]uint8{}
	t.hdrData = [6][]float32{}

	t.textureType = gl.TEXTURE_CUBE_MAP

	t.SetName("TextureCubemap")
	GetInstance().MustAssign(t)

	t.size = size
	t.uploadFunc = t.Upload

	t.internalFormat = TextureFormatToInternal(format)
	t.glFormat = TextureFormatToFormat(format)
	t.storageFormat = TextureFormatToStorage(format)

	return t
}

func (t *TextureCubemap) SetData(data []byte, offset int) {
	if offset > 5 {
		return
	}

	t.data[offset] = data
}

func (t *TextureCubemap) SetHDRData(data []float32, offset int) {
	if offset > 5 {
		return
	}

	t.hdrData[offset] = data
}

func (t *TextureCubemap) Upload() {
	t.Bind()

	if len(t.hdrData[0]) > 0 {
		for i := range t.hdrData {
			gl.TexImage2D(
				gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
				0,
				t.internalFormat,
				t.size.X(),
				t.size.Y(),
				0,
				t.glFormat,
				t.storageFormat,
				gl.Ptr(t.hdrData[i]),
			)
		}
	} else if len(t.data[0]) > 0 {
		for i := range t.data {
			gl.TexImage2D(
				gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
				0,
				t.internalFormat,
				t.size.X(),
				t.size.Y(),
				0,
				t.glFormat,
				t.storageFormat,
				gl.Ptr(t.data[i]),
			)
		}
	} else {
		for i := uint32(0); i < 6; i++ {
			gl.TexImage2D(
				gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
				0,
				t.internalFormat,
				t.size.X(),
				t.size.Y(),
				0,
				t.glFormat,
				t.storageFormat,
				nil,
			)
		}
	}
}
