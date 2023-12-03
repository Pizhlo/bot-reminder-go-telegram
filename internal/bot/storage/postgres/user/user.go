package user

import (
	"context"
	"database/sql"
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

func (db *UserRepo) Add(ctx context.Context, tgid int) (*user.User, error) {
	checkUser, err := db.FindByTelegramID(ctx, tgid)
	if err != nil {
		return nil, fmt.Errorf("error while getting user: %w", err)
	}

	if checkUser.ID != 0 {
		u := &user.User{
			ID:   checkUser.ID,
			TGID: tgid,
		}

		return u, nil
	}

	var id int
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create transaction: %w", err)
	}

	row := tx.QueryRowContext(ctx, `insert into users.users (tg_id) values($1) returning id`, tgid)

	err = row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error while scanning id: %w", err)
	}

	u := &user.User{
		ID:   id,
		TGID: tgid,
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while committing: %w", err)
	}

	return u, nil
}

func (db *UserRepo) Get(ctx context.Context, id int) (*user.User, error) {
	var dbID int
	row := db.db.QueryRowContext(ctx, `select id from users.users where id = $1`, id)
	err := row.Scan(&dbID)
	if err != nil {
		return nil, fmt.Errorf("error while scanning user by id %d: %w", id, err)
	}

	u := &user.User{
		ID:   dbID,
		TGID: id,
	}

	return u, nil
}

func (db *UserRepo) Update(ctx context.Context, id int, updFun func(*user.User) (*user.User, error)) (*user.User, error) {
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

	_, err = tx.ExecContext(ctx, `insert into users.timezones (user_id, timezone) values($1, $2)`, id, newUser.Timezone.Name)
	if err != nil {
		return nil, fmt.Errorf("err while updating user timezone: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("error while committing: %w", err)
	}

	return newUser, nil
}

func (db *UserRepo) FindByTelegramID(ctx context.Context, tgid int) (*user.User, error) {
	var id int
	row := db.db.QueryRowContext(ctx, `select id from users.users where tg_id = $1`, tgid)

	err := row.Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("error while getting user by telegram ID: %w", err)
	}

	u := &user.User{
		ID:   id,
		TGID: tgid,
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
