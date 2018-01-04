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
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"git.dbnservers.net/haakenlabs/forge/internal/math"
)

var _ System = &Window{}

const SysNameWindow = "window"

const (
	DisplayModeWindow DisplayMode = iota
	DisplayModeWindowedFullscreen
	DisplayModeFullscreen
)

const (
	MouseAbsolute MouseMode = iota
	MouseRelative
)

type DisplayMode int
type MouseMode int

type EventKey struct {
	key      glfw.Key
	scancode int
	action   glfw.Action
	mods     glfw.ModifierKey
}

type EventMouseButton struct {
	button glfw.MouseButton
	action glfw.Action
	mod    glfw.ModifierKey
}

type EventJoy struct {
	joystick int
	event    int
}

type DisplayProperties struct {
	Resolution math.IVec2
	Mode       DisplayMode
	Vsync      bool
}

// Window implements a GLFW-based window system.
type Window struct {
	window            *glfw.Window
	ortho             mgl32.Mat4
	resolution        math.IVec2
	mousePos          math.DVec2
	scrollAxis        math.DVec2
	cursorPosition    mgl32.Vec2
	mouseMode         MouseMode
	displayMode       DisplayMode
	mouseButtonEvents []EventMouseButton
	keyEvents         []EventKey
	joystickEvents    []EventJoy
	aspectRatio       float32
	title             string
	vsync             bool
	focus             bool
	cursorEnter       bool
	cursorMoved       bool
	scrollMoved       bool
	windowResized     bool
	shouldClose       bool
	hasEvents         bool
}

func (w *Window) Setup() (err error) {
	var monitor *glfw.Monitor

	if err := glfw.Init(); err != nil {
		return err
	}

	logrus.Debug("[GLFW] Library initialized")

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	w.displayMode = DisplayMode(viper.GetInt("graphics.mode"))
	w.resolution = math.ToIVec2(viper.Get("graphics.resolution"))
	w.vsync = viper.GetBool("graphics.vsync")

	resX := int(w.resolution.X())
	resY := int(w.resolution.Y())

	switch w.displayMode {
	case DisplayModeWindowedFullscreen:
		monitor = glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()

		glfw.WindowHint(glfw.RedBits, mode.RedBits)
		glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
		glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
		glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)

		resX = mode.Width
		resY = mode.Height
	case DisplayModeFullscreen:
		monitor = glfw.GetPrimaryMonitor()
		vidmode := GetRecommendedVideoMode(monitor)

		glfw.WindowHint(glfw.RedBits, vidmode.RedBits)
		glfw.WindowHint(glfw.GreenBits, vidmode.GreenBits)
		glfw.WindowHint(glfw.BlueBits, vidmode.BlueBits)
		glfw.WindowHint(glfw.RefreshRate, vidmode.RefreshRate)

		resX = vidmode.Width
		resY = vidmode.Height
	}

	w.resolution = math.IVec2{int32(resX), int32(resY)}

	if w.window, err = glfw.CreateWindow(resX, resY, CurrentApp().name, monitor, nil); err != nil {
		return err
	}

	if err := w.setupGL(); err != nil {
		return err
	}

	// Register input callbacks
	w.window.SetCharCallback(w.onChar)
	w.window.SetCursorEnterCallback(w.onCursorEnter)
	w.window.SetCursorPosCallback(w.onCursorMove)
	w.window.SetDropCallback(w.onDrop)
	w.window.SetKeyCallback(w.onKey)
	w.window.SetMouseButtonCallback(w.onMouseButton)
	w.window.SetScrollCallback(w.onScroll)
	w.window.SetCloseCallback(w.onClose)
	w.window.SetSizeCallback(w.onWindowResize)
	glfw.SetJoystickCallback(w.onJoystick)

	logrus.Debug("[GLFW] Ready")

	return nil
}

func (w *Window) setupGL() error {
	w.window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return err
	}

	logrus.Debug("[OpenGL] Version: ", gl.GoStr(gl.GetString(gl.VERSION)))

	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_CUBE_MAP_SEAMLESS)
	gl.DepthFunc(gl.LEQUAL)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	w.SetSize(w.resolution)

	w.EnableVsync(w.vsync)

	logrus.Debug("[OpenGL] Ready")

	return nil
}

// Teardown tears down the System.
func (w *Window) Teardown() {
	glfw.Terminate()
}

// Name returns the name of the System.
func (w *Window) Name() string {
	return SysNameWindow
}

func (w *Window) EnableVsync(enable bool) {
	if enable {
		glfw.SwapInterval(1)
	} else {
		glfw.SwapInterval(0)
	}

	w.vsync = enable
}

