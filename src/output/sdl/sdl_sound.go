/*
 * Copyright: ⚛ <0xe2.0x9a.0x9b@gmail.com> 2010
 *
 * The contents of this file can be used freely,
 * except for usages in immoral contexts.
 */

// +build linux freebsd

package sdl_output

import (
	"errors"
	"fmt"
	"github.com/scottferg/Go-SDL/sdl"
	sdl_audio "github.com/scottferg/Go-SDL/sdl/audio"
	"github.com/remogatto/gospeccy/src/spectrum"
	"math"
	"os"
	"sync"
)

func init() {
	const expectedVersion = "⚛SDL audio bindings 1.0"
	actualVersion := sdl_audio.GoSdlAudioVersion()
	if actualVersion != expectedVersion {
		fmt.Fprintf(os.Stderr, "Invalid SDL audio bindings version: expected \"%s\", got \"%s\"\n",
			expectedVersion, actualVersion)
		os.Exit(1)
	}
}

// ======================
// Audio loop (goroutine)
// ======================

// Forward 'AudioData' objects from 'audio.data' to 'audio.playback'
func forwarderLoop(evtLoop *spectrum.EventLoop, audio *SDLAudio) {
	audioDataChannel := audio.data
	playback_closed := false

	shutdown.Add(1)
	for {
		select {
		case <-evtLoop.Pause:
			// Remove all enqueued AudioData objects
		loop:
			for {
				select {
				case <-audio.playback:
				default:
					break loop
				}
			}

			audio.mutex.Lock()
			{
				if !audio.sdlAudioUnpaused {
					// Unpause SDL Audio. This is needed in order to avoid
					// a potential deadlock on 'sdl_audio.SendAudio_int16()'.
					// (If audio is paused, 'sdl_audio.SendAudio_int16()' waits indefinitely.)
					sdl_audio.PauseAudio(false)
					audio.sdlAudioUnpaused = true
				}
			}
			audio.mutex.Unlock()

			close(audio.playback)
			playback_closed = true

			<-audio.playbackLoopFinished

			sdl_audio.CloseAudio()

			audio.mutex.Lock()
			forwarderLoopFinished := audio.forwarderLoopFinished
			audio.mutex.Unlock()
			if forwarderLoopFinished != nil {
				forwarderLoopFinished <- 0
			}

			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if evtLoop.App().Verbose {
				evtLoop.App().PrintfMsg("audio forwarder loop: exit")
			}
			evtLoop.Terminate <- 0
			shutdown.Done()
			return

		case audioData := <-audioDataChannel:
			if audioData != nil {
				if !playback_closed {
					audio.bufferAdd()
					audio.playback <- audioData
				}
			} else {
				// Prevent any future sends via the 'audio.data' channel
				close(audio.data)

				// Replace 'audioDataChannel' with nil,
				// so that future executions of the 'select' statement ignore the "<-audioDataChannel" case
				audioDataChannel = nil

				// Go to the '<-evtLoop.Pause' case
				done := evtLoop.Delete()
				go func() { <-done }()
			}
		}
	}
}

func playbackLoop(app *spectrum.Application, audio *SDLAudio) {
	for audioData := range audio.playback {
		audio.bufferRemove()
		audio.render(audioData)
	}

	if app.Verbose {
		app.PrintfMsg("audio playback loop: exit")
	}
	audio.playbackLoopFinished <- 0
}

// ========
// SDLAudio
// ========

// The default SDL-audio playback frequency
const PLAYBACK_FREQUENCY = 48000

// Ideal number of buffered 'AudioData' objects,
// in order to prevent [SDL buffer underruns] and [Go channel overruns].
const BUFSIZE_IDEAL = 3

const FREQUENCY_CHANGE_RATE = 1.0002

// The function 'bufferRemove' requires a sufficiently high frequency
// so that FREQUENCY_CHANGE_RATE has an actual impact on the frequency.
// In other words: ((FREQUENCY_CHANGE_RATE-1) * MIN_PLAYBACK_FREQUENCY) has to be greater than 1.
const MIN_PLAYBACK_FREQUENCY = 10000

// The ZX Spectrum beeper has only two levels: 0 and 1.
// However, the beeper can produce multi-channel sound if it changes so quickly that
// the speaker (speaker = the physical object that the TV uses to produce the sound)
// is unable to keep up with the rate of changes between the levels. This effectively
// puts the speaker into a position somewhere *between* 0 and 1 (e.g: 0.25, 0.7)
//
// The RESPONSE_FREQUENCY is a frequency at which the multi-channel music in
// certain ZX Spectrum games and demos sounds "good-enough" to the human ear.
//
// It is used only when 'hqAudio' is enabled.
const RESPONSE_FREQUENCY = 12000

