package main

import (
	"net/http"
	"net"
	"log"
	"time"
	"html/template"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/coreos/go-systemd/daemon"
)

var htmlfile = "index.html"
var folderpath = "templates/"

type Webpage struct {
	Title string
	Body  []byte
}

func renderTemplate(w http.ResponseWriter, title string, folderpath string, page *Webpage) {
	temp, _ := template.ParseFiles(folderpath+title)
 	temp.Execute(w, page)
 }

 // Load webpage, else throw an error.
 func loadWebpage(htmlfile string, folderpath string) (*Webpage, error){
   filename := folderpath+htmlfile
   body, err := ioutil.ReadFile(filename)
   if err != nil {
     log.Println(err)
     return nil, err
   }
   return &Webpage{Title: htmlfile, Body: body}, nil
 }

func pageHandler(w http.ResponseWriter, r *http.Request) {
  page, err := loadWebpage(htmlfile, folderpath)
  if err != nil {
    log.Println(err)
    return
  }
  renderTemplate(w, htmlfile, folderpath, page)
}

func main() {
	r := mux.NewRouter().StrictSlash(true)

	t := http.StripPrefix("/templates/", http.FileServer(http.Dir("templates")))
	r.PathPrefix("/templates/").Handler(t)

	// Valid Routes
	routes := []string{
			"/",
		}

	// Assign Handler for Valid Routes
	for _, value := range routes {
			r.HandleFunc(value, pageHandler)
	}

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
