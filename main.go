package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"work-at-olist/storage"

	"work-at-olist/handlers"
)

const (
	DATABASE string = "work-at-olist.db"
	DIALECT  string = "sqlite3"
	HOST     string = "0.0.0.0"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	db, err := storage.NewDB(DIALECT, DATABASE)
	if err != nil {
		logger.Fatal(err)
	}

	db.InitSchema()

	h := handlers.New(db, logger)

	http.HandleFunc("/api/records", h.RecordsHandler)
	http.HandleFunc("/api/records/", h.RecordsHandler)

	http.HandleFunc("/api/bills", h.BillsHandler)
	http.HandleFunc("/api/bills/", h.BillsHandler)

	url := fmt.Sprintf("%s:%s", HOST, os.Getenv("PORT"))

	log.Printf("🚀 Server listening in %s 🚀", url)

	err = http.ListenAndServe(url, nil)
	if err != nil {
		logger.Fatal(err)
	}
}
