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

func (h *Handler) RecordsToCall(rs []storage.Record) error {
	//t1, err := time.Parse(time.RFC3339, rs[0].Timestamp)
	//if err != nil {
	//	return err
	//}
	//
	//t2, err := time.Parse(time.RFC3339, rs[1].Timestamp)
	//if err != nil {
	//	return err
	//}

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

	record := storage.NewRecord(body.Type, body.Timestamp, body.CallId, body.Source, body.Destination)

	w.Header().Set("Content-Type", "application/json")
	if valid, errs := record.IsValid(); !valid {
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := errorResponse{Errors: errs}

		if err = json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	calls, err := h.DB.GetRecordsByType(record.CallId, record.Type)
	if err != nil || len(calls) > 0 {
		err = fmt.Errorf("call already registered")
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := h.DB.CreateRecord(record); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if record.Type == "end" {
		rs, err := h.DB.GetRecordsByCallId(body.CallId)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}

		h.RecordsToCall(rs)
		c := storage.NewCall(body.Destination, "", 0, "", 0.0)
		if err := h.DB.CreateCall(c); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(record); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
