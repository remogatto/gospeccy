include $(GOROOT)/src/Make.inc

TARG=⚛readline
CGOFILES=readline.go
CGO_LDFLAGS=-lreadline

include $(GOROOT)/src/Make.pkg

_cgo_defun.c: $(CGOFILES)
	CGOPKGPATH=$(dir) cgo -- $(CGO_CFLAGS) $(CGOFILES)
	sed s/cgo_readline/cgo_⚛readline/ --in-place _cgo_defun.c
