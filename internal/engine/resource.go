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
	"bufio"
	"bytes"
	"io"
	"path/filepath"
	"strings"
)

const (
	bindataPrefix = "<builtin>:"
)

type ResourceType int

const (
	ResourceFile    ResourceType = iota // ResourceFile is a file located on the local filesystem.
	ResourcePackage                     // ResourcePackage is a file located in a package.
	ResourceBindata                     // ResourceBindata is a file built in to the binary.
)

// Resource is a represents a read-only file that has an added layer of abstraction
// in terms of underlying storage type. The resource itself does not know how to
// read from the path provided, that is left up to a separate resource manager.
type Resource struct {
	resType   ResourceType
	buffer    *bytes.Buffer
	location  string
	container string
}

// NewResource creates a new Resource object for the given filename. The type
// of the resource will be derived from its path.
func NewResource(filename string) (*Resource, error) {
	r := &Resource{
		buffer: bytes.NewBuffer([]byte{}),
	}

	// Detect resource type.
	if strings.HasPrefix(filename, bindataPrefix) {
		r.resType = ResourceBindata
		r.location = r.Path(strings.TrimPrefix(filename, bindataPrefix))
		r.container = "<builtin>"
	} else if IsPackagePath(filename) {
		r.resType = ResourcePackage
		r.container, r.location = SplitPackagePath(filename)
		r.location = r.Path(r.location)
	} else {
		r.resType = ResourceFile
		r.location = filename
	}

	return r, nil
}

// Reader returns a new io.Reader for this Resource.
func (r *Resource) Reader() io.Reader {
	return bufio.NewReader(r.buffer)
}

// Bytes returns a byte slice representation of the Resource.
func (r *Resource) Bytes() []byte {
	return r.buffer.Bytes()
}

// Size returns the byte count of the Resource.
func (r *Resource) Size() int {
	return r.buffer.Len()
}

// Location returns the full path for the Resource.
func (r *Resource) Location() string {
	return r.location
}

// Container returns the name of the object containing this resource. For bindata
// resources, this is "<builtin>". For package resources, this is the name of the
// package. For file resources, an empty string is returned.
func (r *Resource) Container() string {
	return r.container
}

// Type returns the resource type.
func (r *Resource) Type() ResourceType {
	return r.resType
}

// Base returns the last element of the resource's location (the filename).
func (r *Resource) Base() string {
	return filepath.Base(r.location)
}

// Dir returns all but the last element of the resource's location.
func (r *Resource) Dir() string {
	return r.Path(filepath.Dir(r.location))
}

// Dir returns all but the last element of the resource's location.
func (r *Resource) DirPrefix() string {
	if r.resType == ResourceBindata || r.resType == ResourcePackage {
		return r.container + ":" + r.Dir()
	}

	return r.Dir()
}

func (r *Resource) Path(value string) string {
	if r.resType == ResourceFile {
		return value
	}

	if !strings.HasPrefix(value, "..") {
		if strings.HasPrefix(value, ".") {
			value = strings.TrimPrefix(value, ".")
		}
	}

	value = strings.Replace(value, "\\", "/", -1)

	if strings.HasPrefix(value, "/") {
		value = strings.TrimPrefix(value, "/")
	}

	return value
}