type SDLAudio struct {
	// Synchronous Go channel for receiving 'AudioData' objects
	data chan *spectrum.AudioData

	// A buffer with a capacity for multiple 'AudioData' objects.
	// The number of enqueued messages hovers around 'BUFSIZE_IDEAL'.
	playback chan *spectrum.AudioData

	// A channel for properly synchronizing the audio shutdown procedure
	playbackLoopFinished  chan byte
	forwarderLoopFinished chan byte

	// Whether SDL playback is active. Initial value is 'false'.
	// Changed to 'true' after the first 'AudioData' object becomes available.
	sdlAudioUnpaused bool

	// The number of 'AudioData' objects currently enqueued in the 'playback' Go channel
	bufSize uint

	// The playback frequency of the SDL audio device
	freq uint

	// The virtual/effective playback frequency.
	// This frequency is automatically adjusted so that the number
	// of 'AudioData' objects enqueued in the 'playback' Go channel
	// hovers around 'BUFSIZE_IDEAL'.
	virtualFreq uint

	// Sum of fractions which were lost because of integer truncation
	numSamples_cummulativeFraction float32

	// Array for storing samples. It is declared here in order
	// to avoid repetitive allocation of this array in method 'render'.
	samples []float64

	// Array for storing samples. It is declared here in order
	// to avoid repetitive allocation of this array in method 'render'.
	samples_int16 []int16

	// Overflow from previous frame. It is used if 'hqAudio' is enabled.
	overflow []float64

	// Enables higher-quality audio resampling
	hqAudio bool

	// The number of frames seen by this 'SDLAudio' object
	frame uint

	mutex sync.Mutex
}

var sdlAudio_instance *SDLAudio = nil

// Opens SDL audio.
// If 'playbackFrequency' is 0, the frequency will be equivalent to PLAYBACK_FREQUENCY.
func NewSDLAudio(app *spectrum.Application, playbackFrequency uint, hqAudio bool) (*SDLAudio, error) {
	if playbackFrequency == 0 {
		playbackFrequency = PLAYBACK_FREQUENCY
	}

	if playbackFrequency < MIN_PLAYBACK_FREQUENCY {
		return nil, errors.New(fmt.Sprintf("playback frequency of %d Hz is too low", playbackFrequency))
	}

	// Open SDL audio
	var spec sdl_audio.AudioSpec
	{
		spec.Freq = int(playbackFrequency)
		spec.Format = sdl_audio.AUDIO_S16SYS
		spec.Channels = 1
		spec.Samples = uint16(2048 * float32(playbackFrequency) / PLAYBACK_FREQUENCY)
		if sdl_audio.OpenAudio(&spec, &spec) != 0 {
			return nil, errors.New(sdl.GetError())
		}
		if app.Verbose {
			app.PrintfMsg("%#v", spec)
		}
	}

	audio := &SDLAudio{
		data:                  make(chan *spectrum.AudioData),
		playback:              make(chan *spectrum.AudioData, 2*BUFSIZE_IDEAL), // Use a buffered Go channel
		playbackLoopFinished:  make(chan byte),
		forwarderLoopFinished: nil,
		sdlAudioUnpaused:      false,
		bufSize:               0,
		freq:                  uint(spec.Freq),
		virtualFreq:           uint(spec.Freq),
		hqAudio:               hqAudio,
	}

	go forwarderLoop(app.NewEventLoop(), audio)
	go playbackLoop(app, audio)

	return audio, nil
}

// Implement AudioReceiver
func (audio *SDLAudio) GetAudioDataChannel() chan<- *spectrum.AudioData {
	return audio.data
}

func (audio *SDLAudio) Close() {
	audio.mutex.Lock()
	audio.forwarderLoopFinished = make(chan byte)
	audio.mutex.Unlock()

	audio.data <- nil

	<-audio.forwarderLoopFinished

	audio.mutex.Lock()
	audio.forwarderLoopFinished = nil
	audio.mutex.Unlock()
}

// Called when the number of buffered 'AudioData' objects increases by 1
func (audio *SDLAudio) bufferAdd() {
	audio.mutex.Lock()
	{
		audio.bufSize++

		// Unpause SDL audio if we have BUFSIZE_IDEAL 'AudioData' objects
		if !audio.sdlAudioUnpaused && (audio.bufSize == BUFSIZE_IDEAL) {
			sdl_audio.PauseAudio(false)
			audio.sdlAudioUnpaused = true
		}
	}
	audio.mutex.Unlock()
}

