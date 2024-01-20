package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

type UserRepo struct {
	db *sql.DB
}

func New(dbURl string) (*UserRepo, error) {
	db, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}
	return &UserRepo{db}, nil
}

func (db *UserRepo) Save(ctx context.Context, id int64, tz *user.User) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("unable to create transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into users.users (tg_id) values($1)`, id)
	if err != nil {
		return fmt.Errorf("error while saving user in DB: %w", err)
	}

	return tx.Commit()
}

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

		err = rows.Scan(&u.ID, &u.TGID)
		if err != nil {
			return nil, fmt.Errorf("error while scanning user: %w", err)
		}

		res = append(res, u)
	}

	return res, nil
}

func (db *UserRepo) Update(ctx context.Context, id int64, updFun func(*user.User) (*user.User, error)) (*user.User, error) {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create transaction: %w", err)
	}

	oldUser, err := db.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error while getting user: %w", err)
	}

	newUser, err := updFun(oldUser)
	if err != nil {
		return nil, fmt.Errorf("err in update function: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into users.timezones (user_id, timezone, lon, lat) values($1, $2, $3, $4)`,
		id, newUser.Timezone.Name, newUser.Timezone.Lon, newUser.Timezone.Lat)
	if err != nil {
		return nil, fmt.Errorf("err while updating user timezone: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while committing: %w", err)
	}

	return newUser, nil
}

func (db *UserRepo) FindByTelegramID(ctx context.Context, tgid int64) (*user.User, error) {
	var id int
	var lon, lat float64
	var timezone string

	row := db.db.QueryRowContext(ctx, `select users.users.id, users.timezones.timezone, users.timezones.lon, users.timezones.lat
	from users.users join users.timezones on users.timezones.user_id = users.users.id where users.users.tg_id = $1`, tgid)

	err := row.Scan(&id, &timezone, &lon, &lat)
	if err != nil {
		return nil, fmt.Errorf("error while getting user by telegram ID: %w", err)
	}

	u := &user.User{
		ID:   id,
		TGID: tgid,
		Timezone: user.Timezone{
			Name: timezone,
			Lon:  lon,
			Lat:  lat,
		},
	}

	return u, nil
}

// func (db *UserRepo) GetUser(ctx context.Context, tgID int64) (model.User, error) {
// 	var user model.User

// 	err := db.QueryRow(ctx, `select id from "users" where tg_id = $1`, tgID).Scan(&user.ID)
// 	if err != nil {
// 		return user, err
// 	}

// 	if user.ID == 0 {
// 		return user, api_errs.ErrUserNotFound
// 	}

// 	return user, nil
// }

// func (db *UserRepo) SaveUser(ctx context.Context, telegramID int64) (int, error) {
// 	var id int
// 	err := db.QueryRow(context.Background(), `insert into "users"(tg_id) values($1) returning id`, telegramID).Scan(&id)
// 	if err != nil {
// 		return 0, errors.Wrap(err, `unable to save user`)
// 	}

// 	return id, nil
// }
