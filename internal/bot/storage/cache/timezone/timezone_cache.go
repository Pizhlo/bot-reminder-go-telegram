package cache

import (
	"context"
	"sync"
	"time"

	"errors"
)

type TimezoneCache struct {
	mu   sync.Mutex
	data map[int64]*time.Location
}

func New() *TimezoneCache {
	return &TimezoneCache{mu: sync.Mutex{}, data: make(map[int64]*time.Location)}
}

func (c *TimezoneCache) Save(ctx context.Context, id int64, loc *time.Location) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[id] = loc
}

func (c *TimezoneCache) Get(ctx context.Context, id int64) (*time.Location, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	val, ok := c.data[id]
	if !ok {
		return nil, errors.New("no location found for this user")
	}

	return val, nil
}
