include $(GOROOT)/src/Make.$(GOARCH)

gospeccy: _obj _obj/spectrum.a _obj/gospeccy.$(O)
	$(LD) -L./_obj -o $@ _obj/gospeccy.$(O)

clean:
	rm -f gospeccy
	rm -rf _obj

.PHONY:
_obj:
	mkdir _obj

_obj/gospeccy.$(O): src/gospeccy.go _obj/spectrum.a
	$(GC) -I./_obj -o $@ src/gospeccy.go

_obj/spectrum.a:
	$(GC) -o _obj/spectrum.$(O) $^
	gopack grc $@ _obj/spectrum.$(O)

_obj/spectrum.a: src/spectrum/application.go
_obj/spectrum.a: src/spectrum/display.go
_obj/spectrum.a: src/spectrum/keyboard.go
_obj/spectrum.a: src/spectrum/memory.go
_obj/spectrum.a: src/spectrum/port.go
_obj/spectrum.a: src/spectrum/sdldisplay.go
_obj/spectrum.a: src/spectrum/spectrum.go
_obj/spectrum.a: src/spectrum/z80.go
