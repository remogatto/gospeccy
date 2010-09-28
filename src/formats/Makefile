GOSPECCY_ROOT=../..
include $(GOSPECCY_ROOT)/Make.inc-head

TARG=spectrum/formats
GOFILES=$(FORMATS_FILES)

Z80_test.go: $(PRETTYTEST_PKG_LIB)
_obj/$(TARG).a: $(ZIP_PKG_LIB)

include $(GOROOT)/src/Make.pkg
include $(GOSPECCY_ROOT)/Make.inc-tail

.PHONY: uninstall
uninstall:
	rm -f $(INSTALLFILES)
