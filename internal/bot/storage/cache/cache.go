package cache

import (
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
)

type Cache struct {
	data sync.Map
}

type UserTimezone struct {
	Location string
	Long     string
	Lat      string
}

func New() *Cache {
	return &Cache{}
}

func (c *Cache) SaveUserTimezone(id int64, tz UserTimezone) {
	c.data.Store(id, tz)
}

func (c *Cache) GetUserTimezone(id int64) (UserTimezone, error) {
	val, ok := c.data.Load(id)
	if !ok {
		return UserTimezone{}, errors.UserNotFound
	}

	var userTZ UserTimezone

	userTZ = val.(UserTimezone)

	return userTZ, nil
}
