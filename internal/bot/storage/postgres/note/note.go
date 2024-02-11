package note

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type NoteRepo struct {
	db *sql.DB
}

func New(dbURl string) (*NoteRepo, error) {
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}
	return &NoteRepo{db}, nil
}
