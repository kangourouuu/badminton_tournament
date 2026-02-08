package api

import "github.com/uptrace/bun"

type Handler struct {
	DB *bun.DB
}

func NewHandler(db *bun.DB) *Handler {
	return &Handler{DB: db}
}
