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

const defaultUserDir = ".gospeccy"

// Return a valid path for the named snapshot.
//
// The search is performed in this order:
// 1. ./
// 2. $HOME/.gospeccy/sna/
func SnaPath(filename string) string {
	var (
		currDir = path.Join("./", filename)
		userDir = path.Join(os.Getenv("HOME"), defaultUserDir, "sna/", filename)
	)
	
	if _, err := os.Stat(filename); err == nil {
		return filename
	} else if _, err := os.Stat(currDir); err == nil {
		return currDir
 	} else if _, err := os.Stat(userDir); err == nil {
		return userDir
	} else {
		return ""
	}

	return ""
}
 
// Return a valid path for the 48k system ROM.
//
// The search is performed in this order:
// 1. ./roms/48.rom
// 2. $HOME/.gospeccy/roms/48.rom
// 3. $GOROOT/$GOOS_$GOARCH/pkg/gospeccy/roms/48.rom
func SystemRomPath() string {
	var (
		currDir = "./roms/48.rom"
		userDir = path.Join(os.Getenv("HOME"), defaultUserDir, "roms/48.rom")
		distDir = path.Join(runtime.GOROOT(), "pkg", runtime.GOOS + "_" + runtime.GOARCH, "gospeccy/roms/48.rom")
	)
	println(distDir)
	if _, err := os.Stat(currDir); err == nil {
		return currDir
	} else if _, err := os.Stat(userDir); err == nil {
		return userDir
	} else if _, err := os.Stat(distDir); err == nil {
		return distDir
	} else {
		return ""
	}

	return ""
}
