// A simple display stress-test.
// If GoSpeccy audio is enabled, there shouldn't be any audio buffer underruns.
func stress1(verbose bool) {
	for i:=0; i<20; i++ {
		if verbose {
			println("Setting scale to 1")
		}
		scale(1)
		wait(10)

		if verbose {
			println("Setting scale to 2")
		}
		scale(2)
		wait(10)
	}
}
