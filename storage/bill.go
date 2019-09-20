package storage

type BillStore interface {
	GetBillByPeriod(string, string, string) (Bill, error)
}

type Bill struct {
	Id         int     `json:"id"`
	Subscriber string  `json:"subscriber"`
	Calls      []Call  `json:"calls"`
	Mouth      string  `json:"mouth"`
	Year       string  `json:"year"`
	Price      float64 `json:"price"`
}

func NewBill(sb string, m string, y string) *Bill {
	return &Bill{
		Subscriber: sb,
		Mouth:      m,
		Year:       y,
	}
}

func (db *DB) GetBillByPeriod(sb string, m string, y string) (Bill, error) {
	var bill Bill
	err := db.First(&bill, "subscriber = ?", sb, m, y).Error
	return bill, err
}
