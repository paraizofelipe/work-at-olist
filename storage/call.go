package storage

type CallStorer interface {
	CreateCall(*Call) error
}

type Call struct {
	Id          string  `json:"id"`
	Destination string  `json:"destination"`
	Duration    string  `json:"duration"`
	StartDate   string  `json:"start_date"`
	StartTime   string  `json:"start_time"`
	Price       float64 `json:"price"`
}

func NewCall(dst string, dur string, sd string, st string, pri float64) *Call {
	return &Call{
		Destination: dst,
		Duration:    dur,
		StartDate:   sd,
		StartTime:   st,
		Price:       pri,
	}
}

func (db *DB) CreateCall(call *Call) error {
	statement, _ := db.Prepare(`INSERT INTO call 
        (destionation, duration, start_date, start_time, price) 
        VALUES (?, ?, ?, ?, ?);`)

	if _, err := statement.Exec(call.Destination, call.Duration, call.StartDate, call.StartTime, call.Price); err != nil {
		return err
	}

	return nil
}
