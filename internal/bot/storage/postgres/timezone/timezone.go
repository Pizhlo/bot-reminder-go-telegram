package timezone

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/pkg/errors"
)

type TimezoneRepo struct {
	db *sql.DB
}

func New(dbURl string) (*TimezoneRepo, error) {
	conn, err := sql.Open("postgres", dbURl)
	if err != nil {
		return nil, fmt.Errorf("connect open a db driver: %w", err)
	}

	err = conn.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot connect to a db: %w", err)
	}

	return &TimezoneRepo{conn}, nil
}

// SaveUserTimezone сохраняет часовой пояс пользователя. Аргументы: id - id базы данных, timezone - модель часового пояса
func (db *TimezoneRepo) Save(ctx context.Context, id int64, timezone model.UserTimezone) error {
	_, err := db.db.ExecContext(ctx, `insert into timezones(user_id, timezone) values(?, ?)`, id, timezone.Location)
	if err != nil {
		return errors.Wrap(err, `error while saving timezone`)
	}

	return nil
}

func (db *TimezoneRepo) Get(ctx context.Context, userID int64) (model.UserTimezone, error) {
	return model.UserTimezone{}, nil
}
