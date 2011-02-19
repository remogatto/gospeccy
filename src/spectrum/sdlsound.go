/*
 * Copyright: ⚛ <0xe2.0x9a.0x9b@gmail.com> 2010
 *
 * The contents of this file can be used freely,
 * except for usages in immoral contexts.
 */

package spectrum

import (
	"fmt"
	"os"
	"⚛sdl"
	sdl_audio "⚛sdl/audio"
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
func forwarderLoop(evtLoop *EventLoop, audio *SDLAudio) {
	audioDataChannel := audio.data
	playback_closed := false

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
				evtLoop.Delete()
			}
		}
	}
}

func playbackLoop(app *Application, audio *SDLAudio) {
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

const PLAYBACK_FREQUENCY = 48000

// Ideal number of buffered 'AudioData' objects,
// in order to prevent [SDL buffer underruns] and [Go channel overruns].
const BUFSIZE_IDEAL = 3

type SDLAudio struct {
	// Synchronous Go channel for receiving 'AudioData' objects
	data chan *AudioData

	// A buffer with a capacity for multiple 'AudioData' objects.
	// The number of enqueued messages hovers around 'BUFSIZE_IDEAL'.
	playback chan *AudioData

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
	// This frequency is being automatically adjusted so that the number
	// of 'AudioData' objects enqueued in the 'playback' Go channel
	// hovers around 'BUFSIZE_IDEAL'.
	virtualFreq uint

	// Sum of fractions which were lost because of integer truncation
	numSamples_cummulativeFraction float32

	// Array for storing samples. It is declared here in order
	// to avoid repetitive allocation of this array in method 'render'.
	samples_int16 []int16

	mutex sync.Mutex
}

var sdlAudio_instance *SDLAudio = nil

func NewSDLAudio(app *Application) (*SDLAudio, os.Error) {
	// Open SDL audio
	var spec sdl_audio.AudioSpec
	{
		spec.Freq = PLAYBACK_FREQUENCY
		spec.Format = sdl_audio.AUDIO_S16SYS
		spec.Channels = 1
		spec.Samples = 2048
		if sdl_audio.OpenAudio(&spec, &spec) != 0 {
			return nil, os.NewError(sdl.GetError())
		}
		if app.Verbose {
			app.PrintfMsg("%#v", spec)
		}
	}

	audio := &SDLAudio{
		data:                  make(chan *AudioData),
		playback:              make(chan *AudioData, 2*BUFSIZE_IDEAL), // Use a buffered Go channel
		playbackLoopFinished:  make(chan byte),
		forwarderLoopFinished: nil,
		sdlAudioUnpaused:      false,
		bufSize:               0,
		freq:                  uint(spec.Freq),
		virtualFreq:           uint(spec.Freq),
	}

	go forwarderLoop(app.NewEventLoop(), audio)
	go playbackLoop(app, audio)

	return audio, nil
}

// Implement AudioReceiver
func (audio *SDLAudio) getAudioDataChannel() chan<- *AudioData {
	return audio.data
}

func (audio *SDLAudio) close() {
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
			audio.virtualFreq = uint(float32(audio.virtualFreq) * 1.0005)
			changedFreq = true
		} else if audio.bufSize > BUFSIZE_IDEAL+2 {
			// Prevent future buffer overruns
			audio.virtualFreq = uint(float32(audio.virtualFreq) / 1.0005)
			changedFreq = true
		} else if audio.bufSize == BUFSIZE_IDEAL {
			if audio.virtualFreq != audio.freq {
				audio.virtualFreq = audio.freq
				changedFreq = true
			}
		}

		if changedFreq {
			//PrintfMsg("bufSize=%d, virtualFreq=%d", audio.bufSize, audio.virtualFreq)
		}
	}
	audio.mutex.Unlock()
}


type simplifiedBeeperEvent_t struct {
	tstate uint
	level  byte
}

type simplifiedBeeperEvent_array_t struct {
	events []simplifiedBeeperEvent_t
}

func (a *simplifiedBeeperEvent_array_t) Init(n int) {
	a.events = make([]simplifiedBeeperEvent_t, n)
}

func (a *simplifiedBeeperEvent_array_t) Set(i int, _e Event) {
	e := _e.(*BeeperEvent)
	a.events[i] = simplifiedBeeperEvent_t{e.tstate, e.level}
}


func (audio *SDLAudio) render(audioData *AudioData) {
	var events []simplifiedBeeperEvent_t

	if audioData.beeperEvents_orNil != nil {
		var lastEvent *BeeperEvent = audioData.beeperEvents_orNil
		assert(lastEvent.tstate == TStatesPerFrame)

		// Put the events in an array, sorted by T-state value in ascending order
		events_array := &simplifiedBeeperEvent_array_t{}
		EventListToArray_Ascending(lastEvent, events_array, nil)

		events = events_array.events
	} else {
		events = make([]simplifiedBeeperEvent_t, 2)
		events[0] = simplifiedBeeperEvent_t{tstate: 0, level: 0}
		events[1] = simplifiedBeeperEvent_t{tstate: TStatesPerFrame, level: 0}
	}

	numEvents := len(events)

	var numSamples uint
	var samples_int16 []int16
	{
		audio.mutex.Lock()

		numSamples_float := float32(audio.virtualFreq) / audioData.fps
		numSamples = uint(numSamples_float)

		audio.numSamples_cummulativeFraction += numSamples_float - float32(numSamples)
		if audio.numSamples_cummulativeFraction >= 1.0 {
			numSamples += 1
			audio.numSamples_cummulativeFraction -= 1.0
		}

		if len(audio.samples_int16) < int(numSamples) {
			audio.samples_int16 = make([]int16, numSamples)
		}
		samples_int16 = audio.samples_int16

		audio.mutex.Unlock()
	}

	var k float32 = float32(numSamples) / TStatesPerFrame

	samples := make([]float32, numSamples+1)
	for i := 0; i < numEvents-1; i++ {
		start := events[i]

		if start.level > 0 {
			level := Audio16_Table[start.level]
			end := events[i+1]

			var position0 float32 = float32(start.tstate) * k
			var position1 float32 = float32(end.tstate) * k

			pos0 := uint(position0)
			pos1 := uint(position1)

			if pos0 == pos1 {
				samples[pos0] += level * (position1 - position0)
			} else {
				samples[pos0] += level * (float32(pos0+1) - position0)
				for p := pos0 + 1; p < pos1; p++ {
					samples[p] = level
				}
				samples[pos1] += level * (position1 - float32(pos1))
			}
		}
	}

	for i := uint(0); i < numSamples; i++ {
		s := uint(samples[i])
		if s > 0x7fff {
			s = 0x7fff
		}
		samples_int16[i] = int16(s)
	}

	sdl_audio.SendAudio_int16(samples_int16[0:numSamples])
}