func (w *Window) Vsync() bool {
	return w.vsync
}

func (w *Window) CenterWindow() {
	monitor := w.window.GetMonitor()
	if monitor == nil {
		monitor = glfw.GetPrimaryMonitor()
	}

	vmode := monitor.GetVideoMode()

	posX, posY := (vmode.Width/2)-(int(w.resolution.X())/2), (vmode.Height/2)-(int(w.resolution.Y())/2)

	w.window.SetPos(posX, posY)
}

func (w *Window) SetSize(size math.IVec2) {
	w.resolution = size
	w.aspectRatio = getRatio(w.resolution)
	gl.Viewport(0, 0, int32(size.X()), int32(size.Y()))
	w.ortho = mgl32.Ortho2D(0, float32(w.resolution.X()), float32(w.resolution.Y()), 0)
}

// AspectRatio : Get the aspect ratio of the WindowSystem.
func (w *Window) AspectRatio() float32 {
	return w.aspectRatio
}

// Resolution : Get the resolution of the WindowSystem.
func (w *Window) Resolution() math.IVec2 {
	return w.resolution
}

func (w *Window) ClearBuffers() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

// SwapBuffers : Swap front and rear rendering buffers.
func (w *Window) SwapBuffers() {
	w.window.SwapBuffers()
}

func (w *Window) GLFWWindow() *glfw.Window {
	return w.window
}

func (w *Window) OrthoMatrix() mgl32.Mat4 {
	return w.ortho
}

func (w *Window) SetDisplayMode(mode DisplayMode) {
	var monitor *glfw.Monitor
	var refresh int

	posX, posY := w.window.GetPos()
	resX := int(w.resolution.X())
	resY := int(w.resolution.Y())

	switch mode {
	case DisplayModeWindowedFullscreen:
		monitor = glfw.GetPrimaryMonitor()
		mode := monitor.GetVideoMode()

		glfw.WindowHint(glfw.RedBits, mode.RedBits)
		glfw.WindowHint(glfw.GreenBits, mode.GreenBits)
		glfw.WindowHint(glfw.BlueBits, mode.BlueBits)
		glfw.WindowHint(glfw.RefreshRate, mode.RefreshRate)

		resX = mode.Width
		resY = mode.Height
		refresh = mode.RefreshRate
	case DisplayModeFullscreen:
		vidmode := GetRecommendedVideoMode(glfw.GetPrimaryMonitor())

		glfw.WindowHint(glfw.RedBits, vidmode.RedBits)
		glfw.WindowHint(glfw.GreenBits, vidmode.GreenBits)
		glfw.WindowHint(glfw.BlueBits, vidmode.BlueBits)
		glfw.WindowHint(glfw.RefreshRate, vidmode.RefreshRate)

		resX = vidmode.Width
		resY = vidmode.Height
		refresh = vidmode.RefreshRate
	}

	w.window.SetMonitor(monitor, posX, posY, resX, resY, refresh)
}

func (w *Window) GetVideoModes() {
	monitors := glfw.GetMonitors()
	modes := []*glfw.VidMode{}

	for i := range monitors {
		modes = append(modes, monitors[i].GetVideoModes()...)
	}

	for i := range modes {
		fmt.Printf("%dx%dx%d (%d Hz)\n", modes[i].Width, modes[i].Height, modes[i].RedBits+modes[i].GreenBits+modes[i].BlueBits, modes[i].RefreshRate)
	}
}

func (w *Window) ShouldClose() bool {
	return w.shouldClose
}

func (w *Window) KeyDown(key glfw.Key) bool {
	for idx := range w.keyEvents {
		if w.keyEvents[idx].key == key {
			if w.keyEvents[idx].action == glfw.Press || w.keyEvents[idx].action == glfw.Repeat {
				return true
			}
		}
	}

	return false
}

func (w *Window) KeyUp(key glfw.Key) bool {
	for idx := range w.keyEvents {
		if w.keyEvents[idx].key == key {
			if w.keyEvents[idx].action == glfw.Release {
				return true
			}
		}
	}

	return false
}

func (w *Window) KeyPressed() bool {
	return len(w.keyEvents) != 0
}

func (w *Window) MouseDown(button glfw.MouseButton) bool {
	for idx := range w.mouseButtonEvents {
		if w.mouseButtonEvents[idx].button == button {
			if w.mouseButtonEvents[idx].action == glfw.Press {
				return true
			}
		}
	}

	return false
}

