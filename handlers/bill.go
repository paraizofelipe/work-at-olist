package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"work-at-olist/storage"
)

type Bill struct {
	*storage.Bill
}

func (h *Handler) BillsHandler(w http.ResponseWriter, r *http.Request) {
	router := NewRouter(h.Logger)
	router.AddRoute(
		`bills\/(?P<subscriber>[0-9]+)$`,
		http.MethodGet, h.extractBill(h.getBills))

	router.ServeHTTP(w, r)
}

func (h *Handler) extractBill(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var month, year int64

		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		subscriber, _ := ctx.Value("subscriber").(string)

		p, _ := url.ParseQuery(r.URL.RawQuery)
		if val, ok := p["month"]; ok {
			month, err = strconv.ParseInt(val[0], 10, 64)
			if err != nil {
				return
			}
		}
		if val, ok := p["year"]; ok {
			year, err = strconv.ParseInt(val[0], 10, 64)
			if err != nil {
				return
			}
		}

		now := time.Now()
		if now.Year() == int(year) && int(now.Month()) == int(month) {
			month = month - 1
		}

		bill, err := h.DB.GetBillByPeriod(subscriber, int(month), int(year))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if bill.Id == 0 {
			err = fmt.Errorf("bill not found")
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(bill)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) getBills(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("bills")
}
