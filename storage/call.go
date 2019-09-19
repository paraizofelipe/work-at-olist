package storage

import (
	"github.com/jinzhu/gorm"
	"strconv"
)

type CallStorer interface {
	//GetCall(string) (*Call, error)
	GetAllCalls() *gorm.DB
	CreateCall(*Call) error
	//DeleteCall(*Call) error
}

type Call struct {
	ID          int
	Type        string
	Timestamp   string
	CallId      int
	Source      string
	Destination string
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

	if c.Type == "" {
		errs["type"] = []string{"type field can't blank"}
		valid = false
	}

	if c.Destination == "" {
		errs["destination"] = []string{"destination field can't blank"}
		valid = false
	}

	if c.Source == "" {
		errs["source"] = []string{"source field can't blank"}
		valid = false
	}

	if c.Timestamp == "" {
		errs["timestamp"] = []string{"source field can't blank"}
		valid = false
	}

	if c.Type == "start" && (c.Source == "" || c.Destination == "") {
		errs["start_call"] = []string{"call start cannot be null"}
		valid = false
	}

	if c.Type == "start" && (!validPhone(c.Source) || !validPhone(c.Destination)) {
		errs["start_call"] = []string{"call start cannot be null"}
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

func (db *DB) GetAllCalls() *gorm.DB {
	return db.Find(&Call{})
}

func (db *DB) CreateCall(call *Call) (err error) {
	err = db.Create(&call).Error
	return
}
