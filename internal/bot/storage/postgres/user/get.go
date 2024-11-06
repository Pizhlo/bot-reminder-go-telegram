package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

func (db *UserRepo) Get(ctx context.Context, tgID int64) (*user.User, error) {
	var dbID int

	row := db.db.QueryRowContext(ctx, `select id from users.users where tg_id = $1`, tgID)
	err := row.Scan(&dbID)
	if err != nil {
		return nil, fmt.Errorf("error while scanning user by id %d: %w", tgID, err)
	}

	u := &user.User{
		ID:   dbID,
		TGID: tgID,
	}

	return u, nil
}

func (db *UserRepo) GetAll(ctx context.Context) ([]*user.User, error) {
	res := make([]*user.User, 0)

	rows, err := db.db.QueryContext(ctx, `select * from users.users`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error while getting all users from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := &user.User{}

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
