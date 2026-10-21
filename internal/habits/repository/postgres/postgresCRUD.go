package postgres

import (
	"context"
	"errors"
	"log/slog"

	coreErrors "github.com/N1ktarchik/habbit-graph/internal/core/errors"
	"github.com/jackc/pgx/v5"

	"github.com/N1ktarchik/habbit-graph/internal/core/domain"
)

func (r *repository) CreateHabit(ctx context.Context, habit domain.Habit) (*domain.Habit, error) {
	query := `
		INSERT INTO habits 
			(title,description) 
		VALUES ($1,$2) 
		RETURNING id,title,description,created_at`

	model := &modelHabit{}
	err := r.db.QueryRow(ctx, query, habit.Title, habit.Description).Scan(
		&model.ID,
		&model.Title,
		&model.Description,
		&model.CreatedAt)

	if err != nil {

		r.log.Error("error in repository : postgres error : error to create habit", slog.Any("err:", err))
		return nil, err
	}

	result := model.toDomain()

	return &result, nil
}

func (r *repository) GetHabits(ctx context.Context) ([]domain.Habit, error) {
	query := `SELECT id,title,description,created_at FROM habits`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.log.Error("error in repository : postgres error : error to get habits -> query error", slog.Any("err:", err))
		return nil, err

	}
	defer rows.Close()

	arr := make([]domain.Habit, 0, 10)

	for rows.Next() {
		model := &modelHabit{}

		if err := rows.Scan(&model.ID, &model.Title, &model.Description, &model.CreatedAt); err != nil {
			r.log.Error("error in repository : postgres error : error to get habits -> scan error", slog.Any("err:", err))
			return nil, err
		}

		arr = append(arr, model.toDomain())

	}

	if err := rows.Err(); err != nil {
		r.log.Error("error in repository : postgres error : error to get habits -> rows error", slog.Any("err:", err))
		return nil, err
	}

	return arr, nil
}

func (r *repository) UpdateHabit(ctx context.Context, id string, title, desc *string) (*domain.Habit, error) {
	query := `
		UPDATE habits 
		SET 
			title = COALESCE($1, title),
			description = COALESCE($2, description)
		WHERE id = $3
		RETURNING id, title, description, created_at
	`

	model := &modelHabit{}

	err := r.db.QueryRow(ctx, query, title, desc, id).Scan(&model.ID, &model.Title, &model.Description, &model.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Warn("repository: postgres: habit not found for update", slog.String("id", id))
			return nil, coreErrors.NotFound()
		}

		r.log.Error("error in repository : postgres error : error to update habit", slog.Any("err:", err))
		return nil, err
	}

	result := model.toDomain()
	return &result, nil

}

func (r *repository) DeleteHabit(ctx context.Context, id string) error {
	query := `DELETE FROM habits WHERE id=$1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.log.Error("error in repository : postgres error : error to delete habit", slog.Any("err:", err))
		return err
	}

	if tag.RowsAffected() == 0 {
		r.log.Warn("repository: postgres: habit not found for deletion", slog.String("id", id))
		return coreErrors.NotFound()
	}

	return nil
}
