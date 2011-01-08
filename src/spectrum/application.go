package spectrum

import (
	"container/vector"
	"fmt"
	"sync"
	"time"
	"os/signal"
	"runtime"
)


// ===========
// Application
// ===========

type Application struct {
	exitApp       chan byte // If [this channel is closed] then [the whole application should terminate]
	HasTerminated chan byte // This channel is closed after the whole application has terminated

	// A vector of *EventLoop
	eventLoops vector.Vector

	terminationInProgress bool

	mutex sync.Mutex

	messageOutput MessageOutput

	Verbose bool

	CreationTime int64 // The time when this Application object was created, see time.Nanoseconds()
}

func NewApplication() *Application {
	app := &Application{
		exitApp:       make(chan byte),
		HasTerminated: make(chan byte),
		eventLoops:    vector.Vector{},
		CreationTime:  time.Nanoseconds(),
		messageOutput: &stdoutMessageOutput{},
	}

	go appGoroutine(app)

	return app
}

func appGoroutine(app *Application) {
	// Block until there is a request to exit the application
	<-app.exitApp

	var startTime int64
	if app.Verbose {
		startTime = time.Nanoseconds()
	}

	// Cycle until there are no EventLoop objects.
	// Usually, the body of this 'for' statement executes only once
	for {
		// Make a copy of the 'eventLoops' vector, then clear it
		app.mutex.Lock()
		eventLoops := app.eventLoops.Copy()
		app.eventLoops.Cut(0, app.eventLoops.Len())
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
		for _, x := range eventLoops {
			switch e := x.(type) {
			case *EventLoop:
				// Request the event-loop to pause, and wait until it actually pauses
				e.Pause <- 0
				<-e.Pause
			}
		}

		// Terminate all event loops
		for _, x := range eventLoops {
			switch e := x.(type) {
			case *EventLoop:
				// Request the event-loop to terminate, and wait until it actually terminates
				e.Terminate <- 0
				<-e.Terminate
			}
		}

		app.mutex.Lock()
		if app.eventLoops.Len() == 0 {
			app.mutex.Unlock()
			break
		} else {
			// Some new EventLoops were created while we were pausing&terminating the known ones
			app.mutex.Unlock()
		}
	}

	if app.Verbose {
		endTime := time.Nanoseconds()
		app.PrintfMsg("application shutdown completed after %f milliseconds", float(endTime-startTime)/1e6)
		app.PrintfMsg("application has terminated")
	}

	close(app.HasTerminated)

	app.mutex.Lock()
	app.terminationInProgress = false
	app.mutex.Unlock()
}

func (app *Application) addEventLoop(e *EventLoop) {
	app.mutex.Lock()
	app.eventLoops.Push(e)
	app.mutex.Unlock()
}

func (app *Application) RequestExit() {
	app.mutex.Lock()
	{
		if app.terminationInProgress {
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
	if closed(app.exitApp) {
		panic("cannot create a new event-loop because the application has been terminated")
	}

	// By default fill the event name with the caller name.
	pc, _, _, _ := runtime.Caller(1)
	name := runtime.FuncForPC(pc).Name()

	e := &EventLoop{app: app, Name: name, Pause: make(chan byte), Terminate: make(chan byte)}
	app.addEventLoop(e)
	return e
}

func (e *EventLoop) App() *Application {
	e.app_mutex.RLock()
	app := e.app
	e.app_mutex.RUnlock()
	return app
}

// Unregister the EventLoop from the Application
func (e *EventLoop) Delete() {
	e.app_mutex.RLock()
	app := e.app
	e.app_mutex.RUnlock()

	app.mutex.Lock()
	{
		if app.terminationInProgress {
			// Nothing to do here - the EventLoop 'e' will be removed in function 'appGoroutine'
			app.mutex.Unlock()
			return
		}

		found := false
		{
			for i := 0; i < app.eventLoops.Len(); {
				if app.eventLoops.At(i).(*EventLoop) == e {
					// Remove the i-th element
					app.eventLoops.Swap(i, app.eventLoops.Len()-1)
					app.eventLoops.Delete(app.eventLoops.Len() - 1)
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
	}()
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
	var haveMessage bool
	_, haveMessage = <-ticker.C
	for haveMessage {
		_, haveMessage = <-ticker.C
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
	HandleSignal(signal signal.Signal)
}

// Actually, this is a set
var signalHandlers map[SignalHandler]byte = make(map[SignalHandler]byte)

var signalHandlers_mutex sync.Mutex

// Installs the specified handler.
// Trying to re-install an already installed handler is effectively a NOOP.
func InstallSignalHandler(handler SignalHandler) {
	signalHandlers_mutex.Lock()
	signalHandlers[handler] = 0, true
	signalHandlers_mutex.Unlock()
}

// Uninstalls the specified handler.
// Trying to uninstall an non-existent handler is effectively a NOOP.
func UninstallSignalHandler(handler SignalHandler) {
	signalHandlers_mutex.Lock()
	signalHandlers[handler] = 0, false
	signalHandlers_mutex.Unlock()
}

func init() {
	go func() {
		for {
			signal := <-signal.Incoming

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
				handler.HandleSignal(signal)
			}
		}
	}()
}
