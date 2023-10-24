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
	"service-segs/internal/repository"
	"service-segs/internal/service"
	"strconv"
	"time"

	"github.com/caarlos0/env/v9"
)

type dbConfig struct {
	User     string `env:"PG_USER" envDefault:"postgres"`
	Password string `env:"PG_PASSWORD,notEmpty"`
}

type appConfig struct {
	Db            dbConfig
	Port          int    `env:"SERVICE_PORT,notEmpty"`
	CsvVolumePath string `env:"CSV_PATH,required"`
}

func initDatabase(dbCfg dbConfig) *sql.DB {
	database, err := sql.Open(
		"postgres",
		fmt.Sprintf("user=%s password=%s sslmode=disable", dbCfg.User, dbCfg.Password),
	)
	if err != nil {
		log.Fatalln(err)
	}
	return database
}

func main() {
	var appCfg appConfig
	if err := env.Parse(&appCfg); err != nil {
		log.Fatalln(err)
	}

	database := initDatabase(appCfg.Db)

	externalRepository := repository.NewERepository(&struct{}{})
	internalRepository := repository.NewIRepository(database)
	csvVolume := repository.NewCsvRepository(appCfg.CsvVolumePath)

	serviceSegs := service.NewSegmentsService(externalRepository, internalRepository, csvVolume)

	muxer := http.NewServeMux()
	muxer.Handle("/segs", segments.NewHandler(serviceSegs))
	muxer.Handle("/users/", belonging.NewHandler(serviceSegs))
	muxer.Handle("/history/", history.NewHandler(serviceSegs))
	muxer.Handle("/download/", download.NewHandler(serviceSegs))

	server := &http.Server{
		Addr:           ":" + strconv.Itoa(appCfg.Port),
		Handler:        muxer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}

	log.Println("Starting server at", server.Addr, "port")
	log.Fatalln(server.ListenAndServe())
}
