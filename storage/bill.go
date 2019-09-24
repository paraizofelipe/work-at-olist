package storage

type BillStorer interface {
	CreateBill(*Bill) (int64, error)
	GetBill(string) (Bill, error)
	GetBillByPeriod(string, int, int) (Bill, error)
}

type Bill struct {
	Id         int     `json:"id"`
	Subscriber string  `json:"subscriber"`
	Mouth      int     `json:"mouth"`
	Year       int     `json:"year"`
	Price      float64 `json:"price"`
}

func NewBill(sb string, m int, y int) *Bill {
	return &Bill{
		Subscriber: sb,
		Mouth:      m,
		Year:       y,
	}
}

func (db *DB) CreateBill(bill *Bill) (int64, error) {
	statement, _ := db.Prepare(`INSERT INTO bill (subscriber, mouth, year, price) 
        VALUES (?, ?, ?, ?);`)

	result, err := statement.Exec(bill.Subscriber, bill.Mouth, bill.Year, bill.Price)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	return lastId, nil
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

func (db *DB) GetBillByPeriod(sb string, m int, y int) (Bill, error) {
	var bill Bill

	if m == 0 || y == 0 {
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
