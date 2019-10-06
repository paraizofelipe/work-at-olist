package handlers

import (
	"encoding/json"
	"fmt"
	"testing"
	"work-at-olist/storage"
)

func loadData() {
	src := "99988526423"
	dst := "9933468278"

	var data []*storage.Record

	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 70, Type: "start", Timestamp: "2016-02-29T12:00:00Z"})
	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 70, Type: "end", Timestamp: "2016-02-29T14:00:00Z"})

	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 71, Type: "start", Timestamp: "2017-12-11T15:07:13Z"})
	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 71, Type: "end", Timestamp: "2017-12-11T15:14:56Z"})

	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 72, Type: "start", Timestamp: "2017-12-12T22:47:56Z"})
	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 72, Type: "end", Timestamp: "2017-12-12T22:50:56Z"})

	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 73, Type: "start", Timestamp: "2017-12-12T21:57:13Z"})
	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 73, Type: "end", Timestamp: "2017-12-12T22:10:56Z"})

	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 74, Type: "start", Timestamp: "2017-12-12T04:57:13Z"})
	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 74, Type: "end", Timestamp: "2017-12-12T06:10:56Z"})

	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 75, Type: "start", Timestamp: "2017-12-13T21:57:13Z"})
	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 75, Type: "end", Timestamp: "2017-12-14T22:10:56Z"})

	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 76, Type: "start", Timestamp: "2017-12-12T15:07:58Z"})
	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 76, Type: "end", Timestamp: "2017-12-12T15:12:56Z"})

	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 77, Type: "start", Timestamp: "2018-02-28T21:57:13Z"})
	data = append(data, &storage.Record{Source: src, Destination: dst, CallId: 77, Type: "end", Timestamp: "2018-03-01T22:10:56Z"})

	for _, record := range data {
		h.SaveRecord(record)
	}
}

func TestGetBills(t *testing.T) {
	loadData()
	type QueryString struct {
		Subscriber string
		Month      string
		Year       string
	}

	type Expect struct {
		status int
		bill   storage.Bill
	}

	tt := []struct {
		in     QueryString
		expect Expect
	}{
		{
			QueryString{
				Subscriber: "99988526423",
				Month:      "03",
				Year:       "2018",
			},
			Expect{
				200,
				storage.Bill{
					Subscriber: "99988526423",
					Month:      9,
					Year:       2019,
					Price:      200.0,
				},
			},
		},
	}

	for _, test := range tt {
		url := fmt.Sprintf("/api/bills/%s?month=%s&year=%s", test.in.Subscriber, test.in.Month, test.in.Year)
		recorder := makeRequest(t, "GET", url, nil, nil)

		if recorder.Code != test.expect.status {
			t.Errorf("status code: %d", recorder.Code)
			return
		}

		var bill Bill
		err := json.NewDecoder(recorder.Body).Decode(&bill)

		if err != nil {
			t.Errorf(err.Error())
		}

		if bill.Subscriber != test.expect.bill.Subscriber {
			t.Errorf("should return the correct bill id: got %s want %s", bill.Subscriber, test.expect.bill.Subscriber)
		}
	}
}
