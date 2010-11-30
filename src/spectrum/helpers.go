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
	"container/vector"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"
)

var DefaultUserDir = path.Join(os.Getenv("HOME"), ".gospeccy")
var distDir = path.Join(runtime.GOROOT(), "pkg", runtime.GOOS+"_"+runtime.GOARCH, "gospeccy")

var customSearchPaths vector.StringVector
var customSearchPaths_mutex sync.RWMutex

func AddCustomSearchPath(path string) {
	customSearchPaths_mutex.Lock()
	{
		customSearchPaths.Push(path)
	}
	customSearchPaths_mutex.Unlock()
}

func searchForValidPath(paths []string, fileName string) string {
	for _, dir := range paths {
		fullPath := path.Join(dir, fileName)
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}

	return fileName
}

func appendCustomSearchPaths(paths *vector.StringVector) {
	customSearchPaths_mutex.RLock()
	{
		paths.AppendVector(&customSearchPaths)
	}
	customSearchPaths_mutex.RUnlock()
}

// Return a valid path for the named snapshot,
// or the original filename of the search fails.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/sna/
func SnaPath(fileName string) string {
	var (
		currDir = ""
		userDir = path.Join(DefaultUserDir, "sna")
	)

	var paths vector.StringVector
	paths.Push(currDir)
	paths.Push(userDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the named tape file,
// or the original filename of the search fails.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/tape/
func TapePath(fileName string) string {
	var (
		currDir = ""
		userDir = path.Join(DefaultUserDir, "tape")
	)

	var paths vector.StringVector
	paths.Push(currDir)
	paths.Push(userDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the named zip file,
// or the original filename of the search fails.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/zip/
func ZipPath(fileName string) string {
	var (
		currDir = ""
		userDir = path.Join(DefaultUserDir, "zip")
	)

	var paths vector.StringVector
	paths.Push(currDir)
	paths.Push(userDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the file based on its extension,
// or the original filename of the search fails.
func ProgramPath(fileName string) string {
	ext := strings.ToLower(path.Ext(fileName))

	switch ext {
	case ".sna", ".z80":
		return SnaPath(fileName)

	case ".tap":
		return TapePath(fileName)

	case ".zip":
		return ZipPath(fileName)
	}

	return fileName
}

// Returns a valid path for the 48k system ROM,
// or the original filename of the search fails.
//
// The search is performed in this order:
// 1. ./roms/48.rom
// 2. $HOME/.gospeccy/roms/48.rom
// 3. $GOROOT/pkg/$GOOS_$GOARCH/gospeccy/roms/48.rom
func SystemRomPath(fileName string) string {
	var (
		currDir = "roms"
		userDir = path.Join(DefaultUserDir, "roms")
		distDir = path.Join(distDir, "roms")
	)

	var paths vector.StringVector
	paths.Push(currDir)
	paths.Push(userDir)
	paths.Push(distDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}

// Return a valid path for the named script,
// or the original filename of the search fails.
//
// The search is performed in this order:
// 1. ./scripts/
// 2. $HOME/.gospeccy/scripts/
// 3. $GOROOT/pkg/$GOOS_$GOARCH/gospeccy/scripts/
func ScriptPath(fileName string) string {
	var (
		currDir = "scripts"
		userDir = path.Join(DefaultUserDir, "scripts")
		distDir = path.Join(distDir, "scripts")
	)

	var paths vector.StringVector
	paths.Push(currDir)
	paths.Push(userDir)
	paths.Push(distDir)
	appendCustomSearchPaths(&paths)

	return searchForValidPath(paths, fileName)
}


// Panic if condition is false
func assert(condition bool) {
	if !condition {
		panic("internal error")
	}
}
