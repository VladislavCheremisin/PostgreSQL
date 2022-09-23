package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type (
	FirstName string
	Email     string
)
type EmailSearchHint struct {
	FirstName FirstName
	Email     Email
}

func searchName(ctx context.Context, dbpool *pgxpool.Pool, prefix string, limit int) ([]EmailSearchHint, error) {
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
	rows, err := dbpool.Query(ctx, sql, pattern, limit)
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
func main() {
	ctx := context.Background()
	url := "postgres://postgres:P@ssw0rd@localhost:5432/gopher_library"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}
	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()
	limit := 2
	hints, err := searchName(ctx, dbpool, "alex", limit)
	if err != nil {
		log.Fatal(err)
	}
	for _, hint := range hints {
		fmt.Println(hint.Email, hint.FirstName)
	}
}
