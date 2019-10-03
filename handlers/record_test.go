package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"work-at-olist/storage"
)

var (
	h       *Handler
	DB      *sql.DB
	records []*storage.Record
)

func makeRequest(t *testing.T, method string, url string, body io.Reader, header http.Header) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	}
	if header != nil {
		req.Header = header
	}

	var recorder = httptest.NewRecorder()
	h.RecordsHandler(recorder, req)

	return recorder
}

func TestMain(m *testing.M) {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	db, err := storage.NewDB("sqlite3", ":memory:")
	if err != nil {
		logger.Fatal(err)
	}

	DB = db.DB

	db.InitSchema()
	h = New(db, logger)

	exit := m.Run()

	os.Exit(exit)
}

func TestRecordsHandler_SaveRecord(t *testing.T) {

	tt := []struct {
		in     storage.Record
		expect int
	}{
		{
			storage.Record{
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, http.StatusCreated,
		},
		{
			storage.Record{
				Type:        "end",
				Timestamp:   "2016-02-29T15:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, http.StatusCreated,
		},
		{
			storage.Record{
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "41",
				Destination: "4288888888",
			}, http.StatusUnprocessableEntity,
		},
		{
			storage.Record{
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "42",
			}, http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tt {
		jsonBody, _ := json.Marshal(test.in)

		recorder := makeRequest(t, "POST", "/api/records", bytes.NewBuffer(jsonBody), nil)
		if Code := recorder.Code; Code != test.expect {
			t.Errorf("should return a %v status code: got %v", test.expect, Code)
		}
	}

}
