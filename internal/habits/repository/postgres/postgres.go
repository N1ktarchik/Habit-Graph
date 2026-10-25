package postgres

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	log *slog.Logger
	db  *pgxpool.Pool
}

func NewRepository(log *slog.Logger, db *pgxpool.Pool) *repository {
	return &repository{
		log: log,
		db:  db,
	}
}
