package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
)

func (db *UserRepo) Save(ctx context.Context, id int64, tz *model.User) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("unable to create transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into users.users (tg_id, username) values($1, $2) on conflict (tg_id) do nothing`, id, tz.Username)
	if err != nil {
		return fmt.Errorf("error while saving user in DB: %w", err)
	}

	return tx.Commit()
}

func (db *UserRepo) SaveState(ctx context.Context, id int64, state string) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("unable to create transaction: %w", err)
	}

	logrus.Debugf("saving state: %s", state)

	_, err = tx.ExecContext(ctx, `insert into users.states (user_id, state_id) values((select id from users.users where tg_id = $1),
	(select id from users.state_types where name = $2))
	on conflict (user_id) do update set state_id=(select id from users.state_types where name = $2);`, id, state)
	if err != nil {
		return fmt.Errorf("error while saving user's state in DB: %w", err)
	}

	return tx.Commit()
}
