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
	"os"
	"path"
	"runtime"
)

var defaultUserDir = path.Join(os.Getenv("HOME"), ".gospeccy")
var distDir = path.Join(runtime.GOROOT(), "pkg", runtime.GOOS + "_" + runtime.GOARCH, "gospeccy")

func searchForValidPath(paths []string) string {
	
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// Return a valid path for the named snapshot.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/sna/
func SnaPath(fileName string) string {
	var (
		currDir = path.Join(fileName)
		userDir = path.Join(defaultUserDir, "sna", fileName)
	)
	
	path := searchForValidPath([]string{currDir, userDir})

	if path == "" {
		return fileName
	}

	return path
}
 
// Return a valid path for the 48k system ROM.
//
// The search is performed in this order:
// 1. ./roms/48.rom
// 2. $HOME/.gospeccy/roms/48.rom
// 3. $GOROOT/pkg/$GOOS_$GOARCH/gospeccy/roms/48.rom
func SystemRomPath(fileName string) string {
	var (
		currDir = path.Join(fileName)
		userDir = path.Join(defaultUserDir, "roms", fileName)
		distDir = path.Join(distDir, "roms", fileName)
	)

	path := searchForValidPath([]string{currDir, userDir, distDir})

	if path == "" {
		return fileName
	}

	return path
}

// Return a valid path for the named script.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/scripts/
// 3. $GOROOT/pkg/$GOOS_$GOARCH/gospeccy/scripts
func ScriptPath(fileName string) string {
	var (
		currDir = path.Join(fileName)
		userDir = path.Join(defaultUserDir, "scripts", fileName)
		distDir = path.Join(distDir, "scripts", fileName)
	)
	
	path := searchForValidPath([]string{currDir, userDir, distDir})

	if path == "" {
		return fileName
	}

	return path
}
