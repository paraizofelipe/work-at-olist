package storage

import (
	"github.com/jinzhu/gorm"
)

type Datastorer interface {
	RecordStorer
	BillStore
	InitSchema()
}

type DB struct {
	*gorm.DB
}

func NewDB(dialect, dbName string) (*DB, error) {
	db, err := gorm.Open(dialect, dbName)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) InitSchema() {
	db.AutoMigrate(&Record{})
	db.AutoMigrate(&Bill{})
}
