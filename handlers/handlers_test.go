package handlers

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"work-at-olist/storage"
)

var (
	h  *Handler
	DB *sql.DB
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
	h.BillsHandler(recorder, req)

	return recorder
}

func init() {

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	db, err := storage.NewDB("sqlite3", "../test.db")
	if err != nil {
		logger.Fatal(err)
	}

	DB = db.DB

	db.CleanDatabase()
	db.InitSchema()

	h = New(db, logger)
}
