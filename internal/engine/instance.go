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
	"github.com/sirupsen/logrus"
	"math"
	"sync"
)

var _ System = &Instance{}

const SysNameInstance = "instance"

const (
	ErrMaxIDsExceeded        = Error("exceeded maximum number of instance IDs")
	ErrAssignNilObject       = Error("cannot assign nil object")
	ErrObjectAlreadyAssigned = Error("object has already been assigned")
)

type ErrIDAlreadyAssigned uint32
type ErrIDNotFound uint32

func (e ErrIDAlreadyAssigned) Error() string {
	return fmt.Sprintf("object with ID %08X already assigned", e)
}

func (e ErrIDNotFound) Error() string {
	return fmt.Sprintf("object with ID %08X not found", e)
}

// Instance implements a resource tracking system.
type Instance struct {
	objects map[uint32]Object
	next    uint32
	mu      *sync.RWMutex
}

// Setup sets up the System.
func (s *Instance) Setup() error {
	return nil
}

// Setup tears down the System.
func (s *Instance) Teardown() {
	s.ReleaseAll()
}

// Name returns the name of the System.
func (s *Instance) Name() string {
	return SysNameInstance
}

func (s *Instance) Assign(object Object) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if object == nil {
		return ErrAssignNilObject
	}
	if object.ID() != 0 {
		return ErrObjectAlreadyAssigned
	}

	id, err := s.nextID()
	if err != nil {
		return err
	}

	s.objects[id] = object
	object.SetID(id)

	logrus.Debugf("Assigned ID %08X to %s", id, object.Name())

	return nil
}

func (s *Instance) MustAssign(object Object) {
	if err := s.Assign(object); err != nil {
		panic(err)
	}
}

func (s *Instance) Release(ids ...uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, v := range ids {
		if v == 0 {
			continue
		}

		if _, ok := s.objects[v]; !ok {
			logrus.Error(ErrIDNotFound(v))
			continue
		}

		if s.objects[v] == nil {
			logrus.Warnf("Attempted to release nil object %08X", v)
		} else {
			s.objects[v].Dealloc()
			s.objects[v].Release()
		}

		delete(s.objects, v)

		logrus.Debugf("Released ID %08X", v)
	}
}

func (s *Instance) ReleaseAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for v := range s.objects {
		if _, ok := s.objects[v]; !ok {
			logrus.Error(ErrIDNotFound(v))
			continue
		}

		if s.objects[v] == nil {
			logrus.Warnf("Attempted to release nil object %08X", v)
		} else {
			s.objects[v].Dealloc()
			s.objects[v].Release()
		}

		delete(s.objects, v)

		logrus.Debugf("Released ID %08X", v)
	}
}

func (s *Instance) nextID() (uint32, error) {

	if len(s.objects) >= math.MaxUint32 {
		return 0, ErrMaxIDsExceeded
	}

	id := s.next + 1
	_, ok := s.objects[id]

	for ok {
		id := s.next + 1
		_, ok = s.objects[id]
	}

	s.next = id

	return s.next, nil
}

func (s *Instance) Get(id uint32) (Object, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	object, ok := s.objects[id]
	if !ok {
		return nil, ErrIDNotFound(id)
	}

	return object, nil
}

// NewInstance creates a new instance system.
func NewInstance() *Instance {
	s := &Instance{
		objects: make(map[uint32]Object),
		mu:      &sync.RWMutex{},
	}

	return s
}

// GetInstance gets the instance system from the current app.
func GetInstance() *Instance {
	return CurrentApp().MustSystem(SysNameInstance).(*Instance)
}
