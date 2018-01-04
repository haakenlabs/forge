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
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/sirupsen/logrus"
)

const (
	fixedTime    = float64(0.05)
	maxFrameSkip = 5

	builtinAssets = "<builtin>:builtin.json"
)

var (
	app     *App
	appOnce sync.Once
)

// AppConfig provides options for configuring an App.
type AppConfig struct {
	Name string
}

// App is the backbone of any Apex application.
type App struct {
	scenes           map[string]*Scene
	systems          []System
	activeScenes     []string
	preStartFunc     func() error
	postStartFunc    func() error
	preTeardownFunc  func()
	postTeardownFunc func()
	name             string
	running          bool
}

// NewApp creates a new App using the provided AppConfig for customization.
func NewApp(cfg *AppConfig) *App {
	a := &App{
		name: cfg.Name,
	}

	return a
}

// CurrentApp returns the current app.
func CurrentApp() *App {
	return app
}

// SetApp sets the App, but only once.
func setApp(a *App) {
	appOnce.Do(func() {
		app = a
	})
}

// Setup sets up the App.
func (a *App) Setup() error {
	setApp(a)

	LoadGlobalConfig()

	a.scenes = make(map[string]*Scene)
	a.systems = []System{}
	a.activeScenes = []string{}

	// Register required systems.
	a.RegisterSystem(NewWindow())
	a.RegisterSystem(NewInstance())
	a.RegisterSystem(NewAsset())
	a.RegisterSystem(NewTime())

	// Register asset handlers.
	asset := GetAsset()
	asset.RegisterHandler(NewImageHandler())
	asset.RegisterHandler(NewMeshHandler())
	asset.RegisterHandler(NewShaderHandler())
	asset.RegisterHandler(NewSkyboxHandler())

	if a.preStartFunc != nil {
		if err := a.preStartFunc(); err != nil {
			return err
		}
	}

	for i := range a.systems {
		logrus.Debug("Setting up system: ", a.systems[i].Name())

		if err := a.systems[i].Setup(); err != nil {
			return err
		}
	}

	// Load base assets.
	if err := asset.LoadManifest(builtinAssets); err != nil {
		return err
	}

	if a.postStartFunc != nil {
		if err := a.postStartFunc(); err != nil {
			return err
		}
	}

	return nil
}

// Teardown tears down the app.
func (a *App) Teardown() {
	if a.preTeardownFunc != nil {
		a.preTeardownFunc()
	}

	for i := len(a.systems) - 1; i >= 0; i-- {
		logrus.Debug("Tearing down system: ", a.systems[i].Name())

		a.systems[i].Teardown()
	}

	if a.postTeardownFunc != nil {
		a.postTeardownFunc()
	}
}

// Run starts the main loop of the app.
func (a *App) Run() error {
	a.running = true

	a.setupSignalHandler()

	frame := 0
	loops := 0

	time := a.MustSystem(SysNameTime).(*Time)
	window := a.MustSystem(SysNameWindow).(*Window)

	for a.running {
		a.running = !window.ShouldClose()

		time.FrameStart()

		frame++

		if window.KeyDown(glfw.KeyF9) {
			a.debugInfo()
		}

		a.onUpdate()

		loops = 0
		for time.LogicUpdate() && loops < maxFrameSkip {
			time.LogicTick()
			a.onFixedUpdate()
			loops++
		}

		window.ClearBuffers()
		a.onDisplay()
		window.SwapBuffers()

		window.HandleEvents()
		time.FrameEnd()
	}

	return nil
}

// Quit instructs the App to shutdown by setting the running variable to false.
func (a *App) Quit() {
	a.running = false
}

func (a *App) Name() string {
	return a.name
}

// RegisterSystem registers a system with the App. A system can only be added
// once, it is an error to add a system more than once. Systems are initialized
// in the order they are added and torn down in the reverse order.
func (a *App) RegisterSystem(s System) {
	// Check for existing system.
	if a.SystemRegistered(s.Name()) {
		panic(ErrSystemExists(s.Name()))
	}

	// Add system to the systems slice.
	a.systems = append(a.systems, s)

	logrus.Debug("Registered system: ", s.Name())
}

// SystemRegistered returns true if the system with the given name is registered
// with this App.
func (a *App) SystemRegistered(name string) bool {
	for i := range a.systems {
		if a.systems[i].Name() == name {
			return true
		}
	}

	return false
}

// System returns a system by the given name.
func (a *App) System(name string) (System, error) {
	for i := range a.systems {
		if a.systems[i].Name() == name {
			return a.systems[i], nil
		}
	}

	return nil, ErrSystemNotFound(name)
}

// MustSystem is like System, but panics if the system cannot be found.
func (a *App) MustSystem(name string) System {
	s, err := a.System(name)
	if err != nil {
		panic(err)
	}

	return s
}

