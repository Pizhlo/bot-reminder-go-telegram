package timezone

import (
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
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

func (db *TimezoneRepo) Close() {
	if err := db.db.Close(); err != nil {
		logrus.Errorf("error on closing tz repo: %v", err)
	}
}
