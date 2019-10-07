package storage

type BillStorer interface {
	CreateBill(*Bill) (int64, error)
	GetBill(string) (Bill, error)
	GetBillByPeriod(string, int, int) (Bill, error)
	ChangePrice(int, float64) error
}

type Bill struct {
	Id         int     `json:"-"`
	Subscriber string  `json:"subscriber"`
	Month      int     `json:"month"`
	Year       int     `json:"year"`
	Calls      []Call  `json:"calls"`
	Price      float64 `json:"price"`
}

func NewBill(sb string, m int, y int, p float64) *Bill {
	return &Bill{
		Subscriber: sb,
		Month:      m,
		Year:       y,
		Price:      p,
	}
}

func (db *DB) CreateBill(bill *Bill) (int64, error) {
	statement, err := db.Prepare(`INSERT INTO bill (subscriber, month, year, price) 
        VALUES (?, ?, ?, ?);`)
	if err != nil {
		return 0, err
	}

	result, err := statement.Exec(bill.Subscriber, bill.Month, bill.Year, bill.Price)
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	lastId, err := result.LastInsertId()

	return lastId, nil
}

func (db *DB) ChangePrice(id int, price float64) error {
	statement, err := db.Prepare(`UPDATE bill SET price = ? WHERE id = ?`)
	if err != nil {
		return err
	}

	_, err = statement.Exec(price, id)
	if err != nil {
		return err
	}
	defer statement.Close()

	return nil
}

func (db *DB) GetBill(sb string) (Bill, error) {
	var bill Bill

	rows, err := db.Query(`SELECT * FROM bill WHERE subscriber = ? LIMIT 1;`, sb)
	if err != nil {
		return bill, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&bill.Id, &bill.Subscriber, &bill.Month, &bill.Year, &bill.Price); err != nil {
			return bill, err
		}

		bill.Calls, err = db.GetCallsByBillId(bill.Id)
		if err != nil {
			return bill, err
		}
	}
	return bill, nil
}

func (db *DB) GetBillByPeriod(sb string, m int, y int) (Bill, error) {
	var bill Bill

	if m == 0 || y == 0 {
		return db.GetBill(sb)
	}

	rows, err := db.Query(`SELECT * FROM bill 
        WHERE subscriber = ?
        AND month = ? 
        AND year = ?
        LIMIT 1;`, sb, m, y)
	if err != nil {
		return bill, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&bill.Id, &bill.Subscriber, &bill.Month, &bill.Year, &bill.Price); err != nil {
			return bill, err
		}

		bill.Calls, err = db.GetCallsByBillId(bill.Id)
		if err != nil {
			return bill, err
		}
	}
	return bill, nil
}
