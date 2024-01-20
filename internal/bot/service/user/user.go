package user

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	logger         *logrus.Logger
	userEditor     userEditor
	timezoneCache  timezoneEditor // in-memory cache
	timezoneEditor timezoneEditor // db
}

type userEditor interface{}

type timezoneEditor interface {
	Get(ctx context.Context, userID int64) (model.UserTimezone, error)
	Save(ctx context.Context, id int64, tz model.UserTimezone) error
}

func New(userEditor userEditor, timezoneCache timezoneEditor, timezoneEditor timezoneEditor) *UserService {
	return &UserService{userEditor: userEditor, logger: logrus.New(),
		timezoneCache: timezoneEditor, timezoneEditor: timezoneEditor}
}
