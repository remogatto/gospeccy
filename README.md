# GoSpeccy - An evolving ZX Spectrum 48k Emulator

GoSpeccy is yet another ZX Spectrum (Speccy for friends) Emulator. The
interesting fact is that it is written in [Go](http://golang.org) and - AFAIK - it's the
first and only 8-bit machine emulator coded with the new language designed by
certain well known Google employees. Last but not least, GoSpeccy is Free Software.

There are a lot of ZX Spectrum emulators around - so why reinvent
the wheel? Well, mainly for [amarcord](http://en.wikipedia.org/wiki/Amarcord)
reasons and then because it was a great learning experience about emulators
and Go. In addition, thanks to the fundamental contribution of
[⚛](http://github.com/0xe2-0x9a-0x9b), GoSpeccy is now based on a
[concurrent](http://github.com/remogatto/gospeccy/wiki/Architecture)
architecture. We think the concurrency is a strong peculiarity of GoSpeccy
as it opens new interesting scenarios when developing and using the emulator.

Among other things, coding an emulator in Go is very enjoyable. The
language is simple, pragmatic and fast (well... fast enough to run a
Z80 emulator!). It has a lot of features that help in emulators
development:
 
* it has "low-level" similarities with C while allowing more productivity
  than C in certain cases

* it is strongly typed and type safe so you are aware about certain errors at
  compile-time

* it is garbage collected so you haven't to worry about memory leaks

* it has an int16 built-in type that helps dealing with 8/16 bit
  emulation

* as already mentioned, it has goroutines to enable concurrency in the
  emulator design and implementation

The Zilog Z80 emulation is the core of GoSpeccy. The CPU emulation
code is generated using a modified version of the <tt>z80.pl</tt> script
shipped with [FUSE](http://fuse-emulator.sourceforge.net/) (one of the best ZX
Spectrum emulators around). The script has been hacked to generate Go
code rather than C code.

Another source of inspiration was
[JSSpeccy](http://matt.west.co.tt/spectrum/jsspeccy/), a neat
Javascript Speccy emulator.

The Z80 emulation is tested against the excellent test-suite shipped
with FUSE.

If you like this software, please watch it on
[github](http://github.com/remogatto/gospeccy)! Seeing a growing
number of watchers is an excellent motivation for the GoSpeccy team to
keep up this work :) Bug reports and testing are also appreciated! And
don't forget to fork and send patches, of course ;)

# Features

* Complete Zilog Z80 emulation
* Concurrent [architecture](http://github.com/remogatto/gospeccy/wiki/Architecture)
* Beeper support
* Initial support for Kempston joysticks
* An interactive on-screen console interface based on [clingon](http://github.com/remogatto/clingon)
* Snapshot support: SNA, Z80 formats (48k versions)
* Tape support (TAP format)
* Accelerated tape loading
* ZIP files support
* SDL backend
* 2x scaler and fullscreen

# Quick start

On Ubuntu Linux you'll need to install the following packages:

    sudo apt-get install libsdl1.2-dev libsdl-mixer1.2-dev libsdl-image1.2-dev libsdl-ttf2.0-dev git-core

GoSpeccy is using the GOAM build tool. To install it:

    goinstall github.com/0xe2-0x9a-0x9b/goam
    cd $GOROOT/src/pkg/github.com/0xe2-0x9a-0x9b/goam
    make install

To install the dependencies and create the gospeccy executable:

    git clone http://github.com/remogatto/gospeccy.git
    cd gospeccy
    goam make
    ./gospeccy

To install (uninstall) gospeccy and its resource files:

    goam install
    (goam uninstall)

The following dependencies are installed automatically:

* [⚛Go-SDL](http://github.com/0xe2-0x9a-0x9b/Go-SDL)
* [⚛Go-PerfEvents](http://github.com/0xe2-0x9a-0x9b/Go-PerfEvents)
* [clingon](http://github.com/remogatto/clingon")
* [prettytest](http://github.com/remogatto/prettytest")

To make the screen bigger try the "-2x" command line option,
or type "scale(2)" in the interactive console.

To try the classic Hello World try to press the following keys:

    p
    CTRL+p
    hello world
    CTRL+p
    RETURN

And see your shining new ZX Spectrum 48k responding :)

For a nice picture of the speccy keyboard layout visit this
[page](http://www.guybrush.demon.co.uk/spectrum/docs/Basic.htm).

To load a program simply run:

    ./gospeccy file.tap

To enable tape loading acceleration use the <tt>accelerated-load</tt>
option. For a complete list of the command-line options run:

    ./gospeccy -help

If you can't wait to see what this machine can do, try the nice
[Fire104b](http://pouet.net/prod.php?which=54076) intro by Andrew
Gerrand included in the gospeccy distribution! In the gospeccy folder,
run:

    ./gospeccy -2x snapshots/Syntax09nF.z80

For more, try searching the Internet for ZX Spectrum 48k
games and demos in Z80 and TAP format. For example:

* [World of spectrum](http://www.worldofspectrum.org/archive.html)
* [Pouet.net search](http://pouet.net/prodlist.php?platform[]=ZX+Spectrum)
* [Forever](http://forever.zeroteam.sk/download.htm)

# Key bindings

    Host computer   ZX Spectrum
    ---------------------------
    CTRL            SYMBOL SHIFT
    LEFT SHIFT      CAPS SHIFT
    [a-z0-9]        [A-Z0-9]
    SPACE           SPACE

For more info about key bindings see <tt>spectrum/keyboard.go</tt>

# Proprietary games and system ROM

Generally, games/programs are protected by copyright so none of them
is included in GoSpeccy. BTW, you can find tons of games for the ZX
Spectrum on the Internet. The system ROM for Spectrum 48k can be freely
distributed and so it's included in the GoSpeccy distribution.

# Convention over Configuration

Loading files in the emulator relies on a Convention over
Configuration approach. To enjoy it, you should create the following
folder structure:

<pre>
mkdir -p $HOME/.gospeccy/sna		# Snapshots folder
mkdir -p $HOME/.gospeccy/roms		# System roms folder
mkdir -p $HOME/.gospeccy/scripts	# Scripts folder
</pre>

Then put your snapshots, system roms or script files in the
corresponding folder. After this, to load
`$HOME/.gospeccy/sna/somegame.z80` simply execute:

<pre>
gospeccy somegame.z80
</pre>

The same applies to `load()` and `script()` functions in the
interactive console.

The default Spectrum 48k system ROM is copied in
`$GOROOT/pkg/$GOOS_$GOARCH/gospeccy/roms` during the installation
process. This is the default system rom loaded by the emulator. You
can override this behaviour copying your 48k rom in
`$HOME/.gospeccy/roms`.

# Screenshots

Manic Miner running on GoSpeccy.

![Manic Miner running on GoSpeccy](http://sites.google.com/site/remogatto/gospeccy_running_scaled.png)

# To Do

* Fix some memory and I/O contention bugs
* Add support for more file formats (take a look [here](http://www.worldofspectrum.org/faq/reference/formats.htm))
* Add support for tape saving
* Better general performance
* Add more filters and improve the scaler
* Add new graphical backends (Go's exp/draw?)

# Credits

* Thanks to [⚛](http://github.com/0xe2-0x9a-0x9b) for giving a new
  whole direction to this project.
* Thanks to the people on
  [golang-nuts](http://groups.google.com/group/golang-nuts) for giving
  feedback and support.
* Thanks to Andrew Gerrand for the crackling Fire104b demo.

# Contacts

* andrea.fazzi@alcacoop.it
* http://twitter.com/remogatto
* http://www.facebook.com/remochat
* http://remogatto.github.com/

# License

Copyright (c) 2010 Andrea Fazzi

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

