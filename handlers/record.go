package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"work-at-olist/storage"
)

type ValidationMessages map[string]interface{}

// This function validates the integrity of information passed via http request.
func (h *Handler) validateRecord(r *storage.Record) (bool, ValidationMessages) {
	var errs = ValidationMessages{}
	var valid = true

	if r.CallId <= 0 {
		errs["call_id"] = []string{"invalid [call_id] field value"}
		valid = false
	}

	if r.Type == "" {
		errs["type"] = []string{"[type] field can't blank"}
		valid = false
	}

	if r.Timestamp == "" {
		errs["timestamp"] = []string{"[timestamp] field can't blank"}
		valid = false
	}

	if r.Type == "start" && (r.Source == "" || !h.validatePhone(r.Source)) {
		errs["source"] = []string{"invalid [source] field value"}
		valid = false
	}

	if r.Type == "start" && (r.Destination == "" || !h.validatePhone(r.Destination)) {
		errs["destination"] = []string{"invalid [destination] field value"}
		valid = false
	}

	return valid, errs
}

// This function validates phone numbers.
func (h *Handler) validatePhone(ph string) bool {
	if _, err := strconv.Atoi(ph); err != nil {
		return false
	}
	if len(ph) > 11 || len(ph) < 10 {
		return false
	}
	return true
}

// This function checks if a record already exists in the database.
func (h *Handler) recordExist(r *storage.Record) (bool, error) {
	calls, err := h.DB.GetRecordsByType(r.CallId, r.Type)
	if err != nil || len(calls) > 0 {
		return false, err
	}
	return true, nil
}

func (h *Handler) RecordsHandler(w http.ResponseWriter, r *http.Request) {
	router := NewRouter(h.Logger)
	router.AddRoute(
		`records\/?$`,
		http.MethodPost, h.postRecord())

	router.ServeHTTP(w, r)
}

// This function saves the call start and end records required for the API.
func (h *Handler) SaveRecord(record *storage.Record) ErrorResponse {
	respErr := ValidationMessages{}

	if valid, errs := h.validateRecord(record); !valid {
		respErr = errs
		return ErrorResponse{http.StatusUnprocessableEntity, respErr}
	}

	if valid, err := h.recordExist(record); err != nil || !valid {
		respErr["error"] = "record already saved"
		return ErrorResponse{http.StatusUnprocessableEntity, respErr}
	}

	if record.Type == "end" {
		rs, err := h.DB.GetRecordsByCallId(record.CallId)
		if err != nil {
			respErr["error"] = "failed to save record"
			return ErrorResponse{http.StatusInternalServerError, respErr}
		}

		if len(rs) < 1 {
			respErr["error"] = "call not started"
			return ErrorResponse{http.StatusUnprocessableEntity, respErr}
		}

		if err := h.DB.CreateRecord(record); err != nil {
			respErr["error"] = err.Error()
			return ErrorResponse{http.StatusInternalServerError, respErr}
		}

		err = h.SaveCall(rs[0], *record)
		if err != nil {
			respErr["error"] = err.Error()
			return ErrorResponse{http.StatusInternalServerError, respErr}
		}
	} else {
		if err := h.DB.CreateRecord(record); err != nil {
			respErr["error"] = err.Error()
			return ErrorResponse{http.StatusInternalServerError, respErr}
		}
	}

	return ErrorResponse{Errors: respErr}
}

// This function captures all call register POST requests.
func (h *Handler) postRecord() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var body struct {
			Type        string `json:"type"`
			Timestamp   string `json:"timestamp"`
			CallId      int    `json:"call_id"`
			Source      string `json:"source"`
			Destination string `json:"destination"`
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		record := storage.NewRecord(body.Type, body.Timestamp, body.CallId, body.Source, body.Destination)

		if errs := h.SaveRecord(record); errs.Status != 0 {
			w.WriteHeader(errs.Status)
			if err = json.NewEncoder(w).Encode(errs); err != nil {
				http.Error(w, "failed to save record", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusCreated)

		resp := map[string]string{"message": "record successfully saved"}

		if err = json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
