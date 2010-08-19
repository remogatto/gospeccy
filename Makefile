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
	src/spectrum/ula.go\
	src/spectrum/z80.go\
	src/spectrum/z80_gen.go\
	src/spectrum/z80_tables.go\

GOFMT_FILES=\
	src/gospeccy.go\
	src/spectrum/application.go\
	src/spectrum/console.go\
	src/spectrum/keyboard.go\
	src/spectrum/memory.go\
	src/spectrum/port.go\
	src/spectrum/spectrum.go\
	src/spectrum/ula.go\
	src/spectrum/z80.go\
	src/spectrum/z80_test.go\
	src/spectrum/z80_gen.go\
	src/spectrum/z80_tables.go\


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

_obj/spectrum.a: $(SPECTRUM_FILES)
	$(GC) -I./_obj -o _obj/spectrum.$(O) $(SPECTRUM_FILES)
	gopack grc $@ _obj/spectrum.$(O)
