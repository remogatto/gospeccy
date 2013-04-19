
package spectrum

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

// ===========
// Application
// ===========

type Application struct {
	exitApp       chan byte // If [this channel is closed] then [the whole application should terminate]
	HasTerminated chan byte // This channel is closed after the whole application has terminated

	eventLoops []*EventLoop

	terminationInProgress bool
	terminated            bool

	mutex sync.Mutex

	messageOutput MessageOutput

	Verbose         bool
	VerboseShutdown bool

	CreationTime time.Time // The time when this Application object was created
}

func NewApplication() *Application {
	app := &Application{
		exitApp:       make(chan byte),
		HasTerminated: make(chan byte),
		eventLoops:    make([]*EventLoop, 0, 8),
		CreationTime:  time.Now(),
		messageOutput: &stdoutMessageOutput{},
	}

	go appGoroutine(app)

	return app
}

func appGoroutine(app *Application) {
	// Block until there is a request to exit the application
	<-app.exitApp

	var startTime time.Time
	if app.Verbose {
		startTime = time.Now()
	}

	// Cycle until there are no EventLoop objects.
	// Usually, the body of this 'for' statement executes only once
	for {
		// Get the 'app.eventLoops' array, then clear 'app.eventLoops'
		app.mutex.Lock()
		eventLoops := app.eventLoops
		app.eventLoops = make([]*EventLoop, 0, 8)
		app.mutex.Unlock()

		// This is a procedure of two phases:
		//
		//	1. Pause all event loops (i.e: stop all tickers)
		//	2. Terminate all event loops
		//
		// The two phases are required because we have no idea what the relationship
		// among the event-loops might be. We are assuming the worst case scenario
		// in which the relationships among event-loops form a graph - in such a case
		// it is unclear whether an event-loop can be terminated without knowing
		// that all event-loops are paused.

		// Pause all event loops
		for _, e := range eventLoops {
			if app.VerboseShutdown {
				app.PrintfMsg("sending pause message to %s", e.Name)
			}
			// Request the event-loop to pause, and wait until it actually pauses
			e.Pause <- 0
			<-e.Pause
			if app.VerboseShutdown {
				app.PrintfMsg("%s is now paused", e.Name)
			}
		}

		// Terminate all event loops
		for _, e := range eventLoops {
			if app.VerboseShutdown {
				app.PrintfMsg("sending terminate message to %s", e.Name)
			}
			// Request the event-loop to terminate, and wait until it actually terminates
			e.Terminate <- 0
			<-e.Terminate
			if app.VerboseShutdown {
				app.PrintfMsg("%s has terminated", e.Name)
			}
		}

		app.mutex.Lock()
		if len(app.eventLoops) == 0 {
			app.mutex.Unlock()
			break
		} else {
			// Some new EventLoops were created while we were pausing&terminating the known ones
			app.mutex.Unlock()
		}
	}

	if app.Verbose {
		endTime := time.Now()
		app.PrintfMsg("application shutdown completed after %s", endTime.Sub(startTime))
		app.PrintfMsg("application has terminated")
	}

	close(app.HasTerminated)

	app.mutex.Lock()
	app.terminationInProgress = false
	app.terminated = true
	app.mutex.Unlock()
}

// Returns whether the operation succeeded.
// False is returned if the application has already terminated.
func (app *Application) addEventLoop(e *EventLoop) bool {
	app.mutex.Lock()

	if app.terminated {
		// Fail
		app.mutex.Unlock()
		return false
	}

	// At this point: The application is running,
	//                or it is in the process of being terminated.

	app.eventLoops = append(app.eventLoops, e)
	app.mutex.Unlock()
	return true
}

func (app *Application) RequestExit() {
	app.mutex.Lock()
	{
		if app.terminationInProgress || app.terminated {
			app.mutex.Unlock()
			return
		}
		app.terminationInProgress = true
	}
	app.mutex.Unlock()

	close(app.exitApp)
}

func (app *Application) TerminationInProgress() bool {
	app.mutex.Lock()
	a := app.terminationInProgress
	app.mutex.Unlock()
	return a
}

func (app *Application) Terminated() bool {
	app.mutex.Lock()
	a := app.terminated
	app.mutex.Unlock()
	return a
}

// Get the current MessageOutput
func (app *Application) GetMessageOutput() MessageOutput {
	return app.messageOutput
}

// Replaces the MessageOutput, and returns the previous MessageOutput
func (app *Application) SetMessageOutput(out MessageOutput) MessageOutput {
	var prev MessageOutput
	app.mutex.Lock()
	{
		prev = app.messageOutput
		app.messageOutput = out
	}
	app.mutex.Unlock()

	return prev
}

func (app *Application) PrintfMsg(format string, a ...interface{}) {
	app.mutex.Lock()
	out := app.messageOutput
	app.mutex.Unlock()

	out.PrintfMsg(format, a...)
}

// =========
// EventLoop
// =========

type EventLoop struct {
	// The application to which this EventLoop belongs, or nil if
	// this EventLoop was deleted before the whole application
	// terminated.
	app       *Application
	app_mutex sync.RWMutex

	// A symbolic name associated to the EventLoop (useful for
	// debugging).
	Name string

	// If [this channel receives a value] then [this event-loop
	// should pause].  As a response, after this event-loop
	// actually pauses, a value will appear on this channel.
	Pause chan byte

	// Constraint: A value can be sent to this channel only after
	// this event-loop has been paused.
	//
	// If [this channel receives a value] then [this event-loop
	// should terminate].  As a response, after this event-loop
	// actually terminates, a value will appear on this channel.
	Terminate chan byte
}

