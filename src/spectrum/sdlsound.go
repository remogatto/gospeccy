/*
 * Copyright: ⚛ <0xe2.0x9a.0x9b@gmail.com> 2010
 *
 * The contents of this file can be used freely,
 * except for usages in immoral contexts.
 */

package spectrum

import (
	"⚛sdl"
	sdl_audio "⚛sdl/audio"
	"sync"
)


const PLAYBACK_FREQUENCY = 48000


// ======================
// Audio loop (goroutine)
// ======================

// Forward 'AudioData' objects from 'audio.data' to 'audio.playback'
func forwarderLoop(evtLoop *EventLoop, audio *SDLAudio) {
	for {
		select {
		case <-evtLoop.Pause:
			close(audio.playback)
			<-audio.sdlAudioClosed

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
		PrintfMsg("audio playback loop: close SDL audio")
	}
	sdl_audio.CloseAudio()
	audio.sdlAudioClosed <- 0
}


// ========
// SDLAudio
// ========

// Ideal number of buffered 'AudioData' objects,
// in order to prevent [SDL buffer underruns] and [Go channel overruns].
const BUFSIZE_IDEAL = 2

type SDLAudio struct {
	// Synchronous Go channel for receiving 'AudioData' objects
	data chan *AudioData

	// A buffer with a capacity for multiple 'AudioData' objects.
	// The number of enqueued messages hovers around 'BUFSIZE_IDEAL'.
	playback chan *AudioData

	// A channel for properly synchronizing the audio shutdown procedure
	sdlAudioClosed chan byte

	// The number of 'AudioData' objects currently enqueued in the 'playback' Go channel.
	bufSize uint

	// The playback frequency of the SDL audio device
	freq uint

	// The virtual/effective playback frequency.
	// This frequency is being automatically adjusted so that the number
	// of 'AudioData' objects enqueued in the 'playback' Go channel
	// hovers around 'BUFSIZE_IDEAL'.
	virtualFreq uint

	mutex sync.Mutex
}

func NewSDLAudio(app *Application) *SDLAudio {
	// Open SDL audio
	var spec sdl_audio.AudioSpec
	{
		spec.Freq = PLAYBACK_FREQUENCY
		spec.Format = sdl_audio.AUDIO_S16SYS
		spec.Channels = 1
		spec.Samples = 2048
		if sdl_audio.OpenAudio(&spec, &spec) != 0 {
			panic(sdl.GetError())
		}
		if app.Verbose {
			PrintfMsg("%#v", spec)
		}
	}

	audio := &SDLAudio{
		data:           make(chan *AudioData),
		playback:       make(chan *AudioData, 2*BUFSIZE_IDEAL), // Use a buffered Go channel
		sdlAudioClosed: make(chan byte),
		bufSize:        0,
		freq:           uint(spec.Freq),
		virtualFreq:    uint(spec.Freq),
	}

	go forwarderLoop(app.NewEventLoop(), audio)
	go playbackLoop(app, audio)

	// Unpause SDL audio
	sdl_audio.PauseAudio(false)

	return audio
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
	audio.bufSize++
	audio.mutex.Unlock()
}

// Called when the number of buffered 'AudioData' objects decreases by 1
func (audio *SDLAudio) bufferRemove() {
	audio.mutex.Lock()
	{
		audio.bufSize--

		if audio.bufSize < BUFSIZE_IDEAL {
			audio.virtualFreq = uint(float(audio.virtualFreq) * 1.001)
		} else if audio.bufSize > BUFSIZE_IDEAL {
			audio.virtualFreq = uint(float(audio.virtualFreq) / 1.001)
		} else {
			audio.virtualFreq = audio.freq
		}

		//PrintfMsg("bufSize=%d, virtualFreq=%d", audio.bufSize, audio.virtualFreq)
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
		}
	}

	// Note: If 'lastEvent_orNil' is nil, then 'event[numEvents]' is also nil. But this is OK.

	audio.mutex.Lock()
	var numSamplesPerFrame float = float(audio.virtualFreq) / audioData.fps
	audio.mutex.Unlock()

	numSamples := uint(numSamplesPerFrame)
	samples := make([]int16, numSamples)
	for i := 0; i < numEvents; i++ {
		start := events[i]
		end := events[i+1]

		sample_startIndex := uint(float(start.tstate) / TStatesPerFrame * numSamplesPerFrame)
		sample_endIndex := uint(float(end.tstate) / TStatesPerFrame * numSamplesPerFrame)
		if sample_endIndex > numSamples {
			sample_endIndex = numSamples
		}

		var audio_level int16
		if start.level == 0 {
			audio_level = 0
		} else {
			audio_level = 0x7fff
		}

		for i := sample_startIndex; i < sample_endIndex; i++ {
			samples[i] = audio_level
		}
	}

	sdl_audio.SendAudio_int16(samples[0:numSamples])
}
