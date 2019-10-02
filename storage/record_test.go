package storage

import (
	"reflect"
	"testing"
)

type RecordFields struct {
	Ty string
	Ti string
	Ci int
	Sr string
	Ds string
}

func TestNewRecord(t *testing.T) {
	tt := []struct {
		in     RecordFields
		expect Record
	}{
		{
			RecordFields{
				Ty: "start",
				Ti: "2016-02-29T14:00:00Z",
				Ci: 99,
				Sr: "418888888",
				Ds: "419999999",
			},
			Record{
				Id:          0,
				CallId:      99,
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				Destination: "419999999",
				Source:      "418888888",
			},
		},
	}

	for _, test := range tt {
		record := NewRecord(test.in.Ty, test.in.Ti, test.in.Ci, test.in.Sr, test.in.Ds)
		if !reflect.DeepEqual(record, &test.expect) {
			t.Errorf("NewRecord failed expected: %v, received: %v", test.expect, record)
		}
	}
}

func TestValidPhone(t *testing.T) {
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
		valid := validPhone(test.in)
		if valid != test.expect {
			t.Errorf("validPhone failed expected: %t, received: %t", test.expect, valid)
		}
	}
}

func TestIsValid(t *testing.T) {
	tt := []struct {
		in     Record
		expect bool
	}{
		{
			Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, true,
		},
		{
			Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      0,
				Source:      "4199999999",
				Destination: "4288888888",
			}, false,
		},
		{
			Record{
				Id:          0,
				Type:        "",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, false,
		},
		{
			Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "",
				CallId:      1,
				Source:      "4199999999",
				Destination: "4288888888",
			}, false,
		},
		{
			Record{
				Id:          0,
				Type:        "start",
				Timestamp:   "2016-02-29T14:00:00Z",
				CallId:      1,
				Source:      "",
				Destination: "4288888888",
			}, false,
		},
		{
			Record{
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
		valid, _ := test.in.IsValid()
		if valid != test.expect {
			t.Errorf("IsValid failed expected: %t, received: %t", test.expect, valid)
		}
	}
}
