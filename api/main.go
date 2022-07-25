package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	environment := os.Getenv("environment")
	port := os.Getenv("port")

	log.Println(environment, port)

	mux := http.NewServeMux()

	d := NewDispatcher()

	mux.HandleFunc("/api/start/", d.Start)
	mux.HandleFunc("/api/check/", d.Check)
	mux.HandleFunc("/api/finish/", d.Finish)

	if environment == "production" {
		// CORS is handled by nginx
		http.ListenAndServe(port, mux)
	} else {
		handler := cors.Default().Handler(mux)
		http.ListenAndServe(port, handler)
	}
}
