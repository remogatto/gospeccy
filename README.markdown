# GOSpeccy - A naive ZX Spectrum 48k Emulator

GOSpeccy is yet another ZX Spectrum (Speccy for friends) Emulator. The
interesting fact is that it is written in GO and - AFAIK - it's the
first Spectrum/Z80 emulator coded with the new language by Google.

There are a lot of ZX Spectrum emulators around so, why reinventing
the wheel? Well, because it was a great learning experience about
emulators for me :) And it was a good chance to write something real
with GO.

Coding an emulator in GO is very enjoyable. The language is simple,
pragmatic and fast and has a lot of features that help in emulators
development (for example its low-level similarities with C but GO is
*a lot* more agile than C!). More on that in a future blog post maybe
:)

The Zilog Z80 emulation is the core of GOSpeccy. The CPU emulation
code is generated using the <tt>z80.pl</tt> script shipped with
[FUSE](http://fuse-emulator.sourceforge.net/) (one of the best ZX
Spectrum emulator around). The script has been hacked to generate GO
code rather than C code.

Another source of inspiration was
[JSSpeccy](http://matt.west.co.tt/spectrum/jsspeccy/), a neat
Javascript Speccy emulator.

The Z80 emulation is tested using the excellent test-suite shipped
with FUSE.

If you like this software, please watch it on
[github](http://github.com/remogatto/gospeccy)! Seeing a growing number of watchers is
an excellent motivation for me to keep up this work :) Bug reports and
testing are also appreciated! And don't forget to send me patches, of
course ;)

# Features

* Complete Zilog Z80 emulation
* SNA format support
* SDL backend

# Dependencies

* [banthar/Go-SDL](http://github.com/banthar/Go-SDL)

# Quick start

First of all, be sure to install the dependencies. Then:

    git clone http://github.com/remogatto/gospeccy
    cd gospeccy
    make install
    gospeccy -d # Execute in a double sized window

Now try press the following keys:

    p
    CTRL+p
    hello world
    CTRL+p
    RETURN

And see your shining new ZX Spectrum computer responding :)

To load a SNA rom:

    gospeccy -d image.sna

# Key bindings

    Host computer   Zx Spectrum
    ---------------------------
    CTRL            Symbol Shift
    LEFT SHIFT      Caps

For more info about keybindings see <tt>spectrum/keyboard.go</tt>

# ROMs

Generally, roms are protected by copyright so none of them is included
in GOSpeccy (with the exception of the 48k rom that can be freely
distributed). BTW, you can find tons of roms for the ZX Spectrum on
the Internet. Take a look at JSSpeccy
[svn](http://svn.matt.west.co.tt/svn/jsspec/trunk/snapshots/)
repository.

# TODO

* Add sound emulation
* Improve memory contention
* Add support for more file formats (take a look [here](http://www.worldofspectrum.org/faq/reference/formats.htm))
* Better performances
* Add new backends (exp/draw?)

# Contacts

* andrea.fazzi@alcacoop.it
* http://twitter.com/remogatto
* http://freecella.blogspot.com

# LICENSE

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

