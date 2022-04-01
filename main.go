package main

import (
	"log"
	"net/http"
)





func main() {
	// example usage: curl -s 'http://127.0.0.1:42069/?url=http://go-colly.org/'
	addr := ":7777"

	d := dispatcher{}

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/api/", d.Handler.)

	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
