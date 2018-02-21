package main

import (
	"net/http"
	"net"
	"log"

	"github.com/gorilla/mux"
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

	// Serve
  http.Serve(l, r)

}
