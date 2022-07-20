package main

import (
	"net/http"
)

func main() {
	addr := ":7777"

	d := NewDispatcher()

	http.HandleFunc("/api/start/", d.Start)
	http.HandleFunc("/api/check/", d.Check)
	http.HandleFunc("/api/finish/", d.Finish)

	//log.Println("listening on", addr)
	http.ListenAndServe(addr, nil)
}
