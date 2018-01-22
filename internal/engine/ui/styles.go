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

package ui

import "github.com/haakenlabs/forge/internal/engine"

var Styles = struct {
	BackgroundColor     engine.Color
	AltBackgroundColor  engine.Color
	Alt2BackgroundColor engine.Color
	PrimaryTextColor    engine.Color
	SecondaryTextColor  engine.Color
	TertiaryTextColor   engine.Color
	InverseTextColor    engine.Color
}{
	BackgroundColor:     engine.Color{0.1, 0.1, 0.1, 0.75},
	AltBackgroundColor:  engine.ColorBlue,
	Alt2BackgroundColor: engine.ColorGreen,
	PrimaryTextColor:    engine.ColorWhite,
	SecondaryTextColor:  engine.ColorYellow,
	TertiaryTextColor:   engine.ColorGreen,
	InverseTextColor:    engine.ColorBlue,
}
