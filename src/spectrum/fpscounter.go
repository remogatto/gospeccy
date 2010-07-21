package spectrum

import "time"

const (
	second = 1e9
	ms     = 1e6
)

type FpsCounter struct {
	// The channel on which timings are sent from client code
	Timings chan<- int64
	// Client code receives fps values from this channel
	Fps <-chan float
}

func NewFpsCounter(timeInterval int64) *FpsCounter {

	ticker := time.Tick(timeInterval)
	timings := make(chan int64)
	fpsOut := make(chan float)
	fps := make(chan float)

	fpsCounter := &FpsCounter{timings, fpsOut}

	// Non-blocking fps stream
	go func() {
		lastFps := 0.0
		for {
			select {
			case fpsOut <- lastFps:

			case lastFps = <-fps:
				if closed(fps) {
					close(fpsOut)
					return
				}
			}
		}
	}()

	// Calculate average fps and reset variables every tick
	go func() {
		sum := int64(0)
		numSamples := 0
		for {
			select {
			case t := <-timings:
				if closed(timings) {
					close(fps)
					return
				}
				sum += t
				numSamples++

			case <-ticker:
				if numSamples > 0 {
					avgTime := sum / int64(numSamples)
					fps <- 1/(float(avgTime)/second)
					sum, numSamples = 0, 0
				}
			}
		}
	}()

	return fpsCounter
}
