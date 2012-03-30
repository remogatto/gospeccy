/*

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

*/

// GoSpeccy
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/remogatto/gospeccy/src/env"
	"github.com/remogatto/gospeccy/src/formats"
	"github.com/remogatto/gospeccy/src/interpreter"
	"github.com/remogatto/gospeccy/src/spectrum"
	"net/url"
	"os"
	pathutil "path"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	// Pull-in all optional modules into the final executable
	_ "github.com/remogatto/gospeccy/src/pull_modules"
)

type handler_SIGTERM struct {
	app *spectrum.Application
}

func (h *handler_SIGTERM) HandleSignal(s os.Signal) {
	switch ss := s.(type) {
	case syscall.Signal:
		switch ss {
		case syscall.SIGTERM, syscall.SIGINT:
			if h.app.Verbose {
				h.app.PrintfMsg("%v", ss)
			}

			h.app.RequestExit()
		}
	}
}

func newApplication(verbose bool) *spectrum.Application {
	app := spectrum.NewApplication()
	app.Verbose = verbose
	env.Publish(app)
	return app
}

func newEmulationCore(app *spectrum.Application, acceleratedLoad bool) (*spectrum.Spectrum48k, error) {
	romPath, err := spectrum.SystemRomPath("48.rom")
	if err != nil {
		return nil, err
	}

	rom, err := spectrum.ReadROM(romPath)
	if err != nil {
		return nil, err
	}

	speccy := spectrum.NewSpectrum48k(app, *rom)
	if acceleratedLoad {
		speccy.TapeDrive().AcceleratedLoad = true
	}

	env.Publish(speccy)

	return speccy, nil
}

func ftpget_choice(app *spectrum.Application, matches []string, freeware []bool) (string, error) {
	switch len(matches) {
	case 0:
		return "", nil

	case 1:
		if freeware[0] {
			return matches[0], nil
		} else {
			// Not freeware - We want the user to make the choice
		}
	}

	app.PrintfMsg("")
	fmt.Printf("Select a number from the above list (press ENTER to exit GoSpeccy): ")
	in := bufio.NewReader(os.Stdin)

	input, err := in.ReadString('\n')
	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return "", nil
	}

	id, err := strconv.Atoi(input)
	if err != nil {
		return "", err
	}
	if (id < 0) || (id >= len(matches)) {
		return "", errors.New("Invalid selection")
	}

	url := matches[id]
	if app.Verbose {
		app.PrintfMsg("You've selected %s", url)
	}
	return url, nil
}

func wait(app *spectrum.Application) {
	<-app.HasTerminated

	if app.Verbose {
		var memstats runtime.MemStats
		runtime.ReadMemStats(&memstats)
		app.PrintfMsg("GC: %d garbage collections, %s total pause time",
			memstats.NumGC, time.Nanosecond*time.Duration(memstats.PauseTotalNs))
	}

	// Stop host-CPU profiling
	if *cpuProfile != "" {
		pprof.StopCPUProfile() // flushes profile to disk
	}
}

func exit(app *spectrum.Application) {
	app.RequestExit()
	wait(app)
}

var (
	help            = flag.Bool("help", false, "Show usage")
	acceleratedLoad = flag.Bool("accelerated-load", false, "Accelerated tape loading")
	fps             = flag.Float64("fps", spectrum.DefaultFPS, "Frames per second")
	verbose         = flag.Bool("verbose", false, "Enable debugging messages")
	cpuProfile      = flag.String("hostcpu-profile", "", "Write host-CPU profile to the specified file (for 'pprof')")
	wos             = flag.String("wos", "", "Download from WorldOfSpectrum; you must provide a query regex (ex: -wos=jetsetwilly)")
)