// Called when the number of buffered 'AudioData' objects decreases by 1
func (audio *SDLAudio) bufferRemove() {
	audio.mutex.Lock()
	{
		audio.bufSize--

		changedFreq := false
		if audio.bufSize < BUFSIZE_IDEAL-2 {
			// Prevent future buffer underruns
			audio.virtualFreq = uint(float32(audio.virtualFreq) * FREQUENCY_CHANGE_RATE)
			changedFreq = true
		} else if audio.bufSize > BUFSIZE_IDEAL+2 {
			// Prevent future buffer overruns
			audio.virtualFreq = uint(float32(audio.virtualFreq) / FREQUENCY_CHANGE_RATE)
			changedFreq = true
		} else if audio.bufSize == BUFSIZE_IDEAL {
			if audio.virtualFreq != audio.freq {
				audio.virtualFreq = audio.freq
				changedFreq = true
			}
		}

		if changedFreq {
			//fmt.Printf("bufSize=%d, virtualFreq=%d\n", audio.bufSize, audio.virtualFreq)
		}
	}
	audio.mutex.Unlock()
}

func add_lq(samples []float64, x, w, h float64) {
	var position0 float64 = x
	var position1 float64 = (x + w)

	pos0 := uint(position0)
	pos1 := uint(position1)

	if pos0 == pos1 {
		volume_pos0 := h * w
		samples[pos0] += volume_pos0
	} else {
		ceil_position0 := math.Ceil(position0)

		volume_pos0 := h * (ceil_position0 - position0)
		samples[pos0] += volume_pos0

		for p := uint(ceil_position0); p < pos1; p++ {
			samples[p] += h
		}

		volume_pos1 := h * (position1 - float64(pos1))
		samples[pos1] += volume_pos1
	}
}

// A similar algorithm to 'add_lq', but in order to reduce noise
// it spreads the signal across a wider range of samples.
//
// Another difference between 'add_lq' and 'add_hq' is that the concrete
// mapping of 'x' to 'samples' in 'add_lq' is only circumstantial.
// Function 'add_hq' is not pretending that it knows [the mapping of 'x'
// to a particular element of 'samples'].
//
// Parameters:
//    spread1 = (1 / spread)
func add_hq(samples []float64, x, w, h, spread, spread1 float64) {
	// This function contains two implementations of the same algorithm.
	// The SLOW_AND_LESS_ACCURATE code-path contains a more readable version
	// of the algorithm. The other code-path has been derived from the
	// 1st code-path and uses the assumption that Z goes to +infinity.

	const SLOW_AND_LESS_ACCURATE bool = false

	if SLOW_AND_LESS_ACCURATE {
		const Z float64 = 32
		const Z1 float64 = 1 / Z

		for z := float64(0); z < 1; z += Z1 {
			var position0 float64 = x + spread*z
			var position1 float64 = position0 + w

			pos0 := int(position0)
			pos1 := int(position1)

			if pos0 == pos1 {
				volume_pos0 := h * w
				samples[pos0] += volume_pos0 * Z1
			} else {
				ceil_position0 := math.Ceil(position0)

				volume_pos0 := h * (ceil_position0 - position0)
				samples[pos0] += volume_pos0 * Z1

				for p := int(ceil_position0); p < pos1; p++ {
					samples[p] += h * Z1
				}

				volume_pos1 := h * (position1 - float64(pos1))
				samples[pos1] += volume_pos1 * Z1
			}
		}
	} else {
		var position0, position1 float64 = x, x + w
		var pos0, pos1 int = int(position0), int(position1)

		z := float64(1)
		for {
			if pos0 == pos1 {
				pos1_frac := (position1 - float64(pos1))

				numIterations := (1 - pos1_frac)
				numIterations_s := numIterations * spread1

				if numIterations_s < z {
					samples[pos0] += h * w * numIterations_s

					z -= numIterations_s
					pos1++
					position0 += numIterations
					position1 += numIterations
				} else {
					samples[pos0] += h * w * z
					break
				}
			} else {
				pos0_frac := (position0 - float64(pos0))
				pos1_frac := (position1 - float64(pos1))

				var numIterations float64
				var pos0_next, pos1_next int
				if pos0_frac > pos1_frac {
					// (1-pos0_frac) is smaller than (1-pos1_frac)
					numIterations = (1 - pos0_frac)
					pos0_next = pos0 + 1
					pos1_next = pos1
				} else if pos0_frac < pos1_frac {
					numIterations = (1 - pos1_frac)
					pos0_next = pos0
					pos1_next = pos1 + 1
				} else {
					numIterations = (1 - pos0_frac)
					pos0_next = pos0 + 1
					pos1_next = pos1 + 1
				}

				numIterations_s := numIterations * spread1

				step := 0.5 * spread * h

				if numIterations_s < z {
					a := numIterations_s * h
					b := step * numIterations_s * numIterations_s
					samples[pos0] += a*(1-pos0_frac) - b
					samples[pos1] += a*pos1_frac + b

					for p := pos0 + 1; p < pos1; p++ {
						samples[p] += a
					}

					z -= numIterations_s
					pos0 = pos0_next
					pos1 = pos1_next
					position0 += numIterations
					position1 += numIterations
				} else {
					a := z * h
					b := step * z * z
					samples[pos0] += a*(1-pos0_frac) - b
					samples[pos1] += a*pos1_frac + b

					for p := pos0 + 1; p < pos1; p++ {
						samples[p] += a
					}

					break
				}
			}
		}
	}
}

