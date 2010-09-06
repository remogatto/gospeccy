/*
 * Copyright: ⚛ <0xe2.0x9a.0x9b@gmail.com> 2010
 *
 * The contents of this file can be used freely,
 * except for usages in immoral contexts.
 */

package spectrum

// This is the primary structure for sending audio data
// from the Z80 CPU emulation core to an audio device.
type AudioData struct {
	// The FPS (frames per second) value that applies to this AudioData object
	fps float

	beeperEvents *BeeperEvent // Might be nil
}

// Interface to an audio device awaiting audio data
type AudioReceiver interface {
	getAudioDataChannel() chan<- *AudioData

	// Closes the audio device associated with this AudioReceiver
	close()
}
