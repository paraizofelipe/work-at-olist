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

func TestGetBills() {
	type QueryString struct {
		Subscriber string
		Month      string
		Year       string
	}

	type Expect struct {
		Id int
	}

	tt := []struct {
		in     QueryString
		expect Bill
	}{
		{
			QueryString{
				Subscriber: "",
				Month:      "",
				Year:       "",
			},
			Bill{},
		},
	}

	for _, test := range tt {

	}
}
