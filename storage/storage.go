package storage

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Datastorer interface {
	RecordStorer
	BillStorer
	CallStorer
	InitSchema()
}

type DB struct {
	*sql.DB
}

func NewDB(dialect, dbName string) (*DB, error) {
	db, err := sql.Open(dialect, dbName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(1)

	return &DB{db}, nil
}

func (db *DB) billSchema() error {
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS bill(
        id INTEGER PRIMARY KEY NOT NULL,
        subscriber VARCHAR,
        month VARCHAR,
        year VARCHAR,
        price FLOAR)
    `)
	if _, err := statement.Exec(); err != nil {
		return err
	}

	return nil
}

func (db *DB) callSchema() error {
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS call (
        id INTEGER PRIMARY KEY NOT NULL,
        bill_id INTEGER NOT NULL,
        destionation VARCHAR, 
        duration VARCHAR,
        start_date INTEGER,
        start_time VARCHAR,
        price FLOAT,
        FOREIGN KEY(bill_id) REFERENCES bill(id))
    `)

	if _, err := statement.Exec(); err != nil {
		return err
	}

	return nil
}

func (db *DB) recordSchema() error {
	statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS record (
        id INTEGER PRIMARY KEY NOT NULL, 
        type VARCHAR, 
        timestamp VARCHAR,
        call_id INTEGER,
        source VARCHAR,
        destination VARCHAR)
    `)
	if _, err := statement.Exec(); err != nil {
		return err
	}

	return nil
}

func (db *DB) InitSchema() {
	db.billSchema()
	db.callSchema()
	db.recordSchema()
}
