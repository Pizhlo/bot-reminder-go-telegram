package reminder

import (
	"context"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/google/uuid"
)

// GetAllByUserID достает из базы все напоминания пользователя по ID
func (db *ReminderRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.Reminder, error) {
	reminders := make([]model.Reminder, 0)

	rows, err := db.db.QueryContext(ctx, `select reminders.reminders.id, text, created, date, time, name as type from reminders.reminders
	join reminders.types on reminders.types.id = reminders.reminders.type_id
	where user_id = (select id from users.users where tg_id = $1) order by created ASC`, userID)

	if err != nil {
		return nil, fmt.Errorf("error while getting all reminders from DB by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		reminder := model.Reminder{}

		err := rows.Scan(&reminder.ID, &reminder.Name, &reminder.Created, &reminder.Date, &reminder.Time, &reminder.Type)
		if err != nil {
			return nil, fmt.Errorf("error while scanning reminder: %w", err)
		}

		reminders = append(reminders, reminder)
	}

	if len(reminders) == 0 {
		return nil, api_errors.ErrRemindersNotFound
	}

	return reminders, nil
}

func (db *ReminderRepo) GetAllJobs(ctx context.Context, userID int64) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0)

	rows, err := db.db.QueryContext(ctx, `select job_id from reminders.jobs where user_id = (select id from users.users where tg_id = $1)`, userID)

	if err != nil {
		return nil, fmt.Errorf("error while getting all reminders from DB by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		id := uuid.UUID{}

		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error while scanning reminder: %w", err)
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		return nil, api_errors.ErrRemindersNotFound
	}

	return ids, nil
}

func (db *ReminderRepo) GetJobID(ctx context.Context, userID int64, reminderID int) (uuid.UUID, error) {
	var id uuid.UUID

	row := db.db.QueryRowContext(ctx, `select job_id from reminders.jobs where user_id = (select id from users.users where tg_id = $1) and id = $2`, userID, reminderID)

	err := row.Scan(&id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while scanning job ID: %w", err)
	}

	return id, nil
}
