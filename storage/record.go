package storage

type RecordStorer interface {
	CreateRecord(*Record) error
	GetRecordsByCallId(int) ([]Record, error)
	GetRecordsByType(int, string) ([]Record, error)
}

type Record struct {
	Id          int    `json:"id"`
	Type        string `json:"type"`
	Timestamp   string `json:"timestamp"`
	CallId      int    `json:"call_id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func NewRecord(ty string, ti string, ci int, sr string, ds string) *Record {
	return &Record{
		Type:        ty,
		Timestamp:   ti,
		CallId:      ci,
		Source:      sr,
		Destination: ds,
	}
}

func (db *DB) GetRecordsByCallId(callId int) ([]Record, error) {
	var record Record
	var records []Record

	rows, err := db.Query(`SELECT * FROM record
        WHERE call_id = ?;`, callId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&record.Id, &record.Type, &record.Timestamp, &record.CallId, &record.Source, &record.Destination); err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	return records, nil
}

func (db *DB) GetRecordsByType(callId int, callType string) ([]Record, error) {
	var record Record
	var records []Record

	rows, err := db.Query(`SELECT * FROM record 
        WHERE call_id = ?
        AND type = ?;`, callId, callType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&record.Id, &record.Type, &record.Timestamp, &record.CallId, &record.Source, &record.Destination); err != nil {
			return nil, err
		}
		records = append(records, record)

	}
	return records, nil
}

func (db *DB) CreateRecord(call *Record) (err error) {
	statement, err := db.Prepare(`INSERT INTO record (type, timestamp, call_id, source, destination) VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer statement.Close()

	if _, err := statement.Exec(call.Type, call.Timestamp, call.CallId, call.Source, call.Destination); err != nil {
		return err
	}
	return nil
}
