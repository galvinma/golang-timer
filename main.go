package main

import (
	"net/http"
	"net"
	"log"
	"time"

	"github.com/gorilla/mux"
	"github.com/coreos/go-systemd/daemon"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	// Templates
	t := http.StripPrefix("/templates/", http.FileServer(http.Dir("templates")))
	r.PathPrefix("/templates/").Handler(t)

	// Static
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	r.PathPrefix("/static/").Handler(s)

	r.HandleFunc("/", pageHandler)
	r.HandleFunc("/ws", wsHandler)

	// Listen
	l, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		log.Println(err)
	}

	// Tell systemd website operational.
	daemon.SdNotify(false, "READY=1")

	// Heartbeat
	go func() {
    interval, err := daemon.SdWatchdogEnabled(false)
    if err != nil || interval == 0 {
        return
    }
		for {
	    _, err := http.Get("http://127.0.0.1:5000")
	    if err == nil {
	        daemon.SdNotify(false, "WATCHDOG=1")
	    }
	    time.Sleep(interval / 3)
		}
	}()
	// Serve
  http.Serve(l, r)

}
