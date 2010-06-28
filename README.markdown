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
*a lot* more agile than C!). More about that in a future blog post may
be :)

The Zilog Z80 emulation is the core of GOSpeccy. The CPU emulation
code is generated using the <tt>z80.pl</tt> script shipped with FUSE
(one of the best ZX Spectrum emulator around). The script has been
hacked to generate GO code rather than C code.

Another source of inspiration was JSSpeccy, a neat Javascript Speccy
emulator.

The Z80 emulation is tested using the excellent test-suite shipped
with FUSE (see spectrum/z80_test, spectrum/tests.in and
spectrum/tests.expected files).

If you like this software, please watch it on github! Seeing a growing
number of watchers is a good motivation for me to keep up this work :)
And don't forget to send me patches, of course ;)

# Features

* Complete Zilog Z80 emulation
* SNA format support
* SDL backend

# Dependencies

* Go-SDL

# Quick start

First of all, be sure to install the dependencies. Then:

  git clone http://github.com/remogatto/gospeccy
  cd gospeccy
  make install
  gospeccy -d # Execute in a double sized window

To load a SNA rom:

  gospeccy -d image.sna

# ROMs

Generally, roms are protected by copyright so none of them is included
in GOSpeccy (with the exception of the 48k rom that can be freely
distributed). BTW, you can find tons of roms for the ZX Spectrum on
the Internet. Take a look at JSSpeccy svn repository.

# TODO

* Add sound emulation
* Improve memory contention
* Add support for more file formats (tap, szx, etc)
* Better performances
* Add new backends (exp/draw?)

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