func (app *Application) NewEventLoop() *EventLoop {
	// By default fill the event name with the caller name.
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	e := &EventLoop{app: app, Name: name, Pause: make(chan byte), Terminate: make(chan byte)}
	ok := app.addEventLoop(e)

	if !ok {
		// Pause and terminate
		go func() {
			if app.VerboseShutdown {
				app.PrintfMsg("application has terminated before %s was created", e.Name)
				app.PrintfMsg("sending pause message to %s", e.Name)
			}
			// Request the event-loop to pause, and wait until it actually pauses
			e.Pause <- 0
			<-e.Pause
			if app.VerboseShutdown {
				app.PrintfMsg("%s is now paused", e.Name)
			}

			if app.VerboseShutdown {
				app.PrintfMsg("sending terminate message to %s", e.Name)
			}
			// Request the event-loop to terminate, and wait until it actually terminates
			e.Terminate <- 0
			<-e.Terminate
			if app.VerboseShutdown {
				app.PrintfMsg("%s has terminated", e.Name)
			}
		}()
	}

	return e
}

func (e *EventLoop) App() *Application {
	e.app_mutex.RLock()
	app := e.app
	e.app_mutex.RUnlock()
	return app
}

// Unregister the EventLoop from the Application.
// When the process finishes, the returned new channel will receive a value.
func (e *EventLoop) Delete() <-chan byte {
	doneCh := make(chan byte)

	e.app_mutex.RLock()
	app := e.app
	e.app_mutex.RUnlock()

	// Remove 'e' from 'app.eventLoops'
	app.mutex.Lock()
	{
		if app.terminationInProgress {
			// Nothing to do here - the EventLoop 'e' will be removed in function 'appGoroutine'
			app.mutex.Unlock()
			go func() { doneCh <- 0 }()
			return doneCh
		}

		found := false
		{
			num_eventLoops := len(app.eventLoops)
			for i := 0; i < num_eventLoops; {
				if app.eventLoops[i] == e {
					// Remove the i-th element
					app.eventLoops[i] = app.eventLoops[num_eventLoops-1]
					app.eventLoops = app.eventLoops[0 : num_eventLoops-1]
					found = true
					break
				} else {
					i++
				}
			}
		}
		if !found {
			panic("no such event-loop")
		}
	}
	app.mutex.Unlock()

	go func() {
		// Request the event-loop to pause, and wait until it actually pauses
		e.Pause <- 0
		<-e.Pause

		// Request the event-loop to terminate, and wait until it actually terminates
		e.Terminate <- 0
		<-e.Terminate

		e.app_mutex.Lock()
		e.app = nil
		e.app_mutex.Unlock()

		doneCh <- 0
	}()
	return doneCh
}

// =============
// MessageOutput
// =============

type MessageOutput interface {
	// Prints a single-line message.
	// If the format string does not end with the new-line character,
	// the new-line character is appended automatically.
	PrintfMsg(format string, a ...interface{})
}

type stdoutMessageOutput struct {
	mutex sync.Mutex
}

func (out *stdoutMessageOutput) PrintfMsg(format string, a ...interface{}) {
	out.mutex.Lock()
	{
		fmt.Printf(format, a...)

		appendNewLine := false
		if (len(format) == 0) || (format[len(format)-1] != '\n') {
			appendNewLine = true
		}

		if appendNewLine {
			fmt.Println()
		}
	}
	out.mutex.Unlock()
}

// ==============
// Misc functions
// ==============

func Drain(ticker *time.Ticker) {
loop:
	for {
		select {
		case <-ticker.C:
		default:
			break loop
		}
	}
}

// ======================
// (Unix) signal handling
// ======================

type SignalHandler interface {
	// Function to be called upon receiving an os.Signal.
	//
	// A single signal is passed to all installed signal handlers.
	// The [order in which this function is called in respect to other handlers] is unspecified.
	HandleSignal(signal os.Signal)
}

// Actually, this is a set
var signalHandlers = make(map[SignalHandler]bool)

var signalHandlers_mutex sync.Mutex

// Installs the specified handler.
// Trying to re-install an already installed handler is effectively a NOOP.
func InstallSignalHandler(handler SignalHandler) {
	signalHandlers_mutex.Lock()
	signalHandlers[handler] = true
	signalHandlers_mutex.Unlock()
}

// Uninstalls the specified handler.
// Trying to uninstall an non-existent handler is effectively a NOOP.
func UninstallSignalHandler(handler SignalHandler) {
	signalHandlers_mutex.Lock()
	delete(signalHandlers, handler)
	signalHandlers_mutex.Unlock()
}

func init() {
	go func() {
		c := make(chan os.Signal, 10)
		signal.Notify(c)
		for sig := range c {
			signalHandlers_mutex.Lock()
			handlers_copy := make([]SignalHandler, len(signalHandlers))
			{
				i := 0
				for handler, _ := range signalHandlers {
					handlers_copy[i] = handler
					i++
				}
			}
			signalHandlers_mutex.Unlock()

			for _, handler := range handlers_copy {
				handler.HandleSignal(sig)
			}
		}
	}()
}
