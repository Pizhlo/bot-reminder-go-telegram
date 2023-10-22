package cache

import (
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

type TimezoneCache struct {
	data sync.Map
}

func New() *TimezoneCache {
	return &TimezoneCache{}
}

func (c *TimezoneCache) SaveUserTimezone(id int64, tz model.UserTimezone) {
	c.data.Store(id, tz)
}

func (c *TimezoneCache) GetUserTimezone(id int64) (model.UserTimezone, error) {
	val, ok := c.data.Load(id)
	if !ok {
		return model.UserTimezone{}, errors.ErrUserNotFound
	}

	var userTZ model.UserTimezone

	userTZ, ok = val.(model.UserTimezone)
	if !ok {
		return model.UserTimezone{}, errors.ErrUserNotFound
	}

	return userTZ, nil
}
