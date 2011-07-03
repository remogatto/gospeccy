package websocket

import (
	"http"
	"websocket"
	"spectrum"
	"json"
	"path"
	"os"
)

var (
	publicPath string
	renderComplete chan bool
	fwdScreenChannel chan *spectrum.DisplayData
	app *spectrum.Application
)

type WebSocketService struct {
	// Channel for receiving display changes
	screenChannel chan *spectrum.DisplayData
}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path[1:]
	if p == "" {
		p = "index.html"
	}
	p = path.Join(publicPath, p)
	http.ServeFile(w, r, p)
}

func sockHandler(ws *websocket.Conn) {
	// Start the read/write loops
	go readMessages(ws)
	writeMessages(ws)
}

func readMessages(ws *websocket.Conn) {
	for {
		msg := make([]byte, 10)
		if n, err := ws.Read(msg); err != nil {
			if err == os.EOF {
				return
			}
		} else {
			switch string(msg[0:n]) {
			case "RECEIVED":
				renderComplete <- true
			case "CLOSE":
				return
			}
		} 
	}
}

func writeMessages(ws *websocket.Conn) {
	for {
		select {
		case screen := <-fwdScreenChannel:
			var jsonString string
			if screen != nil {
				if data, err := json.Marshal(screen.Bitmap); err != nil {
					panic(err)
				} else {
					jsonString += "{\"bitmap\":" + string(data)
				}
				if data, err := json.Marshal(screen.Attr); err != nil {
					panic(err)
				} else {
					jsonString += ",\"attr\":" + string(data)
				}
				if data, err := json.Marshal(screen.Dirty); err != nil {
					panic(err)
				} else {
					jsonString += ",\"dirty\":" + string(data) + "}"
				}
				ws.Write([]byte(jsonString))
			} else {
				println("NIL")
			}
		}
	}
}

func NewWebSocketService(a *spectrum.Application, addr, path string) *WebSocketService {
	screen := &WebSocketService{
		screenChannel: make(chan *spectrum.DisplayData),
	}

	app = a
	publicPath = path

	renderComplete = make(chan bool)
	fwdScreenChannel = make(chan *spectrum.DisplayData)

	// Start the service
	go func() {
		http.HandleFunc("/", staticHandler)
		http.Handle("/gospeccy", websocket.Handler(sockHandler))
		if err := http.ListenAndServe(addr, nil); err != nil {
			panic("ListenAndServe: " + err.String())
		}
	}()

	// Listen to incoming DisplayData messages from the emulation
	// core forwarding them to the websocket service
	go webSocketForwarderLoop(screen.screenChannel)

	return screen
}

// Implement DisplayReceiver
func (display *WebSocketService) GetDisplayDataChannel() chan<- *spectrum.DisplayData {
	return display.screenChannel
}

func (display *WebSocketService) Close() {
	display.screenChannel <- nil
}

func webSocketForwarderLoop(screenChannel <-chan *spectrum.DisplayData) {
	for {
		select {
		case screen := <-screenChannel:
			if screen != nil {
				fwdScreenChannel <- screen
				<-renderComplete
			}
		} 
	}
}
