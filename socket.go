package main

import (
	"log"
	"fmt"
	"strconv"
	"net/http"
	"time"
	"math"
	"sync"

	"github.com/gorilla/websocket"
)

// Globals
var b bool

// Upgrade to Websocket
var upgrade = websocket.Upgrader{}

// Seconds --> Minutes and Seconds
func timeLeft(sec float64) string {
		// var buffer bytes.Buffer
		m := strconv.Itoa(int(sec) / 60)
		s := strconv.Itoa(int(sec) % 60)
		formatted := fmt.Sprintf("%02v:%02v", m, s)
		return formatted
}

// Mutex check before sending data to the client.
// All socket calls through sendData
func sendData(ws *websocket.Conn, m *sync.Mutex, t []byte) {
	m.Lock()
	err := ws.WriteMessage(websocket.TextMessage, t)
	if err != nil {
		log.Println(err)
	}
	m.Unlock()
}

// Show 25:00:00 on app launch
func launchSend(ws *websocket.Conn, m *sync.Mutex) {
		timer := "25:00"
		sendData(ws, m, []byte(timer))
}

// Countdown from 25:00:00
func clientSend(ws *websocket.Conn, m *sync.Mutex) {
	start := time.Now()
	for {
		time.Sleep(100 * time.Millisecond)
		if b == false {
			now := time.Duration(1500)*time.Second
			difference := (now - time.Since(start)).String()
			parse, _ := time.ParseDuration(difference)
			seconds_left := math.Round(parse.Seconds())
			formatted := timeLeft(seconds_left)
			sendData(ws, m, []byte(formatted))
		} else {
			// Reset boolean and break
			b = false
			break
		}
	}
}

// Wait for timer
func waitTimer(ws *websocket.Conn, m *sync.Mutex, timer *time.Timer) {
		<-timer.C
		b = true
		log.Println("Timer expired")
		sendData(ws, m, []byte("alarm"))
}

// Listen for user to set the timer
func serverRecieve(ws *websocket.Conn, m *sync.Mutex) {
	for {
		_, alarm, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if string(alarm) == "start-timer" {
				// Create a new timer. 1500s for pomodoro.
				timer := time.NewTimer(1500*time.Second)
				log.Println("Starting timer")

				go clientSend(ws, m)
				go waitTimer(ws, m, timer)
		}
		// Add code here for stop
	}
}

// Websocket Handler
func wsHandler(w http.ResponseWriter, r *http.Request) {
	var m sync.Mutex
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	launchSend(ws, &m)
	go serverRecieve(ws, &m)

}
