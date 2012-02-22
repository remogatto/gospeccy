/*

Copyright (c) 2010 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

*/

// FIXME: The content of this file should be moved to a different
// package once we accomplish a better separation between the spectrum
// emulation basecode (package spectrum) and the frontend (package
// main). For example, it doesn't make sense to define path-related
// constants and methods into spectrum package. However, at the
// moment, these helpers are needed by both the frontend and the
// console which is part of the spectrum package. There is a lot of
// duplication too. BTW, as a first iteration we're happy.

package spectrum

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var DefaultUserDir = path.Join(os.Getenv("HOME"), ".gospeccy")
var distDir string

func init() {
	gopaths := os.Getenv("GOPATH")
	gopath0 := strings.Split(gopaths, string(os.PathListSeparator))[0]
	distDir = path.Join(gopath0, "src", "github.com", "remogatto", "gospeccy")
}

var customSearchPaths []string
var customSearchPaths_mutex sync.RWMutex

func AddCustomSearchPath(path string) {
	customSearchPaths_mutex.Lock()
	{
		customSearchPaths = append(customSearchPaths, path)
	}
	customSearchPaths_mutex.Unlock()
}

func searchForValidPath(paths []string, fileName string) (string, error) {
	for _, dir := range paths {
		if _, err := os.Lstat(dir); err == nil {
			_, err = filepath.EvalSymlinks(dir)
			if err != nil {
				return "", errors.New("path \"" + dir + "\" contains one or more invalid symbolic links")
			}
		}

		fullPath := path.Join(dir, fileName)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath, nil
		}
	}

	return fileName, nil
}

func appendCustomSearchPaths(paths *[]string) {
	customSearchPaths_mutex.RLock()
	{
		*paths = append(*paths, customSearchPaths...)
	}
	customSearchPaths_mutex.RUnlock()
}

// Return a valid path for the specified snapshot,
// or the original filename if the search did not find anything.
//
// An error is returned if the search could not proceed.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/sna/
func SnaPath(fileName string) (string, error) {
	var (
		currDir = ""
		userDir = path.Join(DefaultUserDir, "sna")
	)

	var paths []string
	paths = append(paths, currDir, userDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the specified tape file,
// or the original filename if the search did not find anything.
//
// An error is returned if the search could not proceed.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/tape/
func TapePath(fileName string) (string, error) {
	var (
		currDir = ""
		userDir = path.Join(DefaultUserDir, "tape")
	)

	var paths []string
	paths = append(paths, currDir, userDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the specified zip file,
// or the original filename if the search did not find anything.
//
// An error is returned if the search could not proceed.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/zip/
func ZipPath(fileName string) (string, error) {
	var (
		currDir = ""
		userDir = path.Join(DefaultUserDir, "zip")
	)

	var paths []string
	paths = append(paths, currDir, userDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the file based on its extension,
// or the original filename if the search did not find anything.
//
// An error is returned if the search could not proceed.
func ProgramPath(fileName string) (string, error) {
	ext := strings.ToLower(path.Ext(fileName))

	switch ext {
	case ".sna", ".z80":
		return SnaPath(fileName)

	case ".tap":
		return TapePath(fileName)

	case ".zip":
		return ZipPath(fileName)
	}

	return fileName, nil
}

// Returns a valid path for the 48k system ROM,
// or the original filename if the search did not find anything.
//
// An error is returned if the search could not proceed.
//
// The search is performed in this order:
// 1. ./roms/
// 2. $HOME/.gospeccy/roms/
// 3. $GOPATH/src/github.com/remogatto/gospeccy/roms/
func SystemRomPath(fileName string) (string, error) {
	var (
		currDir = "roms"
		userDir = path.Join(DefaultUserDir, "roms")
		distDir = path.Join(distDir, "roms")
	)

	var paths []string
	paths = append(paths, currDir, userDir, distDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the specified script,
// or the original filename if the search did not find anything.
//
// An error is returned if the search could not proceed.
//
// The search is performed in this order:
// 1. ./scripts/
// 2. $HOME/.gospeccy/scripts/
// 3. $GOPATH/src/github.com/remogatto/gospeccy/scripts/
func ScriptPath(fileName string) (string, error) {
	var (
		currDir = "scripts"
		userDir = path.Join(DefaultUserDir, "scripts")
		distDir = path.Join(distDir, "scripts")
	)

	var paths []string
	paths = append(paths, currDir, userDir, distDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the specified font file,
// or the original filename if the search did not find anything.
//
// An error is returned if the search could not proceed.
//
// The search is performed in this order:
// 1. ./fonts/
// 2. $HOME/.gospeccy/fonts/
// 3. $GOPATH/src/github.com/remogatto/gospeccy/fonts/
func FontPath(fileName string) (string, error) {
	var (
		currDir = "fonts"
		userDir = path.Join(DefaultUserDir, "fonts")
		distDir = path.Join(distDir, "fonts")
	)

	var paths []string
	paths = append(paths, currDir, userDir, distDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Reads the 16KB ROM from the specified file
func ReadROM(path string) (*[0x4000]byte, error) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(fileData) != 0x4000 {
		return nil, errors.New(path + ":invalid ROM file")
	}

	var rom [0x4000]byte
	copy(rom[:], fileData)
	return &rom, nil
}

// Panic if condition is false
func Assert(condition bool) {
	if !condition {
		panic("internal error")
	}
}
