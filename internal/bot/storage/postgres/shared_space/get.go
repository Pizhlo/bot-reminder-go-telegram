package sharedspace

import (
	"context"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

func (db *sharedSpaceRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.SharedSpace, error) {
	spaces := make([]model.SharedSpace, 0)

	rows, err := db.db.QueryContext(ctx,
		`select shared_spaces.shared_spaces.id, space_number, name, created, tg_id, username
		from shared_spaces.shared_spaces
		join shared_spaces.shared_spaces_view on shared_spaces.shared_spaces_view.id = shared_spaces.shared_spaces.id
		join users.users on users.users.id = shared_spaces.creator
		where shared_spaces.shared_spaces.creator = (select id from users.users where tg_id = $1)
		order by created ASC;`, userID)

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

func (db *sharedSpaceRepo) getAllParticipants(ctx context.Context, spaceID int) ([]user.User, error) {
	res := make([]user.User, 0)

	rows, err := db.db.QueryContext(ctx, `select tg_id, username from shared_spaces.participants 
join users.users on users.users.id = shared_spaces.participants.user_id
where space_id = $1;`, spaceID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all users for shared space from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		user := user.User{}

		err := rows.Scan(&user.TGID, &user.UsernameSQL)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID while searching all space's participants: %+v", err)
		}

		res = append(res, user)
	}

	return res, nil
}

func (db *sharedSpaceRepo) getAllNotes(ctx context.Context, spaceID int) ([]model.Note, error) {
	res := make([]model.Note, 0)

	rows, err := db.db.QueryContext(ctx, `select notes.notes.id, text, created, last_edit, space_id, tg_id from notes.notes
join users.users on users.users.id = notes.notes.user_id
where space_id = $1;`, spaceID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all notes for shared space from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		note := model.Note{}

		err := rows.Scan(&note.ID, &note.Text, &note.Created, &note.LastEditSql, &note.SpaceID, &note.TgID)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID while searching all space's participants: %+v", err)
		}

		res = append(res, note)
	}

	return res, nil
}

func (db *sharedSpaceRepo) getAllReminders(ctx context.Context, spaceID int) ([]model.Reminder, error) {
	res := make([]model.Reminder, 0)

	rows, err := db.db.QueryContext(ctx, `select reminders.reminders.id, reminder_number, tg_id, text, created, date, time, name as type, reminders.jobs.job_id
	from reminders.reminders
		join reminders.types on reminders.types.id = reminders.reminders.type_id
		join users.users on users.id = reminders.user_id
		join reminders.reminders_view on reminders.reminders_view.id = reminders.reminders.id
		join reminders.jobs on reminders.jobs.reminder_id = reminders.reminders.id
		where space_id = $1
		order by created ASC;`, spaceID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all notes for shared space from DB: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		reminder := model.Reminder{}

		err := rows.Scan(&reminder.ID, &reminder.ViewID, &reminder.TgID, &reminder.Name, &reminder.Created, &reminder.Date, &reminder.Time, &reminder.Type, &reminder.Job.ID)
		if err != nil {
			return nil, fmt.Errorf("error scanning user ID while searching all space's participants: %+v", err)
		}

		res = append(res, reminder)
	}

	return res, nil
}
