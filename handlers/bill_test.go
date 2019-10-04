package handlers

import (
	"log"
	"os"
	"work-at-olist/storage"
)

func init() {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	db, err := storage.NewDB("sqlite3", ":memory:")
	if err != nil {
		logger.Fatal(err)
	}

	DB = db.DB

	db.InitSchema()
	h = New(db, logger)
}
