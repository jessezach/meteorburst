package controllers

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
)

// Subscriber obj
type Subscriber struct {
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

// Event object is sent to the websocket
type Event struct {
	Type    int
	Content string
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

// Send message to websocket
func sendMessage(event Event) {
	publish <- event
}

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
					// Close connection.
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

	data, _ := json.Marshal(event)

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				Leave(ws)
			}
		}
	}
}

func init() {
	go broadcaster()
}
