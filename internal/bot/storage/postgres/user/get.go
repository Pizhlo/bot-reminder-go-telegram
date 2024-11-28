package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (db *UserRepo) GetByID(ctx context.Context, tgID int64) (*model.User, error) {
	var dbID int

	row := db.db.QueryRowContext(ctx, `select id from users.users where tg_id = $1`, tgID)
	err := row.Scan(&dbID)
	if err != nil {
		return nil, fmt.Errorf("error while scanning user by id %d: %w", tgID, err)
	}

	u := &model.User{
		ID:   dbID,
		TGID: tgID,
	}

	return u, nil
}

func (db *UserRepo) GetAll(ctx context.Context) ([]*model.User, error) {
	res := make([]*model.User, 0)

	rows, err := db.db.QueryContext(ctx, `select * from users.users`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error while getting all users from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := &model.User{}

		err = rows.Scan(&u.ID, &u.TGID, &u.UsernameSQL)
		if err != nil {
			return nil, fmt.Errorf("error while scanning user: %w", err)
		}

		res = append(res, u)
	}

	return res, nil
}

func (db *UserRepo) GetState(ctx context.Context, id int64) (string, error) {
	state := ""

	row := db.db.QueryRowContext(ctx, `select name from users.state_types 
	where id = (select state_id from users.states where user_id = (select id from users.users where tg_id = $1)); `, id)

	err := row.Scan(&state)
	if err != nil {
		return "", fmt.Errorf("error while scanning state: %w", err)
	}

	return state, nil
}

func (db *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	row := db.db.QueryRowContext(ctx, `select * from users.users where username = $1`, username)

	err := row.Scan(&user.ID, &user.TGID, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, api_errors.ErrUserNotFound
		}

		return nil, fmt.Errorf("error while scanning user by username %s: %w", username, err)
	}

	return &user, nil
}