func main() {
	var init_waitGroup sync.WaitGroup
	env.PublishName("init WaitGroup", &init_waitGroup)

	// Handle options
	{
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "GoSpeccy - A ZX Spectrum 48k Emulator written in Go\n\n")
			fmt.Fprintf(os.Stderr, "Usage:\n\n")
			fmt.Fprintf(os.Stderr, "\tgospeccy [options] [image.sna]\n\n")
			fmt.Fprintf(os.Stderr, "Options are:\n\n")
			flag.PrintDefaults()
		}

		flag.Parse()

		if *help == true {
			flag.Usage()
			return
		}
	}

	// Start host-CPU profiling (if enabled).
	// The setup code is based on the contents of Go's file "src/pkg/testing/testing.go".
	var pprof_file *os.File
	if *cpuProfile != "" {
		var err error

		pprof_file, err = os.Create(*cpuProfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			return
		}

		err = pprof.StartCPUProfile(pprof_file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to start host-CPU profiling: %s", err)
			pprof_file.Close()
			return
		}
	}

	app := newApplication(*verbose)

	// Use at least 2 OS threads.
	// This helps to prevent audio buffer underflows
	// in case rendering is consuming too much CPU.
	if (os.Getenv("GOMAXPROCS") == "") && (runtime.GOMAXPROCS(-1) < 2) {
		runtime.GOMAXPROCS(2)
	}
	if app.Verbose {
		app.PrintfMsg("using %d OS threads", runtime.GOMAXPROCS(-1))
	}

	// Install SIGTERM handler
	{
		handler := handler_SIGTERM{app}
		spectrum.InstallSignalHandler(&handler)
	}

	speccy, err := newEmulationCore(app, *acceleratedLoad)
	if err != nil {
		app.PrintfMsg("%s", err)
		exit(app)
		return
	}

	// Run startup scripts.
	// The startup scripts may change the display settings or enable/disable the audio.
	// They may also terminate the program.
	{
		interpreter.Init(app, flag.Arg(0), speccy)

		if app.TerminationInProgress() || app.Terminated() {
			exit(app)
			return
		}
	}

	// Optional: Read and categorize the contents
	//           of the file specified on the command-line
	var program_orNil interface{} = nil
	var programName string
	if flag.Arg(0) != "" {
		file := flag.Arg(0)
		programName = file

		var err error
		path, err := spectrum.ProgramPath(file)
		if err != nil {
			app.PrintfMsg("%s", err)
			exit(app)
			return
		}

		program_orNil, err = formats.ReadProgram(path)
		if err != nil {
			app.PrintfMsg("%s", err)
			exit(app)
			return
		}
	} else if *wos != "" {
		var records []spectrum.WosRecord
		records, err := spectrum.WosQuery(app, "regexp="+url.QueryEscape(strings.Replace(*wos, " ", "*", -1)))
		if err != nil {
			app.PrintfMsg("%s", err)
			exit(app)
			return
		}

		var urls []string
		var isFreeware []bool
		for _, record := range records {
			var freeware bool = (strings.ToLower(record.Publication) == "freeware")

			for _, url := range record.FtpFiles {
				urls = append(urls, url)
				isFreeware = append(isFreeware, freeware)
				if freeware {
					app.PrintfMsg("[%d] - [Freeware] %s", len(urls)-1, url)
				} else {
					app.PrintfMsg("[%d] - [Not freeware] %s", len(urls)-1, url)
				}
			}
		}
		if len(urls) != 1 {
			app.PrintfMsg("%d matches", len(urls))
		} else {
			app.PrintfMsg("1 match")
		}

		url, err := ftpget_choice(app, urls, isFreeware)
		if err != nil {
			app.PrintfMsg("%s", err)
			exit(app)
			return
		}

		if url != "" {
			filePath, err := spectrum.WosGet(app, os.Stdout, url)
			if err != nil {
				app.PrintfMsg("get %s: %s", url, err)
				exit(app)
				return
			}

			program_orNil, err = formats.ReadProgram(filePath)
			if err != nil {
				app.PrintfMsg("%s", err)
				exit(app)
				return
			}

			_, programName = pathutil.Split(filePath)
		} else {
			exit(app)
			return
		}
	}

	// Wait until modules are initialized
	init_waitGroup.Wait()

	// Begin speccy emulation
	go speccy.EmulatorLoop()

	// Set the FPS
	speccy.CommandChannel <- spectrum.Cmd_SetFPS{float32(*fps), nil}

	// Optional: Load the program specified on the command-line
	if program_orNil != nil {
		program := program_orNil

		if _, isTAP := program.(*formats.TAP); isTAP {
			romLoaded := make(chan (<-chan bool))
			speccy.CommandChannel <- spectrum.Cmd_Reset{romLoaded}
			<-(<-romLoaded)
		}

		errChan := make(chan error)
		speccy.CommandChannel <- spectrum.Cmd_Load{programName, program, errChan}
		err := <-errChan
		if err != nil {
			app.PrintfMsg("%s", err)
			exit(app)
			return
		}
	}

	wait(app)
}
