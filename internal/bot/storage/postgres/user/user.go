package user

import (
	"context"

	api_errs "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type UserRepo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *UserRepo {
	return &UserRepo{conn}
}

func (db *UserRepo) GetUser(ctx context.Context, tgID int64) (model.User, error) {
	var user model.User

	err := db.QueryRow(ctx, `select id from "users" where tg_id = $1`, tgID).Scan(&user.ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, api_errs.ErrUserNotFound
	}

	return user, nil
}

func (db *UserRepo) SaveUser(ctx context.Context, telegramID int64) (int, error) {
	var id int
	err := db.QueryRow(context.Background(), `insert into "users"(tg_id) values($1) returning id`, telegramID).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, `unable to save user`)
	}

	return id, nil
}
