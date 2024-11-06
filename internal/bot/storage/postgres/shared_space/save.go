package sharedspace

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

func (db *sharedSpaceRepo) Save(ctx context.Context, space model.SharedSpace) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return err
	}

	var id int64
	row := tx.QueryRowContext(ctx, "insert into shared_spaces.shared_spaces (name, created, creator) values($1, $2, (select id from users.users where tg_id=$3)) returning id",
		space.Name, space.Created, space.Creator.TGID)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("error scanning id of new shared space: %+v", err)
	}

	if err := db.SaveParticipant(ctx, tx, id, space.Creator); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error saving participant while creating shared space: %+v", err)
	}

	return tx.Commit()
}

func (db *sharedSpaceRepo) SaveParticipant(ctx context.Context, tx *sql.Tx, spaceID int64, user user.User) error {
	// tx, err := db.tx(ctx)
	// if err != nil {
	// 	return err
	// }

	_, err := tx.ExecContext(ctx, "insert into shared_spaces.participants (space_id, user_id) values($1, (select id from users.users where tg_id = $2))",
		spaceID, user.TGID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "update users.users set username=$1 where tg_id = $2;", user.Username, user.TGID)
	if err != nil {
		return fmt.Errorf("error saving participant's username: %+v", err)
	}

	return nil
}
