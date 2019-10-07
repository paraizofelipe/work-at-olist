package storage

import (
	"reflect"
	"testing"
)

type BillFields struct {
	Sb string
	M  int
	Y  int
}

func TestNewBill(t *testing.T) {
	tt := []struct {
		in     BillFields
		expect Bill
	}{
		{
			BillFields{
				Sb: "4199999999",
				M:  10,
				Y:  2019,
			},
			Bill{
				Subscriber: "4199999999",
				Month:      10,
				Year:       2019,
				Price:      0,
			},
		},
	}

	for _, test := range tt {
		bill := NewBill(test.in.Sb, test.in.M, test.in.Y, 0)
		if !reflect.DeepEqual(bill, &test.expect) {
			t.Errorf("NewBillfailed expected: %v, received: %v", test.expect, bill)
		}
	}
}
