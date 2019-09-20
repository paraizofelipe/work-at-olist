package handlers

import (
	"fmt"
	"net/http"
	"net/url"
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
		ctx := r.Context()

		subscriber, _ := ctx.Value("subscriber").(string)

		p, _ := url.ParseQuery(r.URL.RawQuery)
		mouth := p["year"][0]
		year := p["mouth"][0]

		bill, _ := h.DB.GetBillByPeriod(subscriber, mouth, year)

		fmt.Println(bill.Subscriber)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) getBills(w http.ResponseWriter, r *http.Request) {
	//var err error
	//
	//if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}
	//
	//w.Header().Set("Content-Type", "application/json")
}
