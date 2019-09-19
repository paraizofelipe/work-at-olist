package main

import (
	"log"
	"net/http"
	"os"
	"work-at-olist/storage"

	"work-at-olist/handlers"
)

const (
	DATABASE string = "work-at-olist.db"
	DIALECT  string = "sqlite3"
	PORT     string = ":8989"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	db, err := storage.NewDB(DIALECT, DATABASE)
	if err != nil {
		logger.Fatal(err)
	}

	db.InitSchema()

	h := handlers.New(db, logger)

	http.HandleFunc("/api/calls", h.CallsHandler)
	http.HandleFunc("/api/calls/", h.CallsHandler)

	err = http.ListenAndServe(PORT, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
