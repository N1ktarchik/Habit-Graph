package postgres

import (
	"context"
	"log/slog"
	"time"
)

func (r *repository) GetGraph(ctx context.Context, habitID string, year int) (map[string]int, error) {

	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, time.December, 31, 0, 0, 0, 0, time.UTC)

	query := `
		SELECT date, duration_minutes 
		FROM activities 
		WHERE habit_id = $1 AND date >= $2 AND date <= $3
	`

	rows, err := r.db.Query(ctx, query, habitID, startDate, endDate)
	if err != nil {
		r.log.Error("error in repository : postgres error : error to get graph -> query error", slog.Any("err:", err))
		return nil, err

	}
	defer rows.Close()

	graph := make(map[string]int, 10)

	for rows.Next() {
		var date time.Time
		var minutes int

		if err := rows.Scan(&date, &minutes); err != nil {
			r.log.Error("error in repository : postgres error : error to get graph -> scan error", slog.Any("err:", err))
			return nil, err
		}

		dateStr := date.Format(time.DateOnly)
		graph[dateStr] = minutes
	}

	if err := rows.Err(); err != nil {
		r.log.Error("error in repository : postgres error : error to get graph -> rows error", slog.Any("err:", err))
		return nil, err
	}

	return graph, nil
}

// func (r *repository) UpsertActivity(ctx context.Context, habitID string, userData domain.Track) error {

// }
