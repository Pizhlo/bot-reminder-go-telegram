package reminder

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/google/uuid"
)

// Save сохраняет напоминание в базе данных. Для сохранения требуется: ID пользователя, содержимое напоминания, дата создания
func (db *ReminderRepo) Save(ctx context.Context, reminder *model.Reminder) (int64, error) {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return 0, fmt.Errorf("error while creating transaction: %w", err)
	}

	var id int64

	err = tx.QueryRowContext(ctx,
		`insert into reminders.reminders (user_id, text, created, type_id, date, time) 
	values(
		(select id from users.users where tg_id=$1), 
		$2, $3, (select id from reminders.types where name = $4), 
		$5, $6) returning ID`,
		reminder.TgID, reminder.Name, reminder.Created, reminder.Type, reminder.Date, reminder.Time).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting reminder: %w", err)
	}

	return id, tx.Commit()
}

func (db *ReminderRepo) SaveJob(ctx context.Context, userID, reminderID int64, jobID uuid.UUID) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into reminders.jobs (job_id, reminder_id, user_id) values($1, $2, (select id from users.users where tg_id = $3))`, jobID, reminderID, userID)
	if err != nil {
		return fmt.Errorf("error inserting job: %w", err)
	}

	return tx.Commit()
}
