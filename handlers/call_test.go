package handlers

import (
	"testing"
	"time"
	"work-at-olist/storage"
)

func TestCallHandler_validateCallTimeRange(t *testing.T) {
	tt := []struct {
		in     time.Time
		expect bool
	}{
		{
			time.Date(2019, 10, 03, 6, 0, 0, 0, time.UTC),
			true,
		},
		{
			time.Date(2019, 10, 03, 6, 1, 0, 0, time.UTC),
			true,
		},
		{
			time.Date(2019, 10, 03, 16, 10, 0, 0, time.UTC),
			true,
		},
		{
			time.Date(2019, 10, 03, 15, 10, 0, 0, time.UTC),
			true,
		},
		{
			time.Date(2019, 10, 03, 23, 10, 0, 0, time.UTC),
			false,
		},
		{
			time.Date(2019, 10, 03, 22, 0, 0, 0, time.UTC),
			true,
		},
		{
			time.Date(2019, 10, 03, 22, 1, 0, 0, time.UTC),
			false,
		},
		{
			time.Date(2019, 10, 03, 04, 0, 0, 0, time.UTC),
			false,
		},
		{
			time.Date(2019, 10, 03, 05, 59, 0, 0, time.UTC),
			false,
		},
	}

	for _, test := range tt {
		if valid := h.callInTimeRange(test.in); valid != test.expect {
			t.Errorf("callInTimeRange %v failed expected: %v, received: %v", test.in, test.expect, valid)
		}
	}
}

func TestCallHandler_calculateCall(t *testing.T) {
	tt := []struct {
		in     [2]time.Time
		expect float64
	}{
		{
			[2]time.Time{
				time.Date(2019, 10, 04, 6, 0, 0, 0, time.UTC),
				time.Date(2019, 10, 04, 22, 0, 0, 0, time.UTC),
			},
			86.76,
		},
		{
			[2]time.Time{
				time.Date(2018, 02, 28, 21, 57, 13, 0, time.UTC),
				time.Date(2018, 03, 01, 22, 10, 56, 0, time.UTC),
			},
			86.94,
		},
		{
			[2]time.Time{
				time.Date(2018, 02, 28, 21, 57, 13, 0, time.UTC),
				time.Date(2018, 02, 28, 22, 17, 53, 0, time.UTC),
			},
			0.54,
		},
		{
			[2]time.Time{
				time.Date(2018, 03, 28, 23, 15, 13, 0, time.UTC),
				time.Date(2018, 03, 29, 05, 00, 00, 0, time.UTC),
			},
			0.36,
		},
		{
			[2]time.Time{
				time.Date(2017, 12, 12, 04, 57, 13, 0, time.UTC),
				time.Date(2017, 12, 12, 06, 10, 56, 0, time.UTC),
			},
			1.26,
		},
	}

	for _, test := range tt {
		price, err := h.calculateCallTime(test.in[0], test.in[1])
		if err != nil {
			t.Errorf(err.Error())
		}
		if price != test.expect {
			t.Errorf("calculateCallTime %v failed expected: %f, received: %f", test.in, test.expect, price)
		}
	}
}

func TestCallHandler_SaveCall(t *testing.T) {
	tt := []struct {
		in     [2]storage.Record
		expect error
	}{
		{
			[2]storage.Record{
				{
					Source:      "4199999999",
					Destination: "4288888888",
					Timestamp:   "2016-02-29T14:00:00Z",
				},
				{
					Source:      "4199999999",
					Destination: "4288888888",
					Timestamp:   "2016-02-29T14:00:00Z",
				},
			},
			nil,
		},
		{
			[2]storage.Record{
				{
					Source:      "4199999999",
					Destination: "",
					Timestamp:   "2016-02-29T14:00:00Z",
				},
				{
					Source:      "4199999999",
					Destination: "",
					Timestamp:   "2016-03-01T14:00:00Z",
				},
			},
			nil,
		},
	}

	for _, test := range tt {
		if err := h.SaveCall(test.in[0], test.in[1]); err != test.expect {
			t.Errorf("SaveCall %v failed expected: %v, received: %v", test.in, test.expect, err)
		}
	}
}
