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
	for {
		select {
		case <-evtLoop.Pause:
			// Remove all enqueued AudioData objects
			removed := true
			for removed {
				_, removed = <-audio.playback
			}

			close(audio.playback)
			sdl_audio.CloseAudio()

			<-audio.playbackLoopFinished

			evtLoop.Pause <- 0

		case <-evtLoop.Terminate:
			// Terminate this Go routine
			if evtLoop.App().Verbose {
				PrintfMsg("audio forwarder loop: exit")
			}
			evtLoop.Terminate <- 0
			return

		case audioData := <-audio.data:
			if audioData != nil {
				audio.bufferAdd()
				audio.playback <- audioData
			} else {
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
		PrintfMsg("audio playback loop: exit")
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
	playbackLoopFinished chan byte

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

	numSamples_cummulativeFraction float

	mutex sync.Mutex
}

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
			PrintfMsg("%#v", spec)
		}
	}

	audio := &SDLAudio{
		data:                 make(chan *AudioData),
		playback:             make(chan *AudioData, 2*BUFSIZE_IDEAL), // Use a buffered Go channel
		playbackLoopFinished: make(chan byte),
		sdlAudioUnpaused:     false,
		bufSize:              0,
		freq:                 uint(spec.Freq),
		virtualFreq:          uint(spec.Freq),
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
	audio.data <- nil
}

// Called when the number of buffered 'AudioData' objects increases by 1
func (audio *SDLAudio) bufferAdd() {
	audio.mutex.Lock()
	{
		audio.bufSize++

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
			audio.virtualFreq = uint(float(audio.virtualFreq) * 1.0005)
			changedFreq = true
		} else if audio.bufSize > BUFSIZE_IDEAL+2 {
			// Prevent future buffer overruns
			audio.virtualFreq = uint(float(audio.virtualFreq) / 1.0005)
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

func (audio *SDLAudio) render(audioData *AudioData) {
	var lastEvent_orNil *BeeperEvent = audioData.beeperEvents

	// Determine the number of beeper-events
	numEvents := 0
	for e := lastEvent_orNil; e != nil; e = e.previous_orNil {
		numEvents++
	}

	type simplifiedBeeperEvent_t struct {
		tstate uint
		level  byte
	}

	// Create an array called 'events' and initialize it with
	// the events sorted by T-state value in *ascending* order
	events := make([]simplifiedBeeperEvent_t, numEvents+1)
	var tstate_max1 uint = 0
	{
		i := numEvents - 1
		for e := lastEvent_orNil; e != nil; e = e.previous_orNil {
			events[i] = simplifiedBeeperEvent_t{e.tstate, e.level}
			i--
		}
		// At this point: 'i' should equal to -1

		// The [beeper-level from the last event] lasts until the end of the frame
		if lastEvent_orNil != nil {
			events[numEvents] = simplifiedBeeperEvent_t{TStatesPerFrame, lastEvent_orNil.level}

			// Make sure 'events[numEvents].tstate' is greater than 'events[numEvents-1].tstate'
			if (numEvents > 0) && !(events[numEvents].tstate > events[numEvents-1].tstate) {
				events[numEvents].tstate = events[numEvents-1].tstate + 1
			}

			tstate_max1 = events[numEvents].tstate
		}
	}

	// Note: If 'lastEvent_orNil' is nil, then 'event[numEvents]' is also nil. But this is OK.

	var numSamples uint
	{
		audio.mutex.Lock()

		numSamples_float := float(audio.virtualFreq) / float(audioData.fps)
		numSamples = uint(numSamples_float)

		audio.numSamples_cummulativeFraction += numSamples_float - float(numSamples)
		if audio.numSamples_cummulativeFraction >= 1.0 {
			numSamples += 1
			audio.numSamples_cummulativeFraction -= 1.0
		}

		audio.mutex.Unlock()
	}

	var k float = float(numSamples) / float(tstate_max1)

	samples := make([]float, numSamples+1)
	for i := 0; i < numEvents; i++ {
		start := events[i]
		end := events[i+1]

		if start.level != 0 {
			var position0 float = float(start.tstate) * k
			var position1 float = float(end.tstate) * k

			pos0 := uint(position0)
			pos1 := uint(position1)

			if pos0 == pos1 {
				samples[pos0] += 0x7fff * (position1 - position0)
			} else {
				samples[pos0] += 0x7fff * (float(pos0+1) - position0)
				for p := pos0 + 1; p < pos1; p++ {
					samples[p] = 0x7fff
				}
				samples[pos1] += 0x7fff * (position1 - float(pos1))
			}
		}
	}

	samples_int16 := make([]int16, numSamples)
	for i := uint(0); i < numSamples; i++ {
		s := uint(samples[i])
		if s > 0x7fff {
			s = 0x7fff
		}
		samples_int16[i] = int16(s)
	}

	sdl_audio.SendAudio_int16(samples_int16)
}
