package handlers

import (
	"log"
	"work-at-olist/storage"
)

type ErrorResponse struct {
	Status int                    `json:"-"`
	Errors map[string]interface{} `json:"errors"`
}

type Handler struct {
	DB     storage.Datastorer
	Logger *log.Logger
}

func New(db *storage.DB, logger *log.Logger) *Handler {
	return &Handler{db, logger}
}
