package main

import (
	"net/http"

	"github.com/rs/cors"
)

func main() {
	addr := ":7777"

	mux := http.NewServeMux()

	d := NewDispatcher()

	mux.HandleFunc("/api/start/", d.Start)
	mux.HandleFunc("/api/check/", d.Check)
	mux.HandleFunc("/api/finish/", d.Finish)

	handler := cors.Default().Handler(mux)
	http.ListenAndServe(addr, handler)
}
