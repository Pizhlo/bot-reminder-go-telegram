package timezone

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// SaveUserTimezone сохраняет часовой пояс пользователя. Аргументы: id - id базы данных, timezone - модель часового пояса
func (db *TimezoneRepo) Save(ctx context.Context, id int64, timezone *model.Timezone) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	_, err = tx.ExecContext(ctx, `insert into users.timezones(user_id, timezone) 
	values((select id from users.users where tg_id=$1), $2) 
	on conflict (user_id) do update set timezone=$2`,
		id, timezone.Name)
	if err != nil {
		return fmt.Errorf("error while saving timezone: %w", err)
	}

	return tx.Commit()
}
