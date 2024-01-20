package timezone

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/pkg/errors"
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

// SaveUserTimezone сохраняет часовой пояс пользователя. Аргументы: id - id базы данных, timezone - модель часового пояса
func (db *TimezoneRepo) Save(ctx context.Context, id int64, timezone *user.Timezone) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return errors.Wrap(err, `error while creating transaction`)
	}

	_, err = tx.ExecContext(ctx, `insert into users.timezones(user_id, timezone, lon, lat) values((select id from users.users where tg_id=$1), $2, $3, $4)`, id, timezone.Name, timezone.Lon, timezone.Lat)
	if err != nil {
		return errors.Wrap(err, `error while saving timezone`)
	}

	return tx.Commit()
}

func (db *TimezoneRepo) Get(ctx context.Context, userID int64) (*user.Timezone, error) {
	return &user.Timezone{}, nil
}

func (db *TimezoneRepo) GetAll(ctx context.Context) ([]*user.User, error) {
	res := make([]*user.User, 0)

	rows, err := db.db.QueryContext(ctx, `select * from users.timezones`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("error while getting all users from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := &user.User{}

		err = rows.Scan(&u.ID, &u.TGID, &u.Timezone.Name, &u.Timezone.Lon, &u.Timezone.Lat)
		if err != nil {
			return nil, fmt.Errorf("error while scanning user: %w", err)
		}

		res = append(res, u)
	}

	return res, nil
}
