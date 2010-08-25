include $(GOROOT)/src/Make.inc

TARG=⚛readline
CGOFILES=readline.go
CGO_LDFLAGS=-lreadline

include $(GOROOT)/src/Make.pkg

_cgo_defun.c _cgo_gotypes.go _cgo_export.c _cgo_export.h: $(CGOFILES)
	CGOPKGPATH=$(dir) cgo -- $(CGO_CFLAGS) $(CGOFILES)
	sed s/cgo_readline/cgo_⚛readline/ --in-place _cgo_defun.c
