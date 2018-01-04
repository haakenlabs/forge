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

import "fmt"

// Object represents a generic resource that should be tracked by the instance
// database. All resources requiring tracking should implement this interface.
type Object interface {
	// ID returns the instance ID of this object.
	ID() uint32

	// Name returns the name of this object.
	Name() string

	// SetID sets the instance ID of this object. By default, an object's ID will
	// be zero. Once the ID has been set, it cannot be changed.
	SetID(uint32)

	// SetName sets the name of this object.
	SetName(string)

	// Alloc allocates any resources during object initialization. By default,
	// this function does nothing. This function will be called automatically
	// when the object is registered with the instance database.
	Alloc() error

	// Dealloc de-allocates any resources during object teardown. By default,
	// this function does nothing. This function will be called automatically
	// when the object is un-registered with the instance database.
	Dealloc()

	// Release will set the instance ID of this object to 0.
	Release()
}

// Object is a compliant implementation of the Object interface. All types that
// intend to implement that interface should embed this struct.
type BaseObject struct {
	id   uint32
	name string
}

// ID returns the instance ID of this object.
func (o *BaseObject) ID() uint32 {
	return o.id
}

// Name returns the name of this object.
func (o *BaseObject) Name() string {
	return o.name
}

// SetID sets the instance ID of this object. By default, an object's ID will
// be zero. Once the ID has been set, it cannot be changed.
func (o *BaseObject) SetID(value uint32) {
	if o.id == 0 {
		o.id = value
	}
}

// SetName sets the name of this object.
func (o *BaseObject) SetName(value string) {
	o.name = value
}

// Alloc allocates any resources during object initialization. By default,
// this function does nothing. This function will be called automatically
// when the object is registered with the instance database.
func (o *BaseObject) Alloc() error {
	return nil
}

// Dealloc de-allocates any resources during object teardown. By default,
// this function does nothing. This function will be called automatically
// when the object is un-registered with the instance database.
func (o *BaseObject) Dealloc() {}

func (o *BaseObject) String() string {
	return fmt.Sprintf("Object(%s %08X)", o.name, o.id)
}

// Release will set the instance ID of this object to 0.
func (o *BaseObject) Release() {
	o.id = 0
}
