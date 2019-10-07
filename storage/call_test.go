package storage

import (
	"reflect"
	"testing"
)

type CallFields struct {
	Bid int
	Dst string
	Dur string
	Sd  string
	St  string
	Pri float64
}

func TestNewCall(t *testing.T) {
	tt := []struct {
		in     CallFields
		expect Call
	}{
		{
			CallFields{
				Bid: 0,
				Dst: "4199999999",
				Dur: "2h10m32s",
				Sd:  "2019-10-02",
				St:  "10:00:00",
				Pri: 12.0,
			},
			Call{
				Id:          0,
				BillId:      0,
				Destination: "4199999999",
				Duration:    "2h10m32s",
				StartDate:   "2019-10-02",
				StartTime:   "10:00:00",
				Price:       12.0,
			},
		},
	}

	for _, test := range tt {
		call := NewCall(test.in.Bid, test.in.Dst, test.in.Dur, test.in.Sd, test.in.St, test.in.Pri)
		if !reflect.DeepEqual(call, &test.expect) {
			t.Errorf("NewCall failed expected: %v, received: %v", test.expect, call)
		}
	}
}
