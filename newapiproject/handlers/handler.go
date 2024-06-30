package handlers

import (
	"newapiprojet/database"
)

type Handler struct {
	db database.Database
}

func NewHandler(db database.Database) *Handler {
	return &Handler{db: db}
}
