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
	FPS float32

	BeeperEvents []BeeperEvent
}

const MAX_AUDIO_LEVEL = 3

// Voltage levels on pin 28 of the ULA chip after an out to port 0xFE,
// with no input signal on the EAR socket. Issue 2 Spectrums.
// Source: http://www.worldofspectrum.org/faq/reference/48kreference.htm
//
// The array is indexed by [bits 3&4 of the value sent to the port].
var Voltage_Issue2 = [4]float32{
	0.39,
	0.73,
	3.66,
	3.79,
}

// Voltage levels on pin 28 of the ULA chip after an out to port 0xFE,
// with no input signal on the EAR socket. Issue 3 Spectrums.
// Source: http://www.worldofspectrum.org/faq/reference/48kreference.htm
//
// The array is indexed by [bits 3&4 of the value sent to the port].
var Voltage_Issue3 = [4]float32{
	0.34,
	0.66,
	3.56,
	3.70,
}

// A table for converting the "audio level" to a 16-bit signed value.
// Note: Users of this table can assume that 'Audio16_Table[0]' equals to zero.
var Audio16_Table = [4]float32{
	0,
	0x7fff * (Voltage_Issue2[1] - Voltage_Issue2[0]) / (Voltage_Issue2[3] - Voltage_Issue2[0]),
	0x7fff * (Voltage_Issue2[2] - Voltage_Issue2[0]) / (Voltage_Issue2[3] - Voltage_Issue2[0]),
	0x7fff,
}

// Interface to an audio device awaiting audio data
type AudioReceiver interface {
	GetAudioDataChannel() chan<- *AudioData

	// Closes the audio device associated with this AudioReceiver
	Close()
}
