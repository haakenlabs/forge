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
	"github.com/spf13/viper"

	"git.dbnservers.net/haakenlabs/forge/internal/math"
)

const (
	cfgFilename = "apex.cfg"
	cfgPrefix = "apex"
)

// LoadGlobalConfig sets up viper and reads in the main configuration.
func LoadGlobalConfig() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix(cfgPrefix)
	viper.SetConfigFile(cfgFilename)
	//viper.AddConfigPath(AppDir)
	viper.SetConfigType("json")

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return err
		}
	}

	loadDefaultSettings()

	return nil
}

// loadDefaultSettings sets default settings.
func loadDefaultSettings() {
	// Graphics Options
	viper.SetDefault("graphics.resolution", math.IVec2{1280, 720})
	viper.SetDefault("graphics.windowed", true)
	viper.SetDefault("graphics.mode", 0)
	viper.SetDefault("graphics.vsync", true)
}
