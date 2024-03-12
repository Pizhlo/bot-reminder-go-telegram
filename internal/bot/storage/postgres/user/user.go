package user

import (
	"database/sql"
	"fmt"
)

type UserRepo struct {
	db *sql.DB
}

func New(dbURl string) (*UserRepo, error) {
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}
	return &UserRepo{db}, nil
}