func (audio *SDLAudio) render(audioData *spectrum.AudioData) {
	var events []spectrum.BeeperEvent

	if len(audioData.BeeperEvents) > 0 {
		var firstEvent *spectrum.BeeperEvent = &audioData.BeeperEvents[0]
		spectrum.Assert(firstEvent.TState == 0)

		var lastEvent *spectrum.BeeperEvent = &audioData.BeeperEvents[len(audioData.BeeperEvents)-1]
		spectrum.Assert(lastEvent.TState == spectrum.TStatesPerFrame)

		events = audioData.BeeperEvents
	} else {
		events = make([]spectrum.BeeperEvent, 2)
		events[0] = spectrum.BeeperEvent{TState: 0, Level: 0}
		events[1] = spectrum.BeeperEvent{TState: spectrum.TStatesPerFrame, Level: 0}
	}

	/*
		// A test signal
		const D = spectrum.TStatesPerFrame/128
		events = make([]spectrum.BeeperEvent, spectrum.TStatesPerFrame/D+1)
		for i:=uint(0); i<128; i++ {
			events[i] = spectrum.BeeperEvent{TState: i*D, Level: uint8(i%2)}
		}
		events[len(events)-1] = spectrum.BeeperEvent{TState: spectrum.TStatesPerFrame, Level: 0}
	*/

	numEvents := len(events)

	spread := float64(audio.freq) / RESPONSE_FREQUENCY
	spread1 := 1 / spread

	var numSamples int
	var samples []float64
	var samples_int16 []int16
	var overflow []float64
	{
		audio.mutex.Lock()

		numSamples_float := float32(audio.virtualFreq) / audioData.FPS
		numSamples = int(numSamples_float)

		len_overflow := int(math.Ceil(spread)) + 2

		audio.numSamples_cummulativeFraction += numSamples_float - float32(numSamples)
		if audio.numSamples_cummulativeFraction >= 1.0 {
			numSamples += 1
			audio.numSamples_cummulativeFraction -= 1.0
		}

		if len(audio.samples) < numSamples+len_overflow {
			audio.samples = make([]float64, numSamples+len_overflow)
		}
		samples = audio.samples

		if len(audio.samples_int16) < numSamples {
			audio.samples_int16 = make([]int16, numSamples)
		}
		samples_int16 = audio.samples_int16

		if len(audio.overflow) < len_overflow {
			new_overflow := make([]float64, len_overflow)
			copy(new_overflow, overflow)
			audio.overflow = new_overflow
		}
		overflow = audio.overflow

		audio.mutex.Unlock()
	}

	var k float64 = float64(numSamples) / spectrum.TStatesPerFrame

	{
		for i := 0; i < len(samples); i++ {
			samples[i] = 0
		}

		copy(samples[:], overflow[:])

		for i := 0; i < numEvents-1; i++ {
			start := events[i]
			end := events[i+1]

			level := float64(spectrum.Audio16_Table[start.Level])

			var position0 float64 = float64(start.TState) * k
			var position1 float64 = float64(end.TState) * k

			if audio.hqAudio {
				add_hq(samples, position0+1, position1-position0, level, spread, spread1)
			} else {
				add_lq(samples, position0+1, position1-position0, level)
			}
		}

		copy(overflow[:], samples[numSamples:])
	}

	for i := 0; i < numSamples; i++ {
		const VOLUME_ADJUSTMENT = 0.5
		samples_int16[i] = int16(VOLUME_ADJUSTMENT * samples[i])
	}

	audio.frame++
	sdl_audio.SendAudio_int16(samples_int16[0:numSamples])
}
