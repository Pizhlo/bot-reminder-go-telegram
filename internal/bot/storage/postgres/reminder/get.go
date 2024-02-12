package reminder

import (
	"context"
	"fmt"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// GetAllByUserID достает из базы все напоминания пользователя по ID
func (db *ReminderRepo) GetAllByUserID(ctx context.Context, userID int64) ([]model.Reminder, error) {
	reminders := make([]model.Reminder, 0)

	rows, err := db.db.QueryContext(ctx, `select id, text, created, type, date, time from reminders.reminders where user_id = (select id from users.users where tg_id = $1) order by created ASC`, userID)
	if err != nil {
		return nil, fmt.Errorf("error while getting all reminders from DB by user ID %d: %w", userID, err)
	}
	defer rows.Close()

	for rows.Next() {
		reminder := model.Reminder{}

		err := rows.Scan(&reminder.ID, &reminder.Name, &reminder.Created, &reminder.Type, &reminder.Date, &reminder.Time)
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
