// +build linux freebsd

package sdl_output

type InitialSettings struct {
	scale2x            *bool
	fullscreen         *bool
	showPaintedRegions *bool

	audio     *bool
	audioFreq *uint
	hqAudio   *bool
}

func (s *InitialSettings) Terminated() bool {
	return false
}

func (s *InitialSettings) ResizeVideo(scale2x, fullscreen bool) {
	// Overwrite the command-line settings
	*s.scale2x = scale2x
	*s.fullscreen = fullscreen
}

func (s *InitialSettings) ShowPaintedRegions(enable bool) {
	*s.showPaintedRegions = enable
}

func (s *InitialSettings) EnableAudio(enable bool) {
	// Overwrite the command-line settings
	*s.audio = enable
}

func (s *InitialSettings) SetAudioFreq(freq uint) {
	// Overwrite the command-line settings
	*s.audioFreq = freq
}

func (s *InitialSettings) SetAudioQuality(hqAudio bool) {
	// Overwrite the command-line settings
	*s.hqAudio = hqAudio
}
