package sharedspace

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (db *sharedSpaceRepo) DeleteInvitation(ctx context.Context, from, to model.Participant, spaceID int64) error {
	_, err := db.tx(ctx)
	if err != nil {
		return err
	}

	_, err = db.db.ExecContext(ctx, `delete from shared_spaces.invitations where "from" = (select id from users.users where tg_id = $1) 
	and "to" = (select id from users.users where tg_id = $2)
	and space_id = $3`,
		from.TGID, to.TGID, spaceID)
	if err != nil {
		_ = db.rollback()
		return err
	}

	return nil
}
