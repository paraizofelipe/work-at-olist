package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"work-at-olist/storage"
)

type Call struct {
	*storage.Record
}

type errorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

type CallsJSON struct {
	Calls      []Call `json:"calls"`
	CallsCount int    `json:"callsCount"`
}

func (h *Handler) buildCallJSON(c *storage.Record) Call {
	call := Call{}
	call.Id = c.Id
	call.Type = c.Type
	call.CallId = c.CallId
	call.Source = c.Source
	call.Destination = c.Destination
	call.Timestamp = c.Timestamp

	return call
}

func (h *Handler) RecordsHandler(w http.ResponseWriter, r *http.Request) {
	router := NewRouter(h.Logger)
	router.AddRoute(
		`records\/?$`,
		http.MethodGet, h.setContext(h.getAllRecords))

	router.AddRoute(
		`records\/?$`,
		http.MethodPost, h.setContext(h.SaveRecord))

	router.ServeHTTP(w, r)
}

func (h *Handler) setContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) SaveRecord(w http.ResponseWriter, r *http.Request) {
	var err error
	var body struct {
		Type        string `json:"type"`
		Timestamp   string `json:"timestamp"`
		CallId      int    `json:"call_id"`
		Source      string `json:"source"`
		Destination string `json:"destination"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()

	c := storage.NewRecord(body.Type, body.Timestamp, body.CallId, body.Source, body.Destination)

	w.Header().Set("Content-Type", "application/json")
	if valid, errs := c.IsValid(); !valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := errorResponse{Errors: errs}

		if err = json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	//TODO validate multiple call terminations
	if c.Type == "start" {
		if calls, err := h.DB.GetRecordsByCallId(c.CallId); err != nil || len(calls) > 0 {
			err = fmt.Errorf("call already initialized")
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}

	if err := h.DB.CreateRecord(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) getAllRecords(w http.ResponseWriter, r *http.Request) {
	var err error
	var calls []storage.Record

	query := h.DB.GetAllRecords()

	err = query.Find(&calls).Error

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if len(calls) == 0 {
		err = json.NewEncoder(w).Encode(CallsJSON{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	err = json.NewEncoder(w).Encode(calls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
