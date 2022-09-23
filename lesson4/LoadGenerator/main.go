package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

type AttackResults struct {
	Duration         time.Duration
	Threads          int
	QueriesPerformed uint64
}

type (
	FirstName string
	Email     string
	Title     string
	Author    string
)
type (
	NumberBooks int
)
type EmailSearchHint struct {
	FirstName FirstName
	Email     Email
}

type BooksSearchHint struct {
	Title  Title
	Author Author
}

type NumberSearchHint struct {
	NumberBooks NumberBooks
	Author      Author
}

func search(ctx context.Context, dbpool *pgxpool.Pool, prefix string, limit int) ([]EmailSearchHint, error) {
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

func searchAuthors(ctx context.Context, dbpool *pgxpool.Pool) ([]BooksSearchHint, error) {
	const sql = `
				SELECT title, author
				FROM books_authors, books, authors
				WHERE books_authors.book_id = books.id AND books_authors.author_id = authors.id;
`

	rows, err := dbpool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var hints []BooksSearchHint
	for rows.Next() {
		var hint BooksSearchHint
		err = rows.Scan(&hint.Author, &hint.Title)
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

func searchNumberBooks(ctx context.Context, dbpool *pgxpool.Pool) ([]NumberSearchHint, error) {
	const sql = `
				SELECT authors.author, COUNT(authors.author)
				FROM books_authors, books, authors
				WHERE books_authors.book_id = books.id AND books_authors.author_id = authors.id
				GROUP BY authors.author;
`

	rows, err := dbpool.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	defer rows.Close()

	var hints []NumberSearchHint
	for rows.Next() {
		var hint NumberSearchHint
		err = rows.Scan(&hint.Author, &hint.NumberBooks)
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

func attack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64

	attacker := func(stopAt time.Time) {
		for {
			//select request type
			_, err := searchNumberBooks(ctx, dbpool)
			//_, err := searchAuthors(ctx, dbpool)
			//_, err := search(ctx, dbpool, "alex", 5)
			if err != nil {
				log.Fatal(err)
			}
			atomic.AddUint64(&queries, 1)
			if time.Now().After(stopAt) {
				return
			}
		}
	}
	var wg sync.WaitGroup
	wg.Add(threads)
	startAt := time.Now()
	stopAt := startAt.Add(duration)
	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}
	wg.Wait()
	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}

func main() {
	ctx := context.Background()
	url := "postgres://postgres:P@ssw0rd@localhost:5432/gopher_library"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}
	cfg.MaxConns = 30
	cfg.MinConns = 30
	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()
	duration := time.Duration(10 * time.Second)
	threads := 1000
	fmt.Println("start attack")
	res := attack(ctx, duration, threads, dbpool)
	fmt.Println("duration:", res.Duration)
	fmt.Println("threads:", res.Threads)
	fmt.Println("queries:", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS:", qps)
}
