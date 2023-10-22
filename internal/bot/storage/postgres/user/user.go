package user

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type UserRepo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{conn}
}

func (db *UserRepo) SaveUser(telegramID int64) (int, error) {
	var id int
	err := db.QueryRow(context.Background(), `insert into "users"(tg_id) values(?) returning id`, telegramID).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, `unable to save user`)
	}

	return id, nil
}
