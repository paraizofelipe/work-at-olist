package storage

type BillStore interface {
	GetBill(string) (Bill, error)
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

func (db *DB) GetBill(sb string) (Bill, error) {
	var bill Bill

	rows, err := db.Query(`SELECT * FROM bill WHERE subscriber = ? LIMIT 1;`, sb)
	if err != nil {
		return bill, err
	}

	for rows.Next() {
		if err := rows.Scan(&bill.Id, &bill.Subscriber, &bill.Mouth, &bill.Price); err != nil {
			return bill, err
		}
	}
	return bill, nil
}

func (db *DB) GetBillByPeriod(sb string, m string, y string) (Bill, error) {
	var bill Bill

	if m == "" || y == "" {
		return db.GetBill(sb)
	}

	rows, err := db.Query(`SELECT * FROM bill 
        WHERE subscriber = ?
        AND mouth = ? 
        AND year = ?
        LIMIT 1;`, sb, m, y)
	if err != nil {
		return bill, err
	}

	for rows.Next() {
		if err := rows.Scan(&bill.Id, &bill.Subscriber, &bill.Mouth, &bill.Price); err != nil {
			return bill, err
		}
	}
	return bill, nil
}
