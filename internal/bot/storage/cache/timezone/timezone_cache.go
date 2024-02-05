package cache

import (
	"context"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

type TimezoneCache struct {
	data sync.Map
}

func New() *TimezoneCache {
	return &TimezoneCache{}
}

func (c *TimezoneCache) Save(ctx context.Context, id int64, tz *user.Timezone) error {
	c.data.Store(id, tz)
	return nil
}

func (c *TimezoneCache) GetAll(ctx context.Context) ([]*user.User, error) {
	res := make([]*user.User, 0)
	c.data.Range(func(key, value interface{}) bool {
		u := value.(*user.User)
		res = append(res, u)
		return true
	})

	return res, nil
}

func (c *TimezoneCache) Get(ctx context.Context, id int64) (*user.Timezone, error) {
	val, ok := c.data.Load(id)
	if !ok {
		return nil, errors.ErrUserNotFound
	}

	var userTZ *user.Timezone

	userTZ, ok = val.(*user.Timezone)
	if !ok {
		return nil, errors.ErrUnableCastVariable
	}

	return userTZ, nil
}
