package cache

import (
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

type UserCache struct {
	data sync.Map
}

func New() *UserCache {
	return &UserCache{}
}

// SaveUser хранит значения [telegram_id : id in db]
func (c *UserCache) SaveUser(id int, tgID int64) {
	u := model.User{
		ID: id,
	}
	c.data.Store(tgID, u)
}

func (c *UserCache) GetUser(tgID int64) (model.User, error) {
	val, ok := c.data.Load(tgID)
	if !ok {
		return model.User{}, errors.UserNotFound
	}

	var user model.User

	user, ok = val.(model.User)
	if !ok {
		return model.User{}, errors.UnableCastVariable
	}

	return user, nil
}
