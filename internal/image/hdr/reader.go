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

// Package hdr implements an image.Image-compliant reader for the Radiance HDR
// image format.
package hdr

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"math"
)

const (
	radianceHeader = "#?RADIANCE\n"
)

type decoder struct {
	r             io.Reader
	img           image.Image
	width, height int
	depth         int
	headerSize    int
	tmp           [3 * 256]byte
}

// FormatError reports that the input is not a valid HDR image.
type FormatError string

func (e FormatError) Error() string {
	return "hdr: invalid format: " + string(e)
}

// UnsupportedError reports that the input uses a valid but unimplemented HDR feature.
type UnsupportedError string

func (e UnsupportedError) Error() string {
	return "hdr: unsupported feature: " + string(e)
}

func init() {
	image.RegisterFormat("hdr", radianceHeader, Decode, DecodeConfig)
}

func (d *decoder) parseHeader(b *bufio.Reader) error {
	var x, y string
	var w, h int

	for {
		line, err := b.ReadString('\n')

		if err != nil {
			return err
		}
		if len(line) == 1 {
			break
		}
	}

	line, err := b.ReadString('\n')

	_, err = fmt.Sscanf(line, "%s %d %s %d", &y, &h, &x, &w)
	if err != nil {
		return err
	}

	// TODO: Handle different X, Y scenarios

	d.width = w
	d.height = h

	return nil
}

func (d *decoder) parseData(b *bufio.Reader) error {
	line := make([]byte, d.width*4)

	for y := 0; y < d.height; y++ {
		if err := readLine(b, line); err != nil {
			return err
		}

		for x := 0; x < d.width; x++ {
			i := x * 4

			r := ldexp(line[i+3], line[i+0])
			g := ldexp(line[i+3], line[i+1])
			b := ldexp(line[i+3], line[i+2])

			d.img.(*RGB96).Set(x, y, RGB96Color{
				R: r,
				G: g,
				B: b,
			})
		}
	}

	return nil
}

func readLine(r *bufio.Reader, line []byte) error {
	var code byte
	var value byte

	lineLength := len(line) / 4

	lineHeader, err := r.Peek(4)
	if err != nil {
		return err
	}

	if lineHeader[0] != 2 || lineHeader[1] != 2 || (lineHeader[2]&128) != 0 {
		return readUncompressedData(r, line)
	}

	hlen := (int(lineHeader[2]) << 8) | int(lineHeader[3])
	if hlen != lineLength {
		return FormatError(fmt.Sprintf("scanline length mismatch. have: %d want: %d", hlen, lineLength))
	}

	if _, err := r.Read(lineHeader); err != nil {
		return err
	}

	for i := 0; i < 4; i++ {
		j := 0

		for j < lineLength {
			if code, err = r.ReadByte(); err != nil {
				return err
			}

			if code > 128 {
				code &= 127
				if value, err = r.ReadByte(); err != nil {
					return err
				}

				for k := 0; k < int(code); k++ {
					line[j*4+i] = value
					j++
				}
			} else {
				for k := 0; k < int(code); k++ {
					if value, err = r.ReadByte(); err != nil {
						return err
					}

					line[j*4+i] = value
					j++
				}
			}
		}
	}

	return nil
}

func readUncompressedData(r *bufio.Reader, data []byte) error {
	length := len(data) / 4

	s := make([]byte, 4)

	rshift := uint(0)
	l := 0

	for l < length {
		if _, err := r.Read(s); err != nil {
			return err
		}

		if s[0] == 1 && s[1] == 1 && s[2] == 1 {
			// Encoded
			count := int(s[3]) << rshift
			for i := 0; i < count; i++ {
				data[(l+i)*4+0] = s[0]
				data[(l+i)*4+1] = s[1]
				data[(l+i)*4+2] = s[2]
				data[(l+i)*4+3] = s[3]
			}

			l += count
			rshift += 8
		} else {
			data[l*4+0] = s[0]
			data[l*4+1] = s[1]
			data[l*4+2] = s[2]
			data[l*4+3] = s[3]
			l++
			rshift = 0
		}
	}

	return nil
}

func (d *decoder) Read(p []byte) (int, error) {
	return 0, nil
}

func (d *decoder) checkHeader() error {
	_, err := io.ReadFull(d.r, d.tmp[:len(radianceHeader)])
	if err != nil {
		return err
	}

	if string(d.tmp[:len(radianceHeader)]) != radianceHeader {
		return FormatError("not an HDR file")
	}

	return nil
}

// Decode reads an HDR image from r and returns it as an image.Image.
// The type of Image returned depends on the HDR contents.
func Decode(r io.Reader) (image.Image, error) {
	d := &decoder{
		r: r,
	}

	if err := d.checkHeader(); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}

	b := bufio.NewReader(d.r)

	if err := d.parseHeader(b); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}

	d.img = NewRGB96(
		image.Rectangle{
			Min: image.Point{},
			Max: image.Point{X: d.width, Y: d.height},
		})

	if err := d.parseData(b); err != nil {
		return nil, err
	}

	return d.img, nil
}

// DecodeConfig returns the color model and dimensions of a PNG image without
// decoding the entire image.
func DecodeConfig(r io.Reader) (image.Config, error) {
	d := &decoder{
		r: r,
	}

	if err := d.checkHeader(); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return image.Config{}, err
	}

	b := bufio.NewReader(d.r)

	if err := d.parseHeader(b); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return image.Config{}, err
	}

	return image.Config{
		ColorModel: RGB96Model,
		Width:      d.width,
		Height:     d.height,
	}, nil
}

func ldexp(exp, val uint8) float32 {
	f := float32(math.Ldexp(1.0, int(exp)-int(128+8)))

	return f * float32(val)
}
