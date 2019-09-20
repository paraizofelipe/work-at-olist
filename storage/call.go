package storage

type Call struct {
	Id          string  `json:"id"`
	Destination string  `json:"destination"`
	Duration    string  `json:"duration"`
	StartDate   int     `json:"start_date"`
	StartTime   string  `json:"start_time"`
	Price       float64 `json:"price"`
}
