package handlers

import (
	"encoding/json"
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
		http.MethodGet, h.getBills())

	router.ServeHTTP(w, r)
}

func (h *Handler) SearchBill(subscriber string, month, year int64) (storage.Bill, ErrorResponse) {
	var err error
	var bill storage.Bill

	respErr := ValidationMessages{}

	now := time.Now()
	if now.Year() == int(year) && int(now.Month()) == int(month) {
		month = month - 1
	}

	bill, err = h.DB.GetBillByPeriod(subscriber, int(month), int(year))
	if err != nil {
		respErr["message"] = err.Error()
		return bill, ErrorResponse{http.StatusInternalServerError, respErr}
	}

	if bill.Id == 0 {
		respErr["message"] = "bill not found"
		return bill, ErrorResponse{http.StatusNotFound, respErr}
	}

	return bill, ErrorResponse{Errors: respErr}
}

func (h *Handler) ExtractTime(p url.Values) (error, int64, int64) {
	var err error
	var month, year int64

	if val, ok := p["month"]; ok {
		month, err = strconv.ParseInt(val[0], 10, 64)
		if err != nil {
			return err, 0, 0
		}
	}
	if val, ok := p["year"]; ok {
		year, err = strconv.ParseInt(val[0], 10, 64)
		if err != nil {
			return err, 0, 0
		}
	}

	return err, month, year
}

// This function only returns bills from closed months.
func (h *Handler) getBills() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var month, year int64

		ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		subscriber, _ := ctx.Value("subscriber").(string)

		q, _ := url.ParseQuery(r.URL.RawQuery)
		err, month, year = h.ExtractTime(q)
		//if val, ok := p["month"]; ok {
		//	month, err = strconv.ParseInt(val[0], 10, 64)
		//	if err != nil {
		//		http.Error(w, err.Error(), http.StatusInternalServerError)
		//		return
		//	}
		//}
		//if val, ok := p["year"]; ok {
		//	year, err = strconv.ParseInt(val[0], 10, 64)
		//	if err != nil {
		//		http.Error(w, err.Error(), http.StatusInternalServerError)
		//		return
		//	}
		//}

		bill, errResp := h.SearchBill(subscriber, month, year)
		if errResp.Status != 0 {
			w.WriteHeader(errResp.Status)
			if err = json.NewEncoder(w).Encode(errResp); err != nil {
				http.Error(w, "failed to save record", http.StatusInternalServerError)
			}
			return
		}
		//now := time.Now()
		//if now.Year() == int(year) && int(now.Month()) == int(month) {
		//	month = month - 1
		//}
		//
		//bill, err := h.DB.GetBillByPeriod(subscriber, int(month), int(year))
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	return
		//}
		//
		//if bill.Id == 0 {
		//	http.Error(w, "bill not found", http.StatusNotFound)
		//	return
		//}

		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(bill); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		r = r.WithContext(ctx)
	}
}
