package sharedspace

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (db *sharedSpaceRepo) Save(ctx context.Context, space model.SharedSpace) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "insert into shared_spaces.shared_spaces (name, created, creator) values($1, $2, (select id from users.users where tg_id=$3))",
		space.Name, space.Created, space.Creator.TGID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
