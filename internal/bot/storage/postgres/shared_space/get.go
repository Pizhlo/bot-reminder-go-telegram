package sharedspace

import (
	"context"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (db *sharedSpaceRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.SharedSpace, error) {
	spaces := make([]model.SharedSpace, 0)

	rows, err := db.db.QueryContext(ctx,
		`select shared_spaces.shared_spaces.id, row_number()over(order by created) AS space_number, name, created, users.users.tg_id, users.users.username from shared_spaces.participants
join shared_spaces.shared_spaces on shared_spaces.shared_spaces.id = shared_spaces.participants.space_id
join users.users on users.users.id = shared_spaces.shared_spaces.creator
where shared_spaces.participants.user_id in (select id from users.users where tg_id = $1);`, userID)

	if err != nil {
		return nil, fmt.Errorf("error while getting all shared spaces from DB by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		space := model.SharedSpace{}

		err := rows.Scan(&space.ID, &space.ViewID, &space.Name, &space.Created, &space.Creator.TGID, &space.Creator.Username)
		if err != nil {
			return nil, fmt.Errorf("error scanning shared space: %+v", err)
		}

		participants, err := db.getAllParticipants(ctx, space.ID)
		if err != nil {
			return nil, err
		}

		notes, err := db.getAllNotes(ctx, space.ID)
		if err != nil {
			return nil, err
		}

		reminders, err := db.getAllReminders(ctx, space.ID)
		if err != nil {
			return nil, err
		}

		space.Participants = participants
		space.Notes = notes
		space.Reminders = reminders

		spaces = append(spaces, space)
	}

	if len(spaces) == 0 {
		return nil, api_errors.ErrSharedSpacesNotFound
	}

	return spaces, nil
}

func (db *sharedSpaceRepo) getAllParticipants(ctx context.Context, spaceID int) ([]model.Participant, error) {
	res := make([]model.Participant, 0)

	rows, err := db.db.QueryContext(ctx, `select tg_id, username, state from shared_spaces.participants 
join users.users on users.users.id = shared_spaces.participants.user_id
join shared_spaces.participants_states on shared_spaces.participants_states.id = shared_spaces.participants.state_id
where space_id = $1;`, spaceID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all participants for shared space from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := model.Participant{}

		err := rows.Scan(&user.TGID, &user.UsernameSQL, &user.State)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID while searching all space's participants: %+v", err)
		}

		if user.UsernameSQL.Valid {
			user.Username = user.UsernameSQL.String
		}

		res = append(res, user)
	}

	return res, nil
}

func (db *sharedSpaceRepo) getAllNotes(ctx context.Context, spaceID int) ([]model.Note, error) {
	res := make([]model.Note, 0)

	rows, err := db.db.QueryContext(ctx, `select note_number, text, created, last_edit, username, tg_id from shared_spaces.notes
	join shared_spaces.notes_view on shared_spaces.notes_view.id = shared_spaces.notes.id
	join users.users on users.users.id = shared_spaces.notes.user_id
	where shared_spaces.notes.space_id = $1 
	order by created ASC;`, spaceID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all notes for shared space from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ViewID, &note.Text, &note.Created, &note.LastEditSql, &note.Creator.Username, &note.Creator.TGID)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID while searching all space's notes: %+v", err)
		}

		res = append(res, note)
	}

	return res, nil
}

func (db *sharedSpaceRepo) getAllReminders(ctx context.Context, spaceID int) ([]model.Reminder, error) {
	res := make([]model.Reminder, 0)

	rows, err := db.db.QueryContext(ctx, `select reminder_number, text, created, username, tg_id from shared_spaces.reminders
	join shared_spaces.reminders_view on shared_spaces.reminders_view.id = shared_spaces.reminders.id
	join users.users on users.users.id = shared_spaces.reminders.user_id
	where shared_spaces.reminders.space_id = $1 
	order by created ASC;`, spaceID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all reminders for shared space from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		reminder := model.Reminder{}

		err := rows.Scan(&reminder.ID, &reminder.ViewID, &reminder.TgID, &reminder.Name, &reminder.Created, &reminder.Date, &reminder.Time, &reminder.Type, &reminder.Job.ID)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID while searching all space's reminders: %+v", err)
		}

		res = append(res, reminder)
	}

	return res, nil
}

func (db *sharedSpaceRepo) GetSharedSpaceByName(ctx context.Context, name string) (model.SharedSpace, error) {
	space := model.SharedSpace{}

	rows := db.db.QueryRowContext(ctx, `select shared_spaces.shared_spaces.id, name, created, tg_id, username from shared_spaces.shared_spaces 
join users.users on users.users.id = shared_spaces.shared_spaces.creator
where name  like $1;`, name)

	err := rows.Scan(&space.ID, &space.Name, &space.Created, &space.Creator.TGID, &space.Creator.Username)
	if err != nil {
		return model.SharedSpace{}, err
	}

	return space, nil
}
