package storage

import (
	"strconv"
)

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

type ValidationMessages map[string]interface{}

func NewRecord(ty string, ti string, ci int, sr string, ds string) *Record {
	return &Record{
		Type:        ty,
		Timestamp:   ti,
		CallId:      ci,
		Source:      sr,
		Destination: ds,
	}
}

func (c *Record) IsValid() (bool, map[string]interface{}) {
	var errs = ValidationMessages{}
	var valid = true

	if c.CallId <= 0 {
		errs["call_id"] = []string{"invalid [call_id] field value"}
		valid = false
	}

	if c.Type == "" {
		errs["type"] = []string{"[type] field can't blank"}
		valid = false
	}

	if c.Timestamp == "" {
		errs["timestamp"] = []string{"[timestamp] field can't blank"}
		valid = false
	}

	if c.Type == "start" && (c.Source == "" || !validPhone(c.Source)) {
		errs["source"] = []string{"invalid [source] field value"}
		valid = false
	}

	if c.Type == "start" && (c.Destination == "" || !validPhone(c.Destination)) {
		errs["destination"] = []string{"invalid [destination] field value"}
		valid = false
	}

	return valid, errs
}

func validPhone(ph string) bool {
	if _, err := strconv.Atoi(ph); err != nil {
		return false
	}
	if len(ph) > 11 || len(ph) < 10 {
		return false
	}
	return true
}

func (db *DB) GetRecordsByCallId(callId int) ([]Record, error) {
	var record Record
	var records []Record

	rows, err := db.Query(`SELECT * FROM record
        WHERE call_id = ?;`, callId)
	if err != nil {
		return nil, err
	}

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

	for rows.Next() {
		if err := rows.Scan(&record.Id, &record.Type, &record.Timestamp, &record.CallId, &record.Source, &record.Destination); err != nil {
			return nil, err
		}
		records = append(records, record)

	}
	return records, nil
}

func (db *DB) CreateRecord(call *Record) (err error) {
	statement, _ := db.Prepare(`INSERT INTO record (type, timestamp, call_id, source, destination) VALUES (?, ?, ?, ?, ?)`)

	if _, err := statement.Exec(call.Type, call.Timestamp, call.CallId, call.Source, call.Destination); err != nil {
		return err
	}

	return nil
}
