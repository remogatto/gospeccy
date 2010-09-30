GOSPECCY_ROOT=../..
include $(GOSPECCY_ROOT)/Make.inc-head

TARG=spectrum
GOFILES=\
	z80.go\
	z80_gen.go\
	z80_tables.go\
	opcodes_gen.go\
	memory.go\
	port.go\
	ula.go\
	display.go\
	sdldisplay.go\
	sound.go\
	sdlsound.go\
	keyboard.go\
	spectrum.go\
	helpers.go\
	application.go\
	event.go

include $(GOROOT)/src/Make.pkg
include $(GOSPECCY_ROOT)/Make.inc-tail

sdldisplay_test.go: $(FORMATS_PKG_LIB)
z80_test.go: $(FORMATS_PKG_LIB)
spectrum_test.go: $(PRETTYTEST_PKG_LIB)