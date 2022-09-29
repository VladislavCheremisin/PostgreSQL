//go:build integration
// +build integration

// файлы с интеграционными тестами используют package storage_test,
// поэтому нужно явно ссылаться на storage, хотя все файлы лежат вместе.
package storage_test

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
	"os"
	"pkg/emailhint/storage"
	"testing"
)

func TestIntegrationSearch(t *testing.T) {
	ctx := context.Background()
	dbpool := connect(ctx)
	defer dbpool.Close()
	tests := []struct {
		name    string
		store   *storage.PG
		ctx     context.Context
		prefix  string
		limit   int
		prepare func(*pgxpool.Pool)
		check   func(*testing.T, []storage.EmailSearchHint, error)
	}{
		{
			name:   "success",
			store:  storage.NewPG(dbpool),
			ctx:    context.Background(),
			prefix: "alex",
			limit:  5,
			prepare: func(dbpool *pgxpool.Pool) {
				// Подготовка тестовых данных
				dbpool.Exec(context.Background(), `insert into users ...`)
			},
			check: func(t *testing.T, hints []storage.EmailSearchHint, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, hints)
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare(dbpool)
			hints, err := tt.store.Search(tt.ctx, tt.prefix, tt.limit)
			tt.check(t, hints, err)
		})
	}
}

// Соединение с экземпляром Postgres
func connect(ctx context.Context) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	return dbpool
}
