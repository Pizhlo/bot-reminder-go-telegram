package reminder

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (db *ReminderRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return fmt.Errorf("error while creating transaction DeleteAllByUserID: %w", err)
	}

	_, err = tx.ExecContext(ctx, `delete from reminders.reminders where user_id = (select id from users.users where tg_id = $1)`, userID)
	if err != nil {
		return fmt.Errorf("error deleting all reminders: %w", err)
	}

	return tx.Commit()
}

func (db *ReminderRepo) DeleteReminderByID(ctx context.Context, reminderID uuid.UUID) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return fmt.Errorf("error while creating transaction DeleteReminderByID: %w", err)
	}

	_, err = tx.ExecContext(ctx, `delete from reminders.reminders where id = $1`, reminderID)
	if err != nil {
		return fmt.Errorf("error deleting reminder by ID %v: %w", reminderID, err)
	}

	return tx.Commit()
}

func (db *ReminderRepo) DeleteJobAndReminder(ctx context.Context, jobID uuid.UUID) error {
	tx, err := db.tx(ctx)
	if err != nil {
		return fmt.Errorf("error while creating transaction DeleteJobAndReminder: %w", err)
	}

	_, err = tx.ExecContext(ctx, "delete from reminders.reminders where id = (select reminder_id from reminders.jobs where job_id = $1)", jobID)
	if err != nil {
		return fmt.Errorf("error deleting job by ID: %w", err)
	}

	return tx.Commit()
}
