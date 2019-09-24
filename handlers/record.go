package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"
	"work-at-olist/storage"
)

type Record struct {
	*storage.Record
}

type errorResponse struct {
	Errors map[string]interface{} `json:"errors"`
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

func inTimeRange(check time.Time) bool {
	start := time.Date(check.Year(), check.Month(), check.Day(), 6, 0, 0, 0, time.UTC)
	end := time.Date(check.Year(), check.Month(), check.Day(), 21, 59, 59, 0, time.UTC)

	return check.After(start) && check.Before(end)
}

func (h *Handler) CalculateCall(dateStart, dateEnd time.Time) (float64, error) {
	var hour int
	var start, end time.Time

	//dateStart, err := time.Parse(time.RFC3339, rs[0].Timestamp)
	//if err != nil {
	//    return 0, err
	//}
	//
	//dateEnd, err := time.Parse(time.RFC3339, rs[1].Timestamp)
	//if err != nil {
	//    return 0, err
	//}

	if !inTimeRange(dateStart) {
		hour = dateStart.Hour()
		if hour < 6 || hour >= 22 {
			dateStart = dateStart.AddDate(0, 0, 1)
		}
		start = time.Date(dateStart.Year(), dateStart.Month(), dateStart.Day(), 6, 0, 0, 0, time.UTC)
	} else {
		start = dateStart
	}

	if !inTimeRange(dateEnd) {
		hour = dateEnd.Hour()
		if hour < 6 {
			dateEnd = dateEnd.AddDate(0, 0, -1)
		}
		end = time.Date(dateEnd.Year(), dateEnd.Month(), dateEnd.Day(), 22, 0, 0, 0, time.UTC)
	} else {
		end = dateEnd
	}

	du := end.Sub(start)
	x := math.Floor(du.Hours()/24) * (8 * 60)
	fmt.Println(du.Hours())

	t := float64(int(du.Minutes()) - int(x))
	rst := math.Round((t*0.09+0.36)*100) / 100

	return rst, nil
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

		dateStart, err := time.Parse(time.RFC3339, rs[0].Timestamp)
		if err != nil {
			return
		}

		dateEnd, err := time.Parse(time.RFC3339, rs[1].Timestamp)
		if err != nil {
			return
		}

		cp, err := h.CalculateCall(dateStart, dateEnd)
		if err != nil {
			return
		}

		duration := dateStart.Sub(dateEnd)
		c := storage.NewCall(body.Destination, duration.String(), dateStart.Format("2006-01-02"), dateStart.Format("3:04:05 PM"), cp)
		if err := h.DB.CreateCall(c); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	}

	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(record); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
