package storage

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

type RecordStorer interface {
	GetAllRecords() *gorm.DB
	CreateRecord(*Record) error
	GetRecordsByCallId(int) ([]Record, error)
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

// TODO Use method to validate call type together with call_id
func (db *DB) GetRecordsByCallId(callId int) ([]Record, error) {
	var calls []Record
	err := db.Find(&calls, "call_id = ?", callId).Error
	return calls, err
}

func (db *DB) GetAllRecords() *gorm.DB {
	return db.Find(&Record{})
}

func (db *DB) CreateRecord(call *Record) (err error) {
	err = db.Create(&call).Error
	return
}
