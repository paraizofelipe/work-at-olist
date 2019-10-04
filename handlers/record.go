package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"work-at-olist/storage"
)

type errorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

func (h *Handler) RecordsHandler(w http.ResponseWriter, r *http.Request) {
	router := NewRouter(h.Logger)
	router.AddRoute(
		`records\/?$`,
		http.MethodPost, h.SaveRecord())

	router.ServeHTTP(w, r)
}

func inTimeRange(check time.Time) bool {
	start := time.Date(check.Year(), check.Month(), check.Day(), 5, 59, 59, 0, time.UTC)
	end := time.Date(check.Year(), check.Month(), check.Day(), 22, 1, 0, 0, time.UTC)

	return check.After(start) && check.Before(end)
}

func (h *Handler) SaveRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var body struct {
			Type        string `json:"type"`
			Timestamp   string `json:"timestamp"`
			CallId      int    `json:"call_id"`
			Source      string `json:"source"`
			Destination string `json:"destination"`
		}

		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

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
				http.Error(w, "failed to register call", http.StatusInternalServerError)
			}
			return
		}

		calls, err := h.DB.GetRecordsByType(record.CallId, record.Type)
		if err != nil || len(calls) > 0 {
			http.Error(w, "call already registered", http.StatusUnprocessableEntity)
			return
		}

		if err := h.DB.CreateRecord(record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if record.Type == "end" {
			rs, err := h.DB.GetRecordsByCallId(body.CallId)
			if err != nil {
				http.Error(w, "failed to register call", http.StatusInternalServerError)
				return
			}

			if len(rs) < 2 {
				http.Error(w, "call not started", http.StatusUnprocessableEntity)
				return
			}

			err = h.SaveCall(rs[0], rs[1])
			if err != nil {
				return
			}
		}

		w.WriteHeader(http.StatusCreated)

		if err = json.NewEncoder(w).Encode(record); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		r = r.WithContext(ctx)
	}
}
