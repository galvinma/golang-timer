package main

import (
	"log"
	"net/http"
	"time"
	// "encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

// Upgrade to Websocket
var upgrade = websocket.Upgrader{}

// Get Time
func getTime() (t time.Time) {
	now := time.Now()
	return now
}

func sendData(ws *websocket.Conn, m *sync.Mutex, t []byte) {
	m.Lock()
	err := ws.WriteMessage(websocket.TextMessage, t)
	if err != nil {
		log.Println(err)
	}
	m.Unlock()
}

// Send time to client
func clientSend(ws *websocket.Conn, m *sync.Mutex) {
	for {
		// Wait a Second
		time.Sleep(10 * time.Millisecond)
		current_time := getTime().Format("15:04:05")
		sendData(ws, m, []byte(current_time))
	}
}

func serverRecieve(ws *websocket.Conn, m *sync.Mutex) {
	for {
		_, alarm, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if string(alarm) != "" {
			alarm_time := string(alarm)
			for {
				if alarm_time == getTime().Format("15:04:05") {
					break
				}
			}
			if err != nil {
				log.Println(err)
			}
			sendData(ws, m, []byte("alarm"))

		}
	}
}

// Websocket Handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	var m sync.Mutex
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	go clientSend(ws, &m)
	go serverRecieve(ws, &m)

}
