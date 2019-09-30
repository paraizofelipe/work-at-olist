package storage

type CallStorer interface {
	CreateCall(*Call) error
}

type Call struct {
	Id          int     `json:"id"`
	BillId      int     `json:"bill_id"`
	Destination string  `json:"destination"`
	Duration    string  `json:"duration"`
	StartDate   string  `json:"start_date"`
	StartTime   string  `json:"start_time"`
	Price       float64 `json:"price"`
}

func NewCall(bid int, dst string, dur string, sd string, st string, pri float64) *Call {
	return &Call{
		BillId:      bid,
		Destination: dst,
		Duration:    dur,
		StartDate:   sd,
		StartTime:   st,
		Price:       pri,
	}
}

func (db *DB) GetCallsByBillId(billId int) ([]Call, error) {
	var call Call
	var calls []Call

	rows, err := db.Query(`SELECT * FROM call WHERE bill_id = ?`, billId)
	if err != nil {
		return calls, err
	}

	for rows.Next() {
		if err := rows.Scan(&call.Id, &call.BillId, &call.Destination, &call.Duration, &call.StartDate, &call.StartTime, &call.Price); err != nil {
			return calls, err
		}
		calls = append(calls, call)
	}
	return calls, err
}

func (db *DB) CreateCall(call *Call) error {
	statement, err := db.Prepare(`INSERT INTO call 
        (bill_id, destionation, duration, start_date, start_time, price) 
        VALUES (?, ?, ?, ?, ?, ?);`)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(
		call.BillId,
		call.Destination,
		call.Duration,
		call.StartDate,
		call.StartTime,
		call.Price); err != nil {

		return err
	}

	return nil
}
