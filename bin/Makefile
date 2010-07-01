include $(GOROOT)/src/Make.$(GOARCH)

TARG=gospeccy
GOFILES=\
	gospeccy.go
GOSPECCYDIR=$(HOME)/.gospeccy/

install_res:
	mkdir --parent $(GOSPECCYDIR) && cp -r res/* $(GOSPECCYDIR)/

include $(GOROOT)/src/Make.cmd

