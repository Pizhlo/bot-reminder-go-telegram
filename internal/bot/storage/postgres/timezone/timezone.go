package timezone

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type TimezoneRepo struct {
	*pgxpool.Pool
}

func New(conn *pgxpool.Pool) *TimezoneRepo {
	return &TimezoneRepo{conn}
}

// SaveUserTimezone сохраняет часовой пояс пользователя. Аргументы: id - id базы данных, timezone - модель часовой пояс
func (db *TimezoneRepo) SaveUserTimezone(id int, timezone model.UserTimezone) error {
	_, err := db.Exec(context.Background(), `insert into timezones(user_id, timezone) values(?, ?)`, id, timezone.Location)
	if err != nil {
		return errors.Wrap(err, `error while saving timezone`)
	}

	return nil
}
