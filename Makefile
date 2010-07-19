gospeccy: _obj _obj/spectrum.a _obj/gospeccy.8
	8l -L./_obj -o $@ _obj/gospeccy.8

clean:
	rm -f gospeccy
	rm -rf _obj

.PHONY:
_obj:
	mkdir _obj

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
