package timezone

import (
	"database/sql"
	"fmt"
)

type TimezoneRepo struct {
	db *sql.DB
}

func New(dbURl string) (*TimezoneRepo, error) {
	conn, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}

	return &TimezoneRepo{conn}, nil
}
