// +build linux freebsd

package pull_modules

import (
	"github.com/remogatto/gospeccy/src/output/sdl"
)

func init() {
	go sdl_output.Main()
}
