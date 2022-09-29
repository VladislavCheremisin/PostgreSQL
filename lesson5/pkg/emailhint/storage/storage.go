package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type (
	FirstName string
	Email     string
)
type EmailSearchHint struct {
	FirstName FirstName
	Email     Email
}

type PG struct {
	dbpool *pgxpool.Pool
}

func NewPG(dbpool *pgxpool.Pool) *PG {
	return &PG{dbpool}
}

func (s *PG) Search(ctx context.Context, prefix string, limit int) ([]EmailSearchHint, error) {
	const sql = `
				select
				firstname,
				email
				from users
				where email like $1
				order by email asc
				limit $2;
`
	pattern := prefix + "%"
	rows, err := s.dbpool.Query(ctx, sql, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var hints []EmailSearchHint

	for rows.Next() {
		var hint EmailSearchHint
		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.Email, &hint.FirstName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		hints = append(hints, hint)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}
	return hints, nil
}
