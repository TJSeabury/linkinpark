package main

import (
	"log"
	"net/http"
)

func main() {
	// example usage: curl -s 'http://127.0.0.1:42069/?url=http://go-colly.org/'
	addr := ":7777"

	d := NewDispatcher()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/api/start", d.Start)
	http.HandleFunc("/api/check", d.Check)
	http.HandleFunc("/api/finish", d.Finish)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
