package sharedspace

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sirupsen/logrus"
)

type sharedSpaceRepo struct {
	db *sql.DB
}

func New(dbURl string) (*sharedSpaceRepo, error) {
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}
	return &sharedSpaceRepo{db}, nil
}

func (db *sharedSpaceRepo) tx(ctx context.Context) (*sql.Tx, error) {
	return db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
}

func (db *sharedSpaceRepo) Close() {
	if err := db.db.Close(); err != nil {
		logrus.Errorf("error on closing shared space repo: %v", err)
	}
}