func (w *Window) MouseUp(button glfw.MouseButton) bool {
	for idx := range w.mouseButtonEvents {
		if w.mouseButtonEvents[idx].button == button {
			if w.mouseButtonEvents[idx].action == glfw.Release {
				return true
			}
		}
	}

	return false
}

func (w *Window) MouseWheelX() float64 {
	return w.scrollAxis[0]
}

func (w *Window) MouseWheelY() float64 {
	return w.scrollAxis[1]
}

func (w *Window) MouseWheel() bool {
	return w.scrollMoved
}

func (w *Window) MouseMoved() bool {
	return w.cursorMoved
}

func (w *Window) MousePressed() bool {
	return len(w.mouseButtonEvents) != 0
}

func (w *Window) MousePosition() mgl32.Vec2 {
	return w.cursorPosition
}

func (w *Window) WindowResized() bool {
	return w.windowResized
}

func (w *Window) HandleEvents() {
	w.clearEvents()
	glfw.PollEvents()
}

func (w *Window) HasEvents() bool {
	return w.hasEvents
}

func (w *Window) clearEvents() {
	w.hasEvents = false
	w.mouseButtonEvents = w.mouseButtonEvents[:0]
	w.keyEvents = w.keyEvents[:0]
	w.joystickEvents = w.joystickEvents[:0]
	w.cursorMoved = false
	w.scrollMoved = false
	w.windowResized = false
}

func (w *Window) keyEvent(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	w.hasEvents = true
	w.keyEvents = append(w.keyEvents, EventKey{key, scancode, action, mods})
}

func (w *Window) mouseButtonEvent(button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	w.hasEvents = true
	w.mouseButtonEvents = append(w.mouseButtonEvents, EventMouseButton{button, action, mod})
}

func (w *Window) joystickEvent(joy int, event int) {
	w.hasEvents = true
	w.joystickEvents = append(w.joystickEvents, EventJoy{joy, event})
}

func (w *Window) onChar(_ *glfw.Window, char rune) {
	w.hasEvents = true
}

func (w *Window) onCursorEnter(_ *glfw.Window, entered bool) {
	w.hasEvents = true
	w.cursorEnter = entered
}

func (w *Window) onCursorMove(_ *glfw.Window, xPos float64, yPos float64) {
	w.hasEvents = true
	w.cursorPosition[0] = float32(xPos)
	w.cursorPosition[1] = float32(yPos)
	w.cursorMoved = true
}

func (w *Window) onDrop(_ *glfw.Window, names []string) {
	w.hasEvents = true
	fmt.Printf("onDrop: %v\n", names)
}

func (w *Window) onJoystick(joy int, event int) {
	w.hasEvents = true
	w.joystickEvent(joy, event)
}

func (w *Window) onKey(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	w.hasEvents = true
	w.keyEvent(key, scancode, action, mods)
}

func (w *Window) onMouseButton(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, mod glfw.ModifierKey) {
	w.hasEvents = true
	w.mouseButtonEvent(button, action, mod)
}

func (w *Window) onScroll(_ *glfw.Window, xOff float64, yOff float64) {
	w.hasEvents = true
	w.scrollAxis[0] = xOff
	w.scrollAxis[1] = yOff
	w.scrollMoved = true
}

func (w *Window) onClose(_ *glfw.Window) {
	w.hasEvents = true
	w.shouldClose = true
}

func (w *Window) onWindowResize(_ *glfw.Window, width int, height int) {
	if width > 0 && height > 0 {
		w.hasEvents = true
		w.SetSize(math.IVec2{int32(width), int32(height)})
		w.windowResized = true
	}
}

// NewWindow creates a new window system.
func NewWindow() *Window {
	return &Window{
		focus:             true,
		cursorEnter:       true,
		mouseButtonEvents: make([]EventMouseButton, 4),
		keyEvents:         make([]EventKey, 4),
		joystickEvents:    make([]EventJoy, 4),
	}
}

func GetRecommendedVideoMode(monitor *glfw.Monitor) *glfw.VidMode {
	modes := monitor.GetVideoModes()

	return modes[len(modes)-1]
}

func DefaultDisplayProperties() *DisplayProperties {
	return &DisplayProperties{
		Resolution: math.IVec2{1280, 720},
		Mode:       DisplayModeWindow,
		Vsync:      false,
	}
}

// GetWindow gets the window system from the current app.
func GetWindow() *Window {
	return CurrentApp().MustSystem(SysNameWindow).(*Window)
}

func getRatio(value math.IVec2) float32 {
	return float32(value.X()) / float32(value.Y())
}
