GOSPECCY_ROOT=.
include Make.inc-head

DIST_PATH=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)/gospeccy
SYSTEM_ROM_PATH=$(DIST_PATH)/roms
SCRIPTS_PATH=$(DIST_PATH)/scripts

SPECTRUM_FILES=\
	src/spectrum/application.go\
	src/spectrum/display.go\
	src/spectrum/helpers.go\
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
	src/spectrum/z80_tables.go


PERF_PKG_LIB=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)/⚛perf.a

SDL_PKG_LIB=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)/⚛sdl.a
SDL_AUDIO_PKG_LIB=$(GOROOT)/pkg/$(GOOS)_$(GOARCH)/⚛sdl/audio.a


PKG_LIBS=\
	$(FORMATS_PKG_LIB)\
	$(PERF_PKG_LIB)\
	$(READLINE_PKG_LIB)\
	$(SDL_PKG_LIB)\
	$(SDL_AUDIO_PKG_LIB)


GOFMT_FILES=\
	src/console.go\
	src/gospeccy.go\
	src/spectrum/application.go\
	src/spectrum/keyboard.go\
	src/spectrum/memory.go\
	src/spectrum/port.go\
	src/spectrum/sdlsound.go\
	src/spectrum/sound.go\
	src/spectrum/spectrum.go\
	src/spectrum/ula.go\
	src/spectrum/z80*.go\
	$(FORMATS_FILES)\
	$(FRONTEND_FILES)\
	$(READLINE_FILES)

gospeccy: _obj _obj/gospeccy.$(O)
	$(LD) -L./_obj -o $@ _obj/gospeccy.$(O)

.PHONY: clean
clean: formats-clean readline-clean
	rm -f gospeccy
	rm -f src/spectrum/path_gen_*.go
	rm -rf _obj

.PHONY: gofmt
gofmt:
	gofmt -w -l $(GOFMT_FILES)

.PHONY: install
install: gospeccy
	cp gospeccy $(GOBIN)
	rm -rf $(DIST_PATH)
	mkdir -p $(SYSTEM_ROM_PATH) $(SCRIPTS_PATH)
	cp -a roms/* $(SYSTEM_ROM_PATH)
	cp -a scripts/* $(SCRIPTS_PATH)

.PHONY: uninstall
uninstall: formats-uninstall readline-uninstall
	rm -f $(GOBIN)/gospeccy
	rm -rf $(DIST_PATH)
	rm -rf $(GOROOT)/pkg/$(GOOS)_$(GOARCH)/spectrum

_obj:
	mkdir _obj

_obj/gospeccy.$(O): $(FRONTEND_FILES) _obj/spectrum.a $(PKG_LIBS)
	$(GC) -I./_obj -o $@ $(FRONTEND_FILES)

_obj/spectrum.a: $(SPECTRUM_FILES) $(PKG_LIBS)
	$(GC) -I./_obj -o _obj/spectrum.$(O) $(SPECTRUM_FILES)
	gopack grc $@ _obj/spectrum.$(O)


#
# Installation of external dependencies
#

PERF_URL=github.com/0xe2-0x9a-0x9b/Go-PerfEvents
SDL_URL=github.com/0xe2-0x9a-0x9b/Go-SDL

$(PERF_PKG_LIB):
	goinstall -u $(PERF_URL) || exit 0
	make -C $(GOROOT)/src/pkg/$(PERF_URL) clean
	make -C $(GOROOT)/src/pkg/$(PERF_URL) install

$(SDL_PKG_LIB):
	goinstall -u $(SDL_URL) || exit 0
	make -C $(GOROOT)/src/pkg/$(SDL_URL) clean
	make -C $(GOROOT)/src/pkg/$(SDL_URL) install

$(SDL_AUDIO_PKG_LIBS): $(SDL_PKG_LIB)


include Make.inc-tail
