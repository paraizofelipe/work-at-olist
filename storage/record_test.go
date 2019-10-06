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
