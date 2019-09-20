package storage

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type CallStorer interface {
	//GetCall(string) (*Call, error)
	GetAllCalls() *gorm.DB
	CreateCall(*Call) error
	GetCallsByCallId(int) ([]Call, error)
	//DeleteCall(*Call) error
}

type Call struct {
	Id          int    `json:"id"`
	Type        string `json:"type"`
	Timestamp   string `json:"timestamp"`
	CallId      int    `json:"call_id"`
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type ValidationMessages map[string]interface{}

func NewCall(ty string, ti string, ci int, sr string, ds string) *Call {
	return &Call{
		Type:        ty,
		Timestamp:   ti,
		CallId:      ci,
		Source:      sr,
		Destination: ds,
	}
}

func (c *Call) IsValid() (bool, map[string]interface{}) {
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
func (db *DB) GetCallsByCallId(callId int) ([]Call, error) {
	var calls []Call
	err := db.Find(&calls, "call_id = ?", callId).Error
	return calls, err
}

func (db *DB) GetAllCalls() *gorm.DB {
	return db.Find(&Call{})
}

func (db *DB) CreateCall(call *Call) (err error) {
	err = db.Create(&call).Error
	return
}
