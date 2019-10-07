package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"work-at-olist/storage"
)

func TestRecordsHandler_ValidPhone(t *testing.T) {
	tt := []struct {
		in     string
		expect bool
	}{
		{"4199999999", true},
		{"", false},
		{"99999999", false},
		{"419", false},
		{"99999", false},
		{"419999999999999", false},
		{"41", false},
	}

	for _, test := range tt {
		valid := h.validatePhone(test.in)
		if valid != test.expect {
			t.Errorf("validPhone failed expected: %t, received: %t", test.expect, valid)
		}
	}
}

func TestRecordsHandler_IsValid(t *testing.T) {
	tt := []struct {
		in     *storage.Record
		expect bool
	}{
		{
			&storage.Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, true,
		},
		{
			&storage.Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      0,
				Source:      "4199999999",
				Destination: "4288888888",
			}, false,
		},
		{
			&storage.Record{
				Id:          0,
				Type:        "",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, false,
		},
		{
			&storage.Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, false,
		},
		{
			&storage.Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "",
				Destination: "4288888888",
			}, false,
		},
		{
			&storage.Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "",
			}, false,
		},
	}

	for _, test := range tt {
		valid, _ := h.validateRecord(test.in)
		if valid != test.expect {
			t.Errorf("IsValid failed expected: %t, received: %t", test.expect, valid)
		}
	}
}

func TestRecordsHandler_postRecord(t *testing.T) {

	tt := []struct {
		in     storage.Record
		expect int
	}{
		{
			storage.Record{},
			http.StatusUnprocessableEntity,
		},
		{
			storage.Record{
				Type:        "end",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, http.StatusUnprocessableEntity,
		},
		{
			storage.Record{
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      2,
				Source:      "4199999999",
				Destination: "4288888888",
			}, http.StatusCreated,
		},
		{
			storage.Record{
				Type:        "end",
				Timestamp:   "2016-02-29T15:00:00Z",
				CallId:      2,
				Source:      "4199999999",
				Destination: "4288888888",
			}, http.StatusCreated,
		},
		{
			storage.Record{
				Type:        "start",
				Timestamp:   "2016-02-29T15:00:00Z",
				CallId:      2,
				Source:      "4199999999",
				Destination: "4288888888",
			}, http.StatusUnprocessableEntity,
		},
		{
			storage.Record{
				Type:        "end",
				Timestamp:   "2016-02-29T15:00:00Z",
				CallId:      3,
				Source:      "4199999999",
				Destination: "4288888888",
			}, http.StatusUnprocessableEntity,
		},
		{
			storage.Record{
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      3,
				Source:      "41",
				Destination: "4288888888",
			}, http.StatusUnprocessableEntity,
		},
		{
			storage.Record{
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      4,
				Source:      "4199999999",
				Destination: "42",
			}, http.StatusUnprocessableEntity,
		},
	}

	for _, test := range tt {
		jsonBody, _ := json.Marshal(test.in)

		recorder := makeRequest(t, "POST", "/api/records", bytes.NewBuffer(jsonBody), nil)
		if Code := recorder.Code; Code != test.expect {
			t.Errorf("%v should return a %v status code: got %v", test.in, test.expect, Code)
		}
	}

}
