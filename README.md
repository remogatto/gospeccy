[![Build Status](https://drone.io/github.com/remogatto/gospeccy/status.png)](https://drone.io/github.com/remogatto/gospeccy/latest)

# GoSpeccy - An evolving ZX Spectrum 48k Emulator

GoSpeccy is a free ZX Spectrum (Speccy for friends) emulator
written in [Go](http://golang.org).

# Quick start

Installing and starting GoSpeccy with Go1 is simple:

    go get -v github.com/remogatto/gospeccy/src/gospeccy
    gospeccy
    gospeccy -wos="interlace demo"

# Description

GoSpeccy is based on a [concurrent](http://github.com/remogatto/gospeccy/wiki/Architecture) architecture.
We think the concurrency is a strong peculiarity of GoSpeccy as it opens new
interesting scenarios when developing and using the emulator.

Go has interesting features that help in emulators development:
 
* it has "low-level" similarities with C while allowing more productivity
  than C in certain cases

* it is strongly typed and type safe so you are aware about certain errors at
  compile-time

* it is garbage collected so there is small chance of memory leaks

* it has an uint16 built-in type that helps dealing with 8/16 bit
  emulation

* it has goroutines to enable concurrency in the emulator design
  and implementation

The Zilog Z80 CPU emulation is the core of GoSpeccy. The CPU emulation
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

* Complete (almost) Zilog Z80 emulation
* Concurrent [architecture](http://github.com/remogatto/gospeccy/wiki/Architecture)
* Beeper support
* Initial support for Kempston joysticks
* An interactive on-screen console interface based on [clingon](http://github.com/remogatto/clingon)
* Snapshot support: SNA, Z80 formats (48k versions)
* Tape support (TAP format, read-only)
* Accelerated tape loading
* ZIP files support
* SDL backend
* 2x scaler and fullscreen

# Using GoSpeccy

To make the screen bigger try the "-2x" command line option,
or type "scale(2)" in the interactive console.

To try the classic Hello World try to press the following keys:

    p
    CTRL+p
    hello world
    CTRL+p
    RETURN

For a nice picture of the speccy keyboard layout visit this
[page](http://www.guybrush.demon.co.uk/spectrum/docs/Basic.htm).

To load a program run:

    gospeccy file.tap

To enable tape loading acceleration use the <tt>accelerated-load</tt>
option. For a complete list of the command-line options run:

    gospeccy -help

If you can't wait to see what this machine can do, try the nice
[Fire104b](http://pouet.net/prod.php?which=54076) intro by Andrew
Gerrand included in the gospeccy distribution! In the gospeccy folder,
run:

    gospeccy -2x snapshots/Syntax09nF.z80

To automatically download a program from [World of spectrum](http://www.worldofspectrum.org),
and start it:

    gospeccy -wos="horace*tower"

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

For more info about key bindings see file <tt>spectrum/keyboard.go</tt>

# Proprietary games and system ROM

Generally, games/programs are protected by copyright so none of them
is included in GoSpeccy. However, it is possible to find tons of games for
the ZX Spectrum on the Internet. The system ROM for Spectrum 48k can be freely
distributed and so it's included in the GoSpeccy distribution.

# Convention over Configuration

Loading files in the emulator relies on a Convention over
Configuration approach. To enjoy it, you should create the following
folder structure:

<pre>
mkdir -p $HOME/.config/gospeccy/roms			# System roms folder
mkdir -p $HOME/.config/gospeccy/programs		# Scripts folder
mkdir -p $HOME/.config/gospeccy/scripts			# Scripts folder
</pre>

If you like to add your custom search path, In the scripts folder,
create file `config_local.go` with the following contents:

<pre>
// Search path for programs, scripts, etc
addSearchPath("/home/user/gospeccy")
</pre>

After this, to load `/home/user/gospeccy/programs/somegame.z80` simply execute:

<pre>
gospeccy somegame.z80
</pre>

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

# osx notes
The SDL usage needs X11 because the Quartz implementation is picky about how threads are used (you need to poll events on the main thread or it throws an exception.)
Getting a sdl installed in osx which supports x11 is not that simple though because the bottled versions seem to exclude X11. The following worked for me.

    brew install -vd --build-from-source  sdl --with-x11 --with-test
    brew install sdl_image sdl_ttf sdl_mixer


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

