package timezone

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/pkg/errors"
)

func (db *TimezoneRepo) Get(ctx context.Context, userID int64) (*user.Timezone, error) {
	return &user.Timezone{}, nil
}

func (db *TimezoneRepo) GetAll(ctx context.Context) ([]*user.User, error) {
	res := make([]*user.User, 0)

	rows, err := db.db.QueryContext(ctx, `select users.users.tg_id, users.timezones.timezone from users.timezones join users.users on users.users.id = users.timezones.user_id`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error while getting all users from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := &user.User{}

		err = rows.Scan(&u.TGID, &u.Timezone.Name)
		if err != nil {
			return nil, fmt.Errorf("error while scanning user: %w", err)
		}

		res = append(res, u)
	}

	return res, nil
}
