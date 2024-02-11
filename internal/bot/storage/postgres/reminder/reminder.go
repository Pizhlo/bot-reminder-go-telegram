package reminder

import (
	"database/sql"
	"fmt"
)

type ReminderRepo struct {
	db *sql.DB
}

func New(dbURl string) (*ReminderRepo, error) {
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}
	return &ReminderRepo{db}, nil
}
