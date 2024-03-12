package user

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

func (db *UserRepo) Save(ctx context.Context, id int64, tz *user.User) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("unable to create transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into users.users (tg_id) values($1) on conflict (tg_id) do nothing`, id)
	if err != nil {
		return fmt.Errorf("error while saving user in DB: %w", err)
	}

	return tx.Commit()
}
