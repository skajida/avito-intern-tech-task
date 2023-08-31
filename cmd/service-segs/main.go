package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"service-segs/internal/handler/belonging"
	"service-segs/internal/handler/download"
	"service-segs/internal/handler/history"
	"service-segs/internal/handler/segments"
	c "service-segs/internal/model/constants"
	"service-segs/internal/repository"
	"service-segs/internal/service"
	"time"
)

func main() {
	database, err := sql.Open("postgres", fmt.Sprintf("%s %s %s", c.PgUser, c.PgPass, c.PgSsl))
	if err != nil {
		log.Fatalln(err)
	}

	externalRepository := repository.NewERepository(&struct{}{})
	internalRepository := repository.NewIRepository(database)

	serviceSegs := service.NewSegmentsService(externalRepository, internalRepository)

	muxer := http.NewServeMux()
	muxer.Handle("/segs", segments.NewHandler(serviceSegs))
	muxer.Handle("/users/", belonging.NewHandler(serviceSegs))
	muxer.Handle("/history/", history.NewHandler(serviceSegs))
	muxer.Handle("/download/", download.NewHandler(serviceSegs))

	server := &http.Server{
		Addr:           c.Port,
		Handler:        muxer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	log.Println("Starting server at", server.Addr, "port")
	log.Fatalln(server.ListenAndServe())
}
