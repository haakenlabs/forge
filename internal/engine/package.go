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
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	pkgRe        = regexp.MustCompile(`^([\w\d-_]+):([\w\d-_/.\\]+)$`)
	pkgExtension = ".pkg"
	pkgRoot      = "assets"
)

type Package struct {
	name   string
	path   string
	reader *zip.ReadCloser
}

// ErrPackageNotFound reports that package was not found/mounted.
type ErrPackageNotMounted string

func (e ErrPackageNotMounted) Error() string {
	return "fs: package not mounted: " + string(e)
}

// ErrPackageMounted reports that the package is already mounted.
type ErrPackageMounted string

func (e ErrPackageMounted) Error() string {
	return "fs: package already mounted: " + string(e)
}

type ErrPackageFileNotFound struct {
	pkg  string
	file string
}

func (e ErrPackageFileNotFound) Error() string {
	return fmt.Sprintf("fs: file '%s' in package '%s' not found", e.file, e.pkg)
}

func NewPackage(name string) *Package {
	p := &Package{
		name: name,
	}

	pkgPath := filepath.Join(pkgRoot, p.name)
	if !strings.HasSuffix(pkgPath, pkgExtension) {
		pkgPath = fmt.Sprintf("%s%s", pkgPath, pkgExtension)
	}

	p.path = pkgPath

	return p
}

func (p *Package) Mount() error {
	if p.reader != nil {
		return ErrPackageMounted(p.name)
	}

	reader, err := zip.OpenReader(p.path)
	if err != nil {
		panic(err)
	}

	p.reader = reader

	logrus.Info("Mounted package: ", p.name)

	return nil
}

func (p *Package) Unmount() error {
	err := p.reader.Close()
	p.reader = nil

	logrus.Info("Unmounted package: ", p.name)

	return err
}

func (p *Package) Name() string {
	return p.name
}

func (p *Package) Path() string {
	return p.path
}

func (p *Package) Read(filename string, w io.Writer) error {
	if p.reader == nil {
		return ErrPackageNotMounted(p.name)
	}

	var file *zip.File

	for _, f := range p.reader.File {
		if f.Name != filename {
			continue
		}

		file = f
		break
	}

	if file == nil {
		return ErrPackageFileNotFound{p.name, filename}
	}

	fReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fReader.Close()

	_, err = io.Copy(w, fReader)

	return nil
}

func IsPackagePath(filename string) bool {
	return pkgRe.MatchString(filename)
}

func SplitPackagePath(filename string) (string, string) {
	matches := pkgRe.FindAllStringSubmatch(filename, 1)
	if len(matches) != 1 {
		return "", filename
	}

	return matches[0][1], matches[0][2]
}
