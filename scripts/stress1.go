// A simple display stress-test.
// If GoSpeccy sound is enabled, there shouldn't be any sound buffer underruns.
func stress1() {
	for i:=0; i<20; i++ {
		println("Setting scale to 1")
		scale(1)
		wait(10)

		println("Setting scale to 2")
		scale(2)
		wait(10)
	}
}
