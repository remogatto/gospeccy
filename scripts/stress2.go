// +build ignore

// A simple audio stress-test.
//
// The ideal behavior is:
//  - there are no deadlocks
//  - timing of animations remains the same
//    (animation smoothness may change a little bit because of missed frames
//     and because of altered X server display synchronization caused by messages
//     being printed in the terminal)
func stress2(verbose bool) {
	for i:=0; i<20; i++ {
		if verbose {
			puts("Disabling audio\n")
		}
		audio(false)
		wait(10)

		if verbose {
			puts("Enabling audio\n")
		}
		audio(true)
		wait(10)
	}
}
