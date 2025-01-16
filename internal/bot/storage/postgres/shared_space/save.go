package sharedspace

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

	if err := db.saveParticipant(ctx, tx, id, space.Creator); err != nil {
		_ = db.rollback()
		return fmt.Errorf("error saving participant while creating shared space: %+v", err)
	}

	return db.commit()
}

func (db *sharedSpaceRepo) SaveParticipant(ctx context.Context, spaceID int64, user model.Participant) error {
	var tx *sql.Tx
	var err error

	if db.currentTx == nil {
		tx, err = db.tx(ctx)
		if err != nil {
			return err
		}
	} else {
		tx = db.currentTx
	}

	err = db.saveParticipant(ctx, tx, spaceID, user)
	if err != nil {
		return err
	}

	return db.commit()
}

func (db *sharedSpaceRepo) saveParticipant(ctx context.Context, tx *sql.Tx, spaceID int64, user model.Participant) error {
	_, err := tx.ExecContext(ctx, "insert into shared_spaces.participants (space_id, user_id, state_id) values($1, (select id from users.users where tg_id = $2), (select id from shared_spaces.participants_states where state = $3))",
		spaceID, user.TGID, user.State)
	if err != nil {
		switch e := err.(type) {
		case *pq.Error: // пользователь уже существует в пространстве
			if e.Code == "23505" {
				return api_errors.ErrUserAlreadyExists
			}
		}

		return err
	}

	// _, err = tx.ExecContext(ctx, "update users.users set username=$1 where tg_id = $2;", user.Username, user.TGID)
	// if err != nil {
	// 	return fmt.Errorf("error saving participant's username: %+v", err)
	// }

	return nil
}

func (db *sharedSpaceRepo) SaveNote(ctx context.Context, note model.Note) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return err
	}

	var id uuid.UUID
	row := tx.QueryRowContext(ctx, "insert into shared_spaces.notes (user_id, text, created, space_id) values((select id from users.users where tg_id=$1), $2, $3, $4) returning id",
		note.Creator.TGID, note.Text, note.Created, note.Space.ID)

	err = row.Scan(&id)
	if err != nil {
		return fmt.Errorf("error scanning ID: %+v", err)
	}

	// создаем структуру для сохранения в elastic
	elasticData := elastic.Data{
		Index: elastic.NoteIndex,
		Model: &elastic.Note{
			ID:          id,
			Text:        note.Text,
			TgID:        note.Creator.TGID,
			SharedSpace: note.Space,
		}}

	// сохраняем в elastic
	err = db.elasticClient.Save(ctx, elasticData)
	if err != nil {
		// отменяем транзакцию в случае ошибки (для консистентности данных)
		_ = db.rollback()
		return fmt.Errorf("error saving note for shared space to Elastic: %+v", err)
	}

	return db.commit()
}

func (db *sharedSpaceRepo) ProcessInvitation(ctx context.Context, from, to model.Participant, spaceID int64) error {
	_, err := db.tx(ctx)
	if err != nil {
		return err
	}

	err = db.saveInvitation(ctx, from.TGID, to.TGID, spaceID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return api_errors.ErrInvitationExists
		}

		_ = db.rollback()
		return err
	}

	err = db.saveParticipant(ctx, db.currentTx, spaceID, to)
	if err != nil {
		_ = db.rollback()
		return err
	}

	return db.commit()
}

func (db *sharedSpaceRepo) saveInvitation(ctx context.Context, from, to, spaceID int64) error {
	_, err := db.currentTx.ExecContext(ctx, `insert into shared_spaces.invitations ("from", "to", space_id) values((select id from users.users where tg_id = $1), (select id from users.users where tg_id = $2), $3)`,
		from, to, spaceID)
	return err
}

func (db *sharedSpaceRepo) SetParticipantState(ctx context.Context, user model.Participant, state string, spaceID int64) error {
	_, err := db.tx(ctx)
	if err != nil {
		return err
	}

	_, err = db.currentTx.ExecContext(ctx,
		`insert into shared_spaces.participants (space_id, user_id, state_id) values 
	($1, (select id from users.users where tg_id = $2), 
	(select id from shared_spaces.participants_states where state = $3))
	on conflict (space_id, user_id) do update
	set state_id=(select id from shared_spaces.participants_states where state = $3);`,
		spaceID, user.TGID, state)
	if err != nil {
		_ = db.rollback()
		return err
	}

	return db.commit()
}
