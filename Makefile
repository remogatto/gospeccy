include $(GOROOT)/src/Make.$(GOARCH)

#Run 'make deps' if you add or change new folders.
DIRS=\
	spectrum\
	bin
TEST=\
	spectrum\

all.dirs: $(addsuffix .all, $(DIRS))
clean.dirs: $(addsuffix .clean, $(DIRS))
install.dirs: $(addsuffix .install, $(DIRS))
test.dirs: $(addsuffix .test, $(TEST))
docs.dirs: $(addsuffix .docs, $(DIRS))
format.dirs: $(addsuffix .format, $(DIRS))

%.all:
	+cd $* && make

%.clean:
	+cd $* && make clean

%.install:
	+cd $* && make install

%.test:
	+cd $* && make test

%.docs:
	+godoc -pkgroot="." -html -v $* > docs/$*.html

%.format:
	+cd $* && gofmt -w *.go

clean:	clean.dirs

install: install.dirs
	cd bin && make install_res

test:	test.dirs

docs:	docs.dirs

format: format.dirs

deps:
	./deps.bash

include Make.deps