// SetPreSetup sets a callback which will be invoked prior to the setup function call.
func (a *App) SetPreSetup(fn func() error) {
	a.preStartFunc = fn
}

// SetPostSetup sets a callback which will be invoked after to the setup function call.
func (a *App) SetPostSetup(fn func() error) {
	a.postStartFunc = fn
}

// SetPreTeardown sets a callback which will be invoked prior to the teardown function call.
func (a *App) SetPreTeardown(fn func()) {
	a.preTeardownFunc = fn
}

// SetPostTeardown sets a callback which will be invoked after to the teardown function call.
func (a *App) SetPostTeardown(fn func()) {
	a.postTeardownFunc = fn
}

func (a *App) RegisterScene(scene *Scene) error {
	if a.SceneRegistered(scene.Name()) {
		return fmt.Errorf("register scene: '%s' already registered", scene.Name())
	}

	a.scenes[scene.Name()] = scene

	logrus.Debug("Registered new scene: ", scene.Name())

	return nil
}

func (a *App) LoadScene(name string) error {
	if !a.SceneRegistered(name) {
		return fmt.Errorf("load scene: '%s' not registered", name)
	}

	if !a.scenes[name].Loaded() {
		return a.scenes[name].Load()
	}

	return nil
}

func (a *App) PurgePushScene(name string) error {
	if !a.SceneRegistered(name) {
		return fmt.Errorf("purge push scene: '%s' not registered", name)
	}

	a.activeScenes = a.activeScenes[:0]
	a.PushScene(name)

	return nil
}

func (a *App) ReplaceScene(name string) error {
	if !a.SceneRegistered(name) {
		return fmt.Errorf("replace error: '%s' not registered", name)
	}

	a.PopScene()
	a.PushScene(name)

	return nil
}

func (a *App) PushScene(name string) error {
	if !a.SceneRegistered(name) {
		return fmt.Errorf("push: '%s' not registered", name)
	}

	if err := a.LoadScene(name); err != nil {
		return err
	}

	a.activeScenes = append(a.activeScenes, name)
	a.scenes[name].OnActivate()

	return nil
}

func (a *App) PopScene() string {
	if len(a.activeScenes) != 0 {
		last := a.activeScenes[len(a.activeScenes)-1]

		a.activeScenes = a.activeScenes[:len(a.activeScenes)-1]
		a.scenes[last].OnDeactivate()

		return last
	}

	return ""
}

func (a *App) RemoveAllScenes() {
	for key := range a.scenes {
		delete(a.scenes, key)
	}
	a.activeScenes = a.activeScenes[:0]
}

func (a *App) UnregisterScene(name string) error {
	if !a.SceneRegistered(name) {
		return fmt.Errorf("unregister scene: '%s' not registered", name)
	}

	delete(a.scenes, name)

	return nil
}

func (a *App) SceneRegistered(name string) bool {
	_, ok := a.scenes[name]

	return ok
}

func (a *App) ActiveScene() *Scene {
	if a.ActiveSceneCount() != 0 {
		name := a.activeScenes[len(a.activeScenes)-1]
		return a.scenes[name]
	}

	return nil
}

func (a *App) ActiveSceneName() string {
	if sc := a.ActiveScene(); sc != nil {
		return sc.Name()
	}

	return ""
}

func (a *App) SceneCount() int {
	return len(a.scenes)
}

func (a *App) ActiveSceneCount() int {
	return len(a.activeScenes)
}

func (a *App) setupSignalHandler() {
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt, syscall.SIGTERM)
	go handleSignal(s, a)
}

func handleSignal(s chan os.Signal, a *App) {
	<-s
	a.Quit()
}

func (a *App) onDisplay() {
	if s := a.ActiveScene(); s != nil {
		sg := s.Graph()
		if sg.Dirty() {
			sg.Update()
		}

		cameras := s.cameras
		for i := range cameras {
			cameras[i].Render()
		}

		sg.SendMessage(MessageGUIDisplay)
	}
}

func (a *App) onUpdate() {
	if s := a.ActiveScene(); s != nil {
		sg := s.Graph()
		if sg.Dirty() {
			sg.Update()
		}

		sg.SendMessage(MessageUpdate)
		sg.SendMessage(MessageLateUpdate)
	}
}

func (a *App) onFixedUpdate() {
	if s := a.ActiveScene(); s != nil {
		sg := s.Graph()
		if sg.Dirty() {
			sg.Update()
		}

		sg.SendMessage(MessageFixedUpdate)
	}
}

func (a *App) debugInfo() {
	fmt.Printf("Application debug info -------------\n")
	fmt.Printf("  App name: %s\n", a.Name())
	fmt.Printf("  Active scene: %s\n", a.ActiveSceneName())
	fmt.Printf("  Registered scenes: %d\n", a.SceneCount())
}
