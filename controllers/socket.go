package controllers

import (
	"container/list"
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

var quit chan bool
var running bool
var testStartTime int64

// Subscriber obj
type Subscriber struct {
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

// Constants for type of message
const (
	MESSAGE = 2
	TOTAL   = 3
	P90     = 4
	P99     = 5
	P50     = 6
	RPS     = 7
)

// Event data structure that is sent to the websocket
type Event struct {
	Type    int
	Content string
}

func setStartTime(time int64) {
	testStartTime = time
}

// New event
func newEvent(ep int, msg string) Event {
	return Event{Type: ep, Content: msg}
}

// Join method called for new connections
func Join(ws *websocket.Conn) {
	subscribe <- Subscriber{Conn: ws}
}

// Leave method to remove the client from socketlist
func Leave(ws *websocket.Conn) {
	unsubscribe <- ws
}

var (
	// Channel for new join users.
	subscribe = make(chan Subscriber, 10)
	// Channel for exit users.
	unsubscribe = make(chan *websocket.Conn, 10)
	// Send events here to publish them.
	publish     = make(chan Event)
	subscribers = list.New()
	users       = 0
)

func broadcaster() {
	for {
		select {
		case sub := <-subscribe:
			subscribers.PushBack(sub) // Add user to the end of list.

		case event := <-publish:
			broadcastWebSocket(event)

		case unsub := <-unsubscribe:
			for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
				if sub.Value.(Subscriber).Conn == unsub {
					subscribers.Remove(sub)
					// Clone connection.
					ws := sub.Value.(Subscriber).Conn
					if ws != nil {
						ws.Close()
					}
					break
				}
			}
		}
	}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event Event) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	data, err := json.Marshal(event)

	if err != nil {
		log.Error("Fail to marshal event:", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- ws
			}
		}
	}
}

func init() {
	go broadcaster()
}
