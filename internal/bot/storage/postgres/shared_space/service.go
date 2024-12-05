package sharedspace

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type sharedSpaceRepo struct {
	db            *sql.DB
	elasticClient elasticClient
	currentTx     *sql.Tx
}

type elasticClient interface {
	Save(ctx context.Context, search elastic.Data) error
	// SearchByText производит поиск по тексту (названию). Возвращает ID из базы подходящих записей
	SearchByText(ctx context.Context, search elastic.Data) ([]uuid.UUID, error)
	// SearchByID производит поиск по ID из базы. Возвращает ID  из эластика подходящих записей
	SearchByID(ctx context.Context, search elastic.Data) ([]string, error)
	Delete(ctx context.Context, search elastic.Data) error
	DeleteAllByUserID(ctx context.Context, data elastic.Data) error
	Update(ctx context.Context, search elastic.Data) error
}

func New(dbURl string, elasticClient elasticClient) (*sharedSpaceRepo, error) {
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}
	return &sharedSpaceRepo{db, elasticClient, nil}, nil
}

func (db *sharedSpaceRepo) tx(ctx context.Context) (*sql.Tx, error) {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, err
	}

	db.currentTx = tx

	return tx, nil
}

func (db *sharedSpaceRepo) Close() {
	if err := db.db.Close(); err != nil {
		logrus.Errorf("error on closing shared space repo: %v", err)
	}
}

func (db *sharedSpaceRepo) commit() error {
	tx := db.currentTx
	db.currentTx = nil
	return tx.Commit()
}

func (db *sharedSpaceRepo) rollback() error {
	tx := db.currentTx
	db.currentTx = nil
	return tx.Rollback()
}
