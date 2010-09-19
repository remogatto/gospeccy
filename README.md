# GoSpeccy - An evolving ZX Spectrum 48k Emulator

GoSpeccy is yet another ZX Spectrum (Speccy for friends) Emulator. The
interesting fact is that it is written in Go and - AFAIK - it's the
first Spectrum/Z80 emulator coded with the new language designed by
certain well known Google employees. Last but not least, GoSpeccy is Free Software.

There are a lot of ZX Spectrum emulators around so, why reinventing
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
  machines

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
* An interactive console interface
* SNA format support (48k version)
* SDL backend
* 2x scaler and fullscreen (to be improved)

# Quick start

On Ubuntu Linux you'll need to install the following packages:

    sudo apt-get install libsdl1.2-dev libsdl-mixer1.2-dev libsdl-image1.2-dev libsdl-ttf2.0-dev libreadline6-dev git-core

To install the dependencies and create the gospeccy executable:

    git clone http://github.com/remogatto/gospeccy.git
    cd gospeccy
    gomake
    ./gospeccy

To install (uninstall) gospeccy and its resource files:

    gomake install
    (gomake uninstall)

The following dependencies are installed automatically:

* [⚛Go-SDL](http://github.com/0xe2-0x9a-0x9b/Go-SDL)
* [⚛Go-PerfEvents](http://github.com/0xe2-0x9a-0x9b/Go-PerfEvents)

To make the screen bigger try the "-2x" command line option,
or type "scale(2)" in the interactive console.

To try the classic Hello World try to press the following keys:

    p
    CTRL+p
    hello world
    CTRL+p
    RETURN

And see your shining new ZX Spectrum 48k responding :)

For a nice picture of the speccy keyboard layout visit this [page](http://www.guybrush.demon.co.uk/spectrum/docs/Basic.htm).

To load a SNA file:

    ./gospeccy IMAGE.sna

And if you're curious to see what this machine can do, try the simple
[Fire104b](http://pouet.net/prod.php?which=54076) intro by Andrew
Gerrand included in the gospeccy distribution! In the gospeccy folder,
run:

    ./gospeccy -2x snapshots/Syntax09nF.sna

For more, try searching the Internet find some ZX Spectrum 48k
games and demos in SNA format (well, good luck with that). For example:

* [World of spectrum](http://www.worldofspectrum.org/archive.html)
* [Pouet.net search](http://pouet.net/prodlist.php?platform[]=ZX+Spectrum)

In order to create a SNA file from another 48k format, install the
[FUSE](http://fuse-emulator.sourceforge.net/), load the original file
into FUSE (e.g: fuse -m48 file.tap), wait until it loads, and save it
in SNA format as "file.sna". For example, you can test this procedure on
the [48K](http://pouet.net/prod.php?which=54504) ZX demo. All these things
will hopefully improve as GoSpeccy matures, so that you won't need to resort
to using other ZX Spectrum emulators.

# Key bindings

    Host computer   Zx Spectrum
    ---------------------------
    CTRL            SYMBOL SHIFT
    LEFT SHIFT      CAPS SHIFT
    [a-z0-9]        [A-Z0-9]
    SPACE           SPACE

For more info about keybindings see <tt>spectrum/keyboard.go</tt>

# Proprietary ROMs

Generally, proprietary roms are protected by copyright so none of them
is included in GoSpeccy (with the exception of the 48k rom that can be
freely distributed). BTW, you can find tons of roms for the ZX
Spectrum on the Internet. Take a look at:

* JSSpeccy [svn](http://svn.matt.west.co.tt/svn/jsspec/trunk/snapshots/) repository

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
`$HOME/.gospeccy/sna/somegame.sna` simply execute:

<pre>
gospeccy somegame.sna
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
* Add support for tape saving & loading
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
* http://freecella.blogspot.com

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

