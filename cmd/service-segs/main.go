package main

import (
	"database/sql"
	"log"
	"net/http"
	"service-segs/internal/handler/belonging"
	"service-segs/internal/handler/history"
	"service-segs/internal/handler/segments"
	"service-segs/internal/repository"
	"service-segs/internal/service"
	"time"
)

func main() {
	database, err := sql.Open("postgres", "user=postgres password=root sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	externalRepository := repository.NewERepository(&struct{}{})
	internalRepository := repository.NewIRepository(database)

	serviceSegs := service.NewSegmentsService(externalRepository, internalRepository)

	muxer := http.NewServeMux()
	muxer.Handle("/segs", segments.NewHandler(serviceSegs))
	muxer.Handle("/users/", belonging.NewHandler(serviceSegs))
	bonus1 := history.NewHandler(serviceSegs)
	muxer.Handle("/history/", bonus1)
	muxer.Handle("/download/", bonus1)

	server := &http.Server{
		Addr:           ":8080",
		Handler:        muxer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	log.Println("Starting server at", server.Addr, "port")
	log.Fatalln(server.ListenAndServe())
}
