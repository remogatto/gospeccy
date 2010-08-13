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
	src/perf/perf_amd64.go\


gospeccy: _obj _obj/spectrum.a _obj/gospeccy.$(O)
	$(LD) -L./_obj -o $@ _obj/gospeccy.$(O)

clean:
	rm -f gospeccy
	rm -rf _obj

gofmt:
	gofmt -w -l $(GOFMT_FILES)

_obj:
	mkdir _obj

<<<<<<< HEAD
_obj/gospeccy.8: src/gospeccy.go _obj/spectrum.a
	8g -I./_obj -o $@ src/gospeccy.go

_obj/spectrum.a:
	8g -o _obj/spectrum.8 $^
	gopack grc $@ _obj/spectrum.8

_obj/spectrum.a: src/spectrum/application.go
_obj/spectrum.a: src/spectrum/display.go
_obj/spectrum.a: src/spectrum/keyboard.go
_obj/spectrum.a: src/spectrum/memory.go
_obj/spectrum.a: src/spectrum/port.go
_obj/spectrum.a: src/spectrum/sdldisplay.go
_obj/spectrum.a: src/spectrum/spectrum.go
_obj/spectrum.a: src/spectrum/z80.go
_obj/spectrum.a: src/spectrum/opcodes.go
=======
_obj/gospeccy.$(O): src/gospeccy.go _obj/spectrum.a 
	$(GC) -I./_obj -o $@ src/gospeccy.go

_obj/spectrum.a: $(SPECTRUM_FILES) _obj/perf.a
	$(GC) -I./_obj -o _obj/spectrum.$(O) $(SPECTRUM_FILES)
	gopack grc $@ _obj/spectrum.$(O)

_obj/perf.a: $(PERF_FILES)
	$(GC) -o _obj/perf.$(O) $(PERF_FILES)
	gopack grc $@ _obj/perf.$(O)

>>>>>>> cb3e523d7e46ecd12be20cc5af2b93a66e4f0a46
