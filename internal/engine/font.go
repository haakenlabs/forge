/*
Copyright (c) 2018 HaakenLabs

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
	"image"
	"image/draw"
	"math"
	"sort"
	"unicode"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/golang/freetype/truetype"
	"github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"

	fmath "github.com/haakenlabs/forge/internal/math"
)

var ASCII []rune

func init() {
	ASCII = make([]rune, unicode.MaxASCII-32)
	for i := range ASCII {
		ASCII[i] = rune(32 + i)
	}
}

type Glyph struct {
	Dot     mgl64.Vec2
	Frame   Rect64
	Advance float64
}

type Atlas struct {
	face       font.Face
	texture    *TextureFont
	mapping    map[rune]Glyph
	ascent     float64
	descent    float64
	lineHeight float64
}

type Font struct {
	BaseObject

	ttf     *truetype.Font
	atlases map[float64]*Atlas
	runes   []rune
}

type Rect64 struct {
	Min, Max mgl64.Vec2
}

func R64(minX, minY, maxX, maxY float64) Rect64 {
	return Rect64{
		Min: mgl64.Vec2{minX, minY},
		Max: mgl64.Vec2{maxX, maxY},
	}
}

func (r Rect64) Norm() Rect64 {
	return Rect64{
		Min: mgl64.Vec2{
			math.Min(r.Min.X(), r.Max.X()),
			math.Min(r.Min.Y(), r.Max.Y()),
		},
		Max: mgl64.Vec2{
			math.Max(r.Min.X(), r.Max.X()),
			math.Max(r.Min.Y(), r.Max.Y()),
		},
	}
}

func (r Rect64) Moved(delta mgl64.Vec2) Rect64 {
	return Rect64{
		Min: r.Min.Add(delta),
		Max: r.Max.Add(delta),
	}
}

func (r Rect64) Union(s Rect64) Rect64 {
	return R64(
		math.Min(r.Min.X(), s.Min.X()),
		math.Min(r.Min.Y(), s.Min.Y()),
		math.Max(r.Max.X(), s.Max.X()),
		math.Max(r.Max.Y(), s.Max.Y()),
	)
}

func (r Rect64) W() float64 {
	return r.Max.X() - r.Min.X()
}

func (r Rect64) H() float64 {
	return r.Max.Y() - r.Min.Y()
}

type fixedGlyph struct {
	dot     fixed.Point26_6
	frame   fixed.Rectangle26_6
	advance fixed.Int26_6
}

func NewFont(ttf *truetype.Font, runeSets ...[]rune) *Font {
	f := &Font{
		ttf:     ttf,
		atlases: make(map[float64]*Atlas),
	}

	seen := make(map[rune]struct{})
	runes := []rune{unicode.ReplacementChar}

	for _, set := range runeSets {
		for _, r := range set {
			if _, ok := seen[r]; !ok {
				runes = append(runes, r)
				seen[r] = struct{}{}
			}
		}
	}

	f.runes = runes
	f.SetName("Font")
	GetInstance().MustAssign(f)

	return f
}

func (f *Font) Atlas(size float64) *Atlas {
	if atlas, ok := f.atlases[size]; ok {
		return atlas
	}

	return f.generateAtlas(size)
}

func (f *Font) generateAtlas(size float64) *Atlas {
	if size <= 0 {
		logrus.Errorf("Invalid font size: %f", size)
		return nil
	}

	face := truetype.NewFace(f.ttf, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})

	fixedMapping, fixedBounds := makeSquareMapping(face, f.runes, fixed.I(2))

	atlasImg := image.NewRGBA(image.Rect(
		fixedBounds.Min.X.Floor(),
		fixedBounds.Min.Y.Floor(),
		fixedBounds.Max.X.Ceil(),
		fixedBounds.Max.Y.Ceil(),
	))

	for r, fg := range fixedMapping {
		dr, mask, maskp, _, _ := face.Glyph(fg.dot, r)
		draw.Draw(atlasImg, dr, mask, maskp, draw.Src)
	}

	bounds := R64(
		i2f(fixedBounds.Min.X),
		i2f(fixedBounds.Min.Y),
		i2f(fixedBounds.Max.X),
		i2f(fixedBounds.Max.Y),
	)

	mapping := make(map[rune]Glyph)
	for r, fg := range fixedMapping {

		mapping[r] = Glyph{
			Dot: mgl64.Vec2{
				i2f(fg.dot.X),
				i2f(fg.dot.Y),
			},
			Frame: R64(
				i2f(fg.frame.Min.X),
				i2f(fg.frame.Min.Y)+(i2f(fg.frame.Max.Y-fg.frame.Min.Y))-bounds.Min.Y(),
				i2f(fg.frame.Max.X),
				i2f(fg.frame.Min.Y)-bounds.Min.Y(),
			).Norm(),
			Advance: i2f(fg.advance),
		}
	}

	atlas := &Atlas{
		face:       face,
		mapping:    mapping,
		ascent:     i2f(face.Metrics().Ascent),
		descent:    i2f(face.Metrics().Descent),
		lineHeight: i2f(face.Metrics().Height),
	}

	atlas.texture = NewTextureFont(fmath.IVec2{
		int32(atlasImg.Bounds().Dx()),
		int32(atlasImg.Bounds().Dy()),
	})
	atlas.texture.SetData(atlasImg.Pix)
	atlas.texture.Alloc()

	f.atlases[size] = atlas

	return atlas
}

func (f *Font) HasSize(size float64) bool {
	_, ok := f.atlases[size]

	return ok
}

func (f *Font) DrawText(text string, size float64) ([]Vertex, mgl32.Vec2) {
	var atlas *Atlas
	var dot mgl64.Vec2
	var prev rune
	var vi int
	var boundings Rect64

	if text == "" {
		return nil, mgl32.Vec2{}
	}

	if !f.HasSize(size) {
		atlas = f.generateAtlas(size)
	} else {
		atlas = f.atlases[size]
	}

	if atlas == nil {
		return nil, mgl32.Vec2{}
	}

	verts := make([]Vertex, 6*len(text))
	tw := float32(atlas.Texture().Width())
	th := float32(atlas.Texture().Height())

	for _, r := range text {
		var rect, frame, bounds Rect64
		rect, frame, bounds, dot = atlas.DrawRune(prev, r, dot)

		prev = r

		ul := Vertex{
			V: mgl32.Vec3{float32(rect.Min.X()), float32(rect.Min.Y()), 0},
			U: mgl32.Vec2{float32(frame.Min.X()) / tw, float32(frame.Min.Y()) / th},
		}
		ur := Vertex{
			V: mgl32.Vec3{float32(rect.Max.X()), float32(rect.Min.Y()), 0},
			U: mgl32.Vec2{float32(frame.Max.X()) / tw, float32(frame.Min.Y()) / th},
		}
		lr := Vertex{
			V: mgl32.Vec3{float32(rect.Max.X()), float32(rect.Max.Y()), 0},
			U: mgl32.Vec2{float32(frame.Max.X()) / tw, float32(frame.Max.Y()) / th},
		}
		ll := Vertex{
			V: mgl32.Vec3{float32(rect.Min.X()), float32(rect.Max.Y()), 0},
			U: mgl32.Vec2{float32(frame.Min.X()) / tw, float32(frame.Max.Y()) / th},
		}

		verts[vi] = ul
		verts[vi+1] = lr
		verts[vi+2] = ur
		verts[vi+3] = ul
		verts[vi+4] = ll
		verts[vi+5] = lr

		if boundings.W()*boundings.H() == 0 {
			boundings = bounds
		} else {
			boundings = boundings.Union(bounds)
		}

		vi += 6
	}

	return verts, mgl32.Vec2{}
}

func (a *Atlas) Texture() *TextureFont {
	return a.texture
}

// Contains reports whether r in contained within the Atlas.
func (a *Atlas) Contains(r rune) bool {
	_, ok := a.mapping[r]
	return ok
}

// Glyph returns the description of r within the Atlas.
func (a *Atlas) Glyph(r rune) (Glyph, bool) {
	g, ok := a.mapping[r]

	return g, ok
}

// Kern returns the kerning distance between runes r0 and r1. Positive distance means that the
// glyphs should be further apart.
func (a *Atlas) Kern(r0, r1 rune) float64 {
	return i2f(a.face.Kern(r0, r1))
}

// Ascent returns the distance from the top of the line to the baseline.
func (a *Atlas) Ascent() float64 {
	return a.ascent
}

// Descent returns the distance from the baseline to the bottom of the line.
func (a *Atlas) Descent() float64 {
	return a.descent
}

// LineHeight returns the recommended vertical distance between two lines of text.
func (a *Atlas) LineHeight() float64 {
	return a.lineHeight
}

func (a *Atlas) DrawRune(prev, r rune, dot mgl64.Vec2) (rect, frame, bounds Rect64, newDot mgl64.Vec2) {
	if !a.Contains(r) {
		logrus.Errorf("unknown rune: %s", string(r))
		r = unicode.ReplacementChar
	}
	if !a.Contains(unicode.ReplacementChar) {
		logrus.Errorf("unknown replacement rune: %s", string(unicode.ReplacementChar))
		return Rect64{}, Rect64{}, Rect64{}, dot
	}
	if !a.Contains(prev) {
		prev = unicode.ReplacementChar
	}

	if prev >= 0 {
		dot[0] += a.Kern(prev, r)
	}

	glyph, _ := a.Glyph(r)

	rect = glyph.Frame.Moved(dot.Sub(glyph.Dot))
	bounds = rect

	if bounds.W()*bounds.H() != 0 {
		bounds = R64(
			bounds.Min.X(),
			dot.Y()-a.Descent(),
			bounds.Max.X(),
			dot.Y()+a.Ascent(),
		)
	}

	dot[0] += glyph.Advance

	return rect, glyph.Frame, bounds, dot
}

func makeSquareMapping(face font.Face, runes []rune, padding fixed.Int26_6) (map[rune]fixedGlyph, fixed.Rectangle26_6) {
	width := sort.Search(int(fixed.I(1024*1024)), func(i int) bool {
		width := fixed.Int26_6(i)
		_, bounds := makeMapping(face, runes, padding, width)
		return bounds.Max.X-bounds.Min.X >= bounds.Max.Y-bounds.Min.Y
	})
	return makeMapping(face, runes, padding, fixed.Int26_6(width))
}

func makeMapping(face font.Face, runes []rune, padding, width fixed.Int26_6) (map[rune]fixedGlyph, fixed.Rectangle26_6) {
	mapping := make(map[rune]fixedGlyph)
	bounds := fixed.Rectangle26_6{}

	dot := fixed.P(0, 0)

	for _, r := range runes {
		b, advance, ok := face.GlyphBounds(r)
		if !ok {
			logrus.Error("Missing rune: %v", r)
			continue
		}

		// this is important for drawing, artifacts arise otherwise
		frame := fixed.Rectangle26_6{
			Min: fixed.P(b.Min.X.Floor(), b.Min.Y.Floor()),
			Max: fixed.P(b.Max.X.Ceil(), b.Max.Y.Ceil()),
		}

		dot.X -= frame.Min.X
		frame = frame.Add(dot)

		mapping[r] = fixedGlyph{
			dot:     dot,
			frame:   frame,
			advance: advance,
		}
		bounds = bounds.Union(frame)

		dot.X = frame.Max.X

		// padding + align to integer
		dot.X += padding
		dot.X = fixed.I(dot.X.Ceil())

		// width exceeded, new row
		if frame.Max.X >= width {
			dot.X = 0
			dot.Y += face.Metrics().Ascent + face.Metrics().Descent

			// padding + align to integer
			dot.Y += padding
			dot.Y = fixed.I(dot.Y.Ceil())
		}
	}

	return mapping, bounds
}

func i2f(i fixed.Int26_6) float64 {
	return float64(i) / (1 << 6)
}
