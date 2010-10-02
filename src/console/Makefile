GOSPECCY_ROOT=../..
include $(GOSPECCY_ROOT)/Make.inc-head

TARG=spectrum/console
GOFILES=$(CONSOLE_FILES)

DEPS+="../spectrum"

include $(GOROOT)/src/Make.pkg
include $(GOSPECCY_ROOT)/Make.inc-tail

.PHONY: uninstall
uninstall:
	rm -f $(INSTALLFILES)
