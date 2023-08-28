package main

import (
	"log"
	"net/http"
	"service-segs/internal/handler"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/segs", &handler.SegmentsHandler{})
	mux.Handle("/segs/", &handler.UserHandler{})

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}

	log.Println("Starting server at", server.Addr, "port")
	server.ListenAndServe()
}
