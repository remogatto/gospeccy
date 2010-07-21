package spectrum

import (
	"testing"
	"time"
)

func loopFor(timeInterval int64, block func(elapsedTime int64)) {
	var elapsedTime int64
	startTime := time.Nanoseconds()
	for elapsedTime < timeInterval {
		block(elapsedTime)
		elapsedTime = time.Nanoseconds() - startTime
	}
}

func TestFpsCounter(t *testing.T) {
	// Collect timings every second
	var timeInterval int64 = 1 * second
	// Create a new FpsCounter service with the given timeInterval
	fpsCounter := NewFpsCounter(timeInterval)

	// Simulate a three seconds emulator loop
	loopFor(3*second, func(elapsedTime int64) {

		fpsCounter.Timings <- 20*ms // Send a dummy time of
		// 20ms (50 fps)

		fps := <-fpsCounter.Fps // Receive fps from the
		// FpsCounter service

		// After 2 seconds average fps should be 50
		if elapsedTime > 2*second {
			if fps != 50 {
				t.Errorf("fps should be 50 but got %f", fps)
			}
		}
	})
}
