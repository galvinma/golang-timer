package main

import (
	"log"
	"fmt"
	"strconv"
	"net/http"
	"time"
	"sync"

	"github.com/a-h/round"
	"github.com/gorilla/websocket"
)

// Globals
// When trigger is true, sound the alarm
// When stop is true, stop the timer
var alarmtrigger bool
var alarmstop bool
var countdown bool
var savetimer bool
var pomo_duration = 1500
var saved_duration int
// var short_break_duration = 300
// var long_break_duration = 1500

// Upgrade to Websocket
var upgrade = websocket.Upgrader{}

// Helper function to turn seconds into human readable minutes + seconds.
func timeLeft(sec float64) string {
		// var buffer bytes.Buffer
		m := strconv.Itoa(int(sec) / 60)
		s := strconv.Itoa(int(sec) % 60)
		// Seconds --> Minutes and Seconds
		formatted := fmt.Sprintf("%02v:%02v", m, s)
		return formatted
}

// Show 25:00:00 on app launch
func launchSend(ws *websocket.Conn, m *sync.Mutex) {
		timer := "25:00"
		sendData(ws, m, []byte(timer))
}

// Mutex check before sending data to the client.
// All socket calls to client funneled through sendData
func sendData(ws *websocket.Conn, m *sync.Mutex, t []byte) {
	m.Lock()
	err := ws.WriteMessage(websocket.TextMessage, t)
	if err != nil {
		log.Println(err)
	}
	m.Unlock()
}

type Save struct {
	Save float64
}
// Function sends the time left in session to the client
// Set alarmtrigger to true if time runs out
// Countdown from 25:00:00
func clientSend(ws *websocket.Conn, m *sync.Mutex) {
	start := time.Now()
	for {
			time.Sleep(100 * time.Millisecond)
			now := time.Duration(pomo_duration)*time.Second
			difference := (now - time.Since(start)).String()
			parse, _ := time.ParseDuration(difference)
			seconds_left := round.AwayFromZero(parse.Seconds(), 0)
			saved_duration = int(seconds_left)
			if seconds_left <= 0 {
				alarmtrigger = true
			}
			formatted := timeLeft(seconds_left)
			if countdown == true {
				sendData(ws, m, []byte(formatted))
			}	else {
				break
			}
	}
}

// Wait for timer
func waitTimer(ws *websocket.Conn, m *sync.Mutex, timer *time.Timer) {
		for {
			// If no alarm and no stop, keep listening
			if alarmtrigger  == false && alarmstop == false {
					// wait for something to happen
					time.Sleep(100 * time.Millisecond)
			// if alarm trigger is false, but user issued stop command
			} else if alarmstop == true {
					// stop and do not send alarm. reset booleans
					log.Println("User stopped the timer.")
					countdown = false;
					break
			// if alarm trigger is true, sound the alarm
			} else if alarmtrigger == true {
					log.Println("Timer expired. Sounding the alarm.")
					sendData(ws, m, []byte("alarm"))
					break
			}
		}
		// reset
		alarmtrigger = false
		alarmstop = false
		countdown = false
}

// Listen for user to set the timer
func serverRecieve(ws *websocket.Conn, m *sync.Mutex) {
	for {
		// sleep
		time.Sleep(100 * time.Millisecond)
		// read
		_, alarm, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if string(alarm) == "start-timer" {
				// pomo_duration = 1500
				pomo_dtime := time.Duration(pomo_duration) * time.Second
				timer := time.NewTimer(pomo_dtime)
				log.Println("Starting timer")
				// reset the booleans
				alarmtrigger = false
				alarmstop = false
				countdown = true
				go clientSend(ws, m)
				go waitTimer(ws, m, timer)
		}
		if string(alarm) == "stop-timer" {
				alarmstop = true
				savetimer = true
		}
		if string(alarm) == "resume-timer" {
				alarmtrigger = false
				alarmstop = false
				countdown = true
				log.Println("Resuming timer")
				pomo_duration = saved_duration
				pomo_dtime := time.Duration(pomo_duration) * time.Second
				timer := time.NewTimer(pomo_dtime)
				go clientSend(ws, m)
				go waitTimer(ws, m, timer)
		}
		if string(alarm) == "reset-timer" {
				log.Println("Reset timer")
				alarmstop = true
		}
		if string(alarm) == "25" {
				pomo_duration = 1500
		}
		if string(alarm) == "10" {
				pomo_duration = 600
		}
		if string(alarm) == "5" {
				pomo_duration = 300
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
	launchSend(ws, &m)
	go serverRecieve(ws, &m)

}
