GOSPECCY_ROOT=../..
include $(GOSPECCY_ROOT)/Make.inc-head

TARG=spectrum/formats
GOFILES=$(FORMATS_FILES)

include $(GOROOT)/src/Make.pkg
include $(GOSPECCY_ROOT)/Make.inc-tail

.PHONY: uninstall
uninstall:
	rm -f $(INSTALLFILES)
