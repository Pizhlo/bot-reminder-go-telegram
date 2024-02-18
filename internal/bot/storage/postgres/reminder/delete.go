package reminder

import (
	"context"
	"database/sql"
	"fmt"
)

func (db *ReminderRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `delete from reminders.reminders where user_id = (select id from users.users where tg_id = $1)`, userID)
	if err != nil {
		return fmt.Errorf("error deleting all reminders: %w", err)
	}

	return tx.Commit()
}

func (db *ReminderRepo) DeleteReminderByID(ctx context.Context, userID int64, reminderID int) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `delete from reminders.reminders where user_id = (select id from users.users where tg_id = $1) and id = $2`, userID, reminderID)
	if err != nil {
		return fmt.Errorf("error deleting reminder by ID %d: %w", reminderID, err)
	}

	return tx.Commit()
}
