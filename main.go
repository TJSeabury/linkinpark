package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static
var static embed.FS

func main() {
	addr := ":7777"

	d := NewDispatcher()

	contentStatic, _ := fs.Sub(static, "static")
	http.Handle("/", http.FileServer(http.FS(contentStatic)))

	http.HandleFunc("/api/start/", d.Start)
	http.HandleFunc("/api/check/", d.Check)
	http.HandleFunc("/api/finish/", d.Finish)

	//log.Println("listening on", addr)
	http.ListenAndServe(addr, nil)
}
