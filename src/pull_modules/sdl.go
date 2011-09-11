package pull_modules

import "spectrum/output/sdl_output"

func init() {
	go sdl_output.Main()
}
