package handlers

import (
	"log"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"work-at-olist/storage"
)

type Handler struct {
	DB     storage.Datastorer
	Logger *log.Logger
}

func New(db *storage.DB, logger *log.Logger) *Handler {
	return &Handler{db, logger}
}
