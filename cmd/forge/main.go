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

package main

import (
	"flag"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"

	"git.dbnservers.net/haakenlabs/forge/internal/engine"

	"git.dbnservers.net/haakenlabs/forge/cmd/forge/scene"
)

var (
	debug bool
)

func init() {
	// Set logrus formatter.
	logrus.SetFormatter(&prefixed.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.StampMicro,
	})

	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

// parseArgs parses command line arguments.
func parseArgs() {
	flag.BoolVar(&debug, "d", false, "enable debug mode")

	flag.Parse()

	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug logging enabled. PID: ", os.Getpid())
	}
}

func makeApp() *engine.App {
	// Make the app.
	a := engine.NewApp(&engine.AppConfig{
		Name: "Forge",
	})

	// Set the PostSetup func.
	a.SetPostSetup(func() error {
		// Make the scenes.
		scenes := []*engine.Scene{
			scene.NewStartScene(),
			scene.NewMenuScene(),
			scene.NewOptionsScene(),
			scene.NewEditorScene(),
		}

		// Register the scenes.
		for i := range scenes {
			if err := a.RegisterScene(scenes[i]); err != nil {
				return err
			}
		}

		// Push the start stage.
		if err := a.PushScene(scene.NameStart); err != nil {
			return err
		}

		return nil
	})

	return a
}

func main() {
	// Parse cli arguments.
	parseArgs()

	// Make the app.
	app := makeApp()

	// Setup the app.
	if err := app.Setup(); err != nil {
		logrus.Fatal(err)
	}
	defer app.Teardown()

	// Run the app.
	if err := app.Run(); err != nil {
		logrus.Fatal(err)
	}
}
