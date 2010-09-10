include $(GOROOT)/src/Make.inc

PATH_FILE=src/spectrum/path_gen_$(GOOS)_$(GOARCH).go

SPECTRUM_FILES=\
	src/spectrum/application.go\
	src/spectrum/console.go\
	src/spectrum/display.go\
	src/spectrum/keyboard.go\
	src/spectrum/memory.go\
	src/spectrum/opcodes_gen.go\
	src/spectrum/port.go\
	src/spectrum/sdldisplay.go\
	src/spectrum/sdlsound.go\
	src/spectrum/sound.go\
	src/spectrum/spectrum.go\
	src/spectrum/ula.go\
	src/spectrum/z80.go\
	src/spectrum/z80_gen.go\
	src/spectrum/z80_tables.go\


PERF_PKG_LIB=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)/⚛perf.a

READLINE_FILES=src/readline/readline.go
READLINE_ARCHIVE=src/readline/_obj/⚛readline.a
READLINE_PKG_LIB=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)/⚛readline.a

SDL_PKG_LIB=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)/⚛sdl.a
SDL_AUDIO_PKG_LIB=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)/⚛sdl/audio.a


PKG_LIBS=\
	$(PERF_PKG_LIB)\
	$(READLINE_PKG_LIB)\
	$(SDL_PKG_LIB)\
	$(SDL_AUDIO_PKG_LIB)


GOFMT_FILES=\
	src/gospeccy.go\
	src/readline/readline.go\
	src/spectrum/application.go\
	src/spectrum/console.go\
	src/spectrum/keyboard.go\
	src/spectrum/memory.go\
	src/spectrum/port.go\
	src/spectrum/sdlsound.go\
	src/spectrum/sound.go\
	src/spectrum/spectrum.go\
	src/spectrum/ula.go\
	src/spectrum/z80*.go\


gospeccy: _obj _obj/spectrum.a _obj/gospeccy.$(O) $(PKG_LIBS)
	$(LD) -L./_obj -o $@ _obj/gospeccy.$(O)

.PHONY: clean
clean:
	rm -f gospeccy
	rm -f src/spectrum/path_gen_*.go
	rm -rf _obj
	make -C src/readline clean

.PHONY: gofmt
gofmt:
	gofmt -w -l $(GOFMT_FILES)

_obj:
	mkdir _obj

_obj/gospeccy.$(O): src/gospeccy.go _obj/spectrum.a
	$(GC) -I./_obj -o $@ src/gospeccy.go

_obj/spectrum.a: $(SPECTRUM_FILES) $(PATH_FILE) $(PKG_LIBS)
	$(GC) -I./_obj -o _obj/spectrum.$(O) $(SPECTRUM_FILES) $(PATH_FILE)
	gopack grc $@ _obj/spectrum.$(O)

$(PATH_FILE):
	@echo Generating $@
	@echo "// Automatically generated file - DO NOT EDIT" >> $(PATH_FILE)
	@echo "package spectrum" >> $(PATH_FILE)
	@echo "const GOOS = \"$(GOOS)\"" >> $(PATH_FILE)
	@echo "const GOARCH = \"$(GOARCH)\"" >> $(PATH_FILE)


#
# Installation of external dependencies and internal libraries
#

$(PERF_PKG_LIB):
	goinstall -u github.com/0xe2-0x9a-0x9b/Go-PerfEvents || exit 0
	make -C $(GOROOT)/src/pkg/github.com/0xe2-0x9a-0x9b/Go-PerfEvents clean
	make -C $(GOROOT)/src/pkg/github.com/0xe2-0x9a-0x9b/Go-PerfEvents install

$(READLINE_PKG_LIB): $(READLINE_ARCHIVE)
	make -C src/readline clean
	make -C src/readline install

$(READLINE_ARCHIVE): $(READLINE_FILES)
	make -C src/readline

$(SDL_PKG_LIB):
	goinstall -u github.com/0xe2-0x9a-0x9b/Go-SDL || exit 0
	make -C $(GOROOT)/src/pkg/github.com/0xe2-0x9a-0x9b/Go-SDL clean
	make -C $(GOROOT)/src/pkg/github.com/0xe2-0x9a-0x9b/Go-SDL install

$(SDL_AUDIO_PKG_LIBS): $(SDL_PKG_LIB)
