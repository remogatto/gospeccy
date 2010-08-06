package spectrum

import (
	"container/vector"
	"sync"
	"time"
)


// ===========
// Application
// ===========

type Application struct {
	exitApp       chan byte // If [this channel is closed] then [the whole application should terminate]
	HasTerminated chan byte // This channel is closed after the whole application has terminated

	// A vector of *EventLoop
	eventLoops       vector.Vector
	eventLoops_mutex sync.Mutex

	Verbose bool
}

func NewApplication() *Application {
	app := &Application{exitApp: make(chan byte), HasTerminated: make(chan byte), eventLoops: vector.Vector{}, Verbose: false}

	go appGoroutine(app)

	return app
}

func appGoroutine(app *Application) {
	// Block until there is a request to exit the application
	<-app.exitApp

	// Cycle until there are no EventLoop objects.
	// Usually, the body of this 'for' statement executes only once
	for {
		// Make a copy of the 'eventLoops' vector, then clear it
		app.eventLoops_mutex.Lock()
		eventLoops := app.eventLoops.Copy()
		app.eventLoops.Cut(0, app.eventLoops.Len())
		app.eventLoops_mutex.Unlock()

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

		app.eventLoops_mutex.Lock()
		if app.eventLoops.Len() == 0 {
			app.eventLoops_mutex.Unlock()
			break
		} else {
			// Some new EventLoops were created while we were pausing&terminating the known ones
			app.eventLoops_mutex.Unlock()
		}
	}

	if app.Verbose {
		println("application has terminated")
	}
	close(app.HasTerminated)
}

func (app *Application) addEventLoop(e *EventLoop) {
	app.eventLoops_mutex.Lock()
	app.eventLoops.Push(e)
	app.eventLoops_mutex.Unlock()
}

func (app *Application) RequestExit() {
	close(app.exitApp)
}


// =========
// EventLoop
// =========

type EventLoop struct {
	// The application to which this EventLoop belongs,
	// or nil if this EventLoop was deleted before the whole application terminated.
	app       *Application
	app_mutex sync.RWMutex

	// If [this channel receives a value] then [this event-loop should pause].
	// As a response, after this event-loop actually pauses, a value will appear on this channel.
	Pause chan byte

	// Constraint: A value can be sent to this channel only after this event-loop has been paused.
	//
	// If [this channel receives a value] then [this event-loop should terminate].
	// As a response, after this event-loop actually terminates, a value will appear on this channel.
	Terminate chan byte
}

func (app *Application) NewEventLoop() *EventLoop {
	if closed(app.exitApp) {
		panic("cannot create a new event-loop because the application has been terminated")
	}

	e := &EventLoop{app: app, Pause: make(chan byte), Terminate: make(chan byte)}
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

	app.eventLoops_mutex.Lock()
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
	app.eventLoops_mutex.Unlock()

	if !found {
		panic("no such event-loop")
	}

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


// ==============
// Misc functions
// =============

func Drain(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C: // No action
		default:
			return
		}
	}
}
