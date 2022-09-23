package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type (
	UserID      int
	PositionID  int
	CharacterID int
)

type (
	Email string
)

type User struct {
	FirstName    string
	LastName     string
	Email        Email
	LoadingBooks PositionID
	CharacterID  CharacterID
}

func insert(ctx context.Context, dbpool *pgxpool.Pool, user User) (UserID, error) {
	const sql = `
				insert into users (firstname, lastname, email, loading_books, character_id) 
				values ($1, $2, $3, $4, $5)
				returning id;
`

	var id UserID
	err := dbpool.QueryRow(ctx, sql,
		// Параметры должны передаваться в том порядке,
		// в котором перечислены столбцы в SQL запросе.
		user.FirstName,
		user.LastName,
		user.Email,
		user.LoadingBooks,
		user.CharacterID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert employee: %w", err)
	}
	return id, nil
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
	employee := User{
		FirstName:    "Alexander4",
		LastName:     "Alexandrov4",
		Email:        "alexandr1234@empire.ru",
		LoadingBooks: 6,
		CharacterID:  2,
	}
	id, err := insert(ctx, dbpool, employee)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(id)
}
