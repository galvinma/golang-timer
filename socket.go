package main

import (
	"log"
	"net/http"
	"time"
	"encoding/json"

  "github.com/gorilla/websocket"
)

// Upgrade to Websocket
var upgrade = websocket.Upgrader{}

// Get Time
func getTime() (t time.Time){
	now := time.Now()
	return now
}

// Websocket Handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Send to client
	for {
		// Wait a Second
		time.Sleep(time.Millisecond)
		// Get Time
		current_time := getTime().Format("Mon Jan _2 15:04:05 2006")

		// Send as JSON
		json, err := json.Marshal(current_time)
		if err != nil {
			log.Println(err)
			}
		err = ws.WriteMessage(websocket.TextMessage, json)
		if err != nil {
			log.Println(err)
			}
		}
}
