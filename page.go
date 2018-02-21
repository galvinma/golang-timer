package main

import (
	"net/http"
	"log"
	"html/template"
	"io/ioutil"

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
