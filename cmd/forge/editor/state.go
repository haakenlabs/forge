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

package editor

import (
	"encoding/gob"
	"io"
	"os"
)

// State represents the editor state.
type State struct {
	Name  string
	Model string
	EnvCfg *EnvironmentConfig
	env    *Environment
}

// Setup sets up the editor state.
func (s *State) Setup() error {
	s.env = NewEnvironment(s.EnvCfg)

	// TODO: load model

	return nil
}

// NewState creates a new editor state according to the parameters given.
func NewState(name, model string) (*State, error) {
	s := &State{
		Name: name,
		Model: model,
		EnvCfg: DefaultEnvironmentConfig(),
	}

	if err := s.Setup(); err != nil {
		return nil, err
	}

	return s, nil
}

// LoadState loads a state from the reader.
func LoadState(r io.Reader) (*State, error) {
	var s *State

	dec := gob.NewDecoder(r)
	if err := dec.Decode(s); err != nil {
		return nil, err
	}

	if err := s.Setup(); err != nil {
		return nil, err
	}

	return s, nil
}

// Save state writes a state to the writer.
func SaveState(w io.Writer, s *State) error {
	enc := gob.NewEncoder(w)
	if err := enc.Encode(s); err != nil {
		return err
	}

	return nil
}

// LoadStateFromFile loads a state from file.
func LoadStateFromFile(filename string) (*State, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return LoadState(r)
}

// SaveStateToFile saves a state to file.
func SaveStateToFile(filename string, s *State) error {
	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	return SaveState(w, s)
}
