package spectrum

import (
	"container/vector"
	"time"
)

type Application struct {
	exitApp       chan byte // If [this channel is closed] then [the whole application should terminate]
	HasTerminated chan byte // This channel is closed after the whole application has terminated

	// A vector of *EventLoop
	eventLoops vector.Vector

	Verbose bool
}

func NewApplication() *Application {
	app := &Application{make(chan byte), make(chan byte), vector.Vector{}, false}

	go func() {
		select {
		case <-app.exitApp:
			// Make a copy of the 'eventLoops' array
			eventLoops := app.eventLoops.Slice(0, app.eventLoops.Len())

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
			for _, x := range *eventLoops {
				switch e := x.(type) {
				case *EventLoop:
					// Request the event-loop to pause, and wait until it actually pauses
					e.Pause <- 0
					<-e.Pause
				}
			}

			// Terminate all event loops
			for _, x := range *eventLoops {
				switch e := x.(type) {
				case *EventLoop:
					// Request the event-loop to terminate, and wait until it actually terminates
					e.Terminate <- 0
					<-e.Terminate
				}
			}
		}

		if app.Verbose {
			println("application has terminated")
		}
		close(app.HasTerminated)
	}()

	return app
}

func (app *Application) addEventLoop(e *EventLoop) {
	app.eventLoops.Push(e)
}

func (app *Application) RequestExit() {
	close(app.exitApp)
}

type EventLoop struct {
	App *Application

	// If [this channel receives a value] then [this event-loop should pause].
	// As a response, when this event-loop actually pauses, a value will appear on this channel.
	Pause chan byte

	// Constraint: A value can be sent to this channel only after this event-loop was paused.
	//
	// If [this channel receives a value] then [this event-loop should terminate].
	// As a response, when this event-loop actually terminates, a value will appear on this channel.
	Terminate chan byte
}

func (app *Application) NewEventLoop() *EventLoop {
	if closed(app.exitApp) {
		panic("cannot create a new event-loop because the application has been terminated")
	}

	e := &EventLoop{app, make(chan byte), make(chan byte)}
	app.addEventLoop(e)
	return e
}

func Drain(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C: // No action
		default:
			return
		}
	}
}
