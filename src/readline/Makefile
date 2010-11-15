include $(GOROOT)/src/Make.inc

TARG=spectrum/readline
CGOFILES=readline.go
CGO_LDFLAGS=-lreadline

include $(GOROOT)/src/Make.pkg

.PHONY: uninstall
uninstall:
	rm -f $(INSTALLFILES)
