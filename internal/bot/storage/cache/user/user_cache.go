package memory

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

type Memory struct {
	store map[int]*user.User   // хранить по тг базы
	uniq  map[int64]*user.User // хранить по тг айди
	lock  *sync.RWMutex
	seq   atomic.Int64
}

//var _ user.Repo = (*Memory)(nil)

func New() *Memory {
	return &Memory{
		store: map[int]*user.User{},
		uniq:  map[int64]*user.User{},
		lock:  &sync.RWMutex{},
	}
}

func (p *Memory) Save(ctx context.Context, id int64, u *user.User) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	if known, ok := p.uniq[u.TGID]; ok {
		return fmt.Errorf("a user %d (%d) alredy exists", known.TGID, u.TGID)
	}

	p.store[u.ID] = u
	p.uniq[u.TGID] = u

	return nil
}

func (p *Memory) Get(ctx context.Context, id int64) (*user.User, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if u, ok := p.uniq[id]; ok {
		return p.clone(u), nil
	}

	return nil, user.ErrNotFound
}

func (p *Memory) GetAll(ctx context.Context) ([]*user.User, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	res := make([]*user.User, 0)

	for _, v := range p.store {
		res = append(res, v)
	}

	return res, nil
}

func (p *Memory) Update(ctx context.Context, id int, updFun func(*user.User) (*user.User, error)) (*user.User, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	u, ok := p.store[id]
	if !ok {
		return nil, user.ErrNotFound
	}

	uu, err := updFun(p.clone(u))
	if err != nil {
		return nil, fmt.Errorf("cannot update a user: %w", err)
	}

	u.Timezone = uu.Timezone

	return p.clone(u), nil
}

func (p *Memory) FindByTelegramID(ctx context.Context, tgid int64) (*user.User, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	if id, ok := p.uniq[tgid]; ok {
		if u, uok := p.store[id.ID]; uok {
			return p.clone(u), nil
		}
	}

	return nil, user.ErrNotFound
}

func (p *Memory) clone(o *user.User) *user.User {
	u := *o
	return &u
}
