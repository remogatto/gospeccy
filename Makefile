include $(GOROOT)/src/Make.$(GOARCH)

SPECTRUM_FILES=\
	src/spectrum/application.go\
	src/spectrum/console.go\
	src/spectrum/display.go\
	src/spectrum/keyboard.go\
	src/spectrum/memory.go\
	src/spectrum/opcodes.go\
	src/spectrum/port.go\
	src/spectrum/sdldisplay.go\
	src/spectrum/spectrum.go\
	src/spectrum/z80.go\

PERF_FILES=\
	src/perf/perf.go\
	src/perf/perf_$(GOARCH).go\
	src/perf/types.$(O).go\

GOFMT_FILES=\
	src/gospeccy.go\
	src/spectrum/application.go\
	src/spectrum/console.go\
	src/spectrum/keyboard.go\
	src/spectrum/memory.go\
	src/spectrum/port.go\
	src/spectrum/spectrum.go\
	src/spectrum/z80.go\
	src/spectrum/z80_test.go\
	src/perf/perf.go\
	src/perf/perf_386.go\


gospeccy: _obj _obj/spectrum.a _obj/gospeccy.$(O)
	$(LD) -L./_obj -o $@ _obj/gospeccy.$(O)

clean:
	rm -f gospeccy
	rm -rf _obj

gofmt:
	gofmt -w -l $(GOFMT_FILES)

_obj:
	mkdir _obj

_obj/gospeccy.$(O): src/gospeccy.go _obj/spectrum.a 
	$(GC) -I./_obj -o $@ src/gospeccy.go

_obj/spectrum.a: $(SPECTRUM_FILES) _obj/perf.a
	$(GC) -I./_obj -o _obj/spectrum.$(O) $(SPECTRUM_FILES)
	gopack grc $@ _obj/spectrum.$(O)

_obj/perf.a: $(PERF_FILES)
	$(GC) -o _obj/perf.$(O) $(PERF_FILES)
	gopack grc $@ _obj/perf.$(O)

