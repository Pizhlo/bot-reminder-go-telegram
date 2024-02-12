package reminder

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// Save сохраняет напоминание в базе данных. Для сохранения требуется: ID пользователя, содержимое напоминания, дата создания
func (db *ReminderRepo) Save(ctx context.Context, reminder *model.Reminder) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into reminders.reminders (user_id, text, created, type, date, time) values((select id from users.users where tg_id=$1), $2, $3, $4, $5) returning id`,
		reminder.TgID, reminder.Text, reminder.Created, reminder.Type, reminder.Date, reminder.Time)
	if err != nil {
		return fmt.Errorf("error inserting reminder: %w", err)
	}

	return tx.Commit()
}
