package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"work-at-olist/storage"
)

type Call struct {
	Id          int    `json:"id"`
	Type        string `json:"type"`
	Timestamp   string `json:"timestamp"`
	CallId      int    `json:"call_id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type errorResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

type CallJSON struct {
	Call `json:"call"`
}

type CallsJSON struct {
	Calls      []Call `json:"calls"`
	CallsCount int    `json:"callsCount"`
}

func (h *Handler) buildCallJSON(c *storage.Call) Call {
	call := Call{
		Id:          c.ID,
		Type:        c.Type,
		CallId:      c.CallId,
		Source:      c.Source,
		Destination: c.Destination,
		Timestamp:   c.Timestamp,
	}

	return call
}

func (h *Handler) CallsHandler(w http.ResponseWriter, r *http.Request) {
	router := NewRouter(h.Logger)
	router.AddRoute(
		`calls\/?$`,
		http.MethodGet, h.setContext(h.getCall))

	router.AddRoute(
		`calls\/(?P<id>[0-9a-zA-Z\-]+)$`,
		http.MethodGet, h.setContext(h.getAllCalls))

	router.AddRoute(
		`calls\/?$`,
		http.MethodPost, h.setContext(h.createCall))

	router.ServeHTTP(w, r)
}

func (h *Handler) setContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) getCall(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "AllCalls")
	if err != nil {
		log.Fatal(err)
	}
}

func (h *Handler) createCall(w http.ResponseWriter, r *http.Request) {
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

	c := storage.NewCall(body.Type, body.Timestamp, body.CallId, body.Source, body.Destination)

	if valid, errs := c.IsValid(); !valid {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		errorResponse := errorResponse{Errors: errs}

		err = json.NewEncoder(w).Encode(errorResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := h.DB.CreateCall(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	callJSON := CallJSON{
		Call: h.buildCallJSON(c),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(callJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) getAllCalls(w http.ResponseWriter, r *http.Request) {
	var err error
	var calls []storage.Call

	query := h.DB.GetAllCalls()

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
