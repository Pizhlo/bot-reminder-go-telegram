package account

import (
	"context"
	"errors"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

type Service interface {
	GetUser(ctx context.Context, id int) (*user.User, error)
	AddUser(ctx context.Context, tgid int) (*user.User, error)
	UpdateUser(ctx context.Context, id int, person *user.User) (*user.User, error)
	FindUserByTelegramID(ctx context.Context, tgid int) (*user.User, error)
}

var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")

type Standard struct {
	memoryUsers   memoryRepo
	postgresUsers postgresRepo
}

type postgresRepo interface {
	Add(ctx context.Context, tgid int) (*user.User, error)
	Get(ctx context.Context, id int) (*user.User, error)
	Update(ctx context.Context, id int, updFun func(*user.User) (*user.User, error)) (*user.User, error)
	FindByTelegramID(ctx context.Context, tgid int) (*user.User, error)
}

type memoryRepo interface {
	Add(ctx context.Context, u *user.User) (*user.User, error)
	Get(ctx context.Context, id int) (*user.User, error)
	Update(ctx context.Context, id int, updFun func(*user.User) (*user.User, error)) (*user.User, error)
	FindByTelegramID(ctx context.Context, tgid int) (*user.User, error)
}

var _ Service = (*Standard)(nil)

func NewStandard(memory memoryRepo, postgres postgresRepo) *Standard {
	return &Standard{
		memoryUsers:   memory,
		postgresUsers: postgres,
	}
}

func (p *Standard) GetUser(ctx context.Context, id int) (*user.User, error) {
	var u *user.User

	u, err := p.memoryUsers.Get(ctx, id)
	if err != nil {
		var dbErr error
		u, dbErr = p.postgresUsers.Get(ctx, id)
		if dbErr != nil {
			if errors.Is(dbErr, user.ErrNotFound) {
				return nil, fmt.Errorf("cannot get a %d user: %w", id, ErrUserNotFound)
			}

			return nil, fmt.Errorf("cannot get a %d user: %w", id, err)
		}
	}

	return u, nil
}

func (p *Standard) AddUser(ctx context.Context, tgid int) (*user.User, error) {
	u, err := p.memoryUsers.FindByTelegramID(ctx, tgid)
	if err != nil && !errors.Is(err, user.ErrNotFound) {
		return nil, fmt.Errorf("cannot add a user (%d): %w", tgid, err)
	}

	if u != nil {
		return nil, fmt.Errorf("cannot add a user (%d): %w", tgid, ErrUserAlreadyExists)
	}

	dbUser, err := p.postgresUsers.Add(ctx, tgid)
	if err != nil {
		return nil, fmt.Errorf("cannot add a user into postgres (%d): %w", tgid, err)
	}

	u, err = p.memoryUsers.Add(ctx, dbUser)
	if err != nil {
		return nil, fmt.Errorf("cannot add a user to memory (%d): %w", tgid, err)
	}

	return u, nil
}

func (p *Standard) UpdateUser(ctx context.Context, id int, person *user.User) (*user.User, error) {
	u, err := p.memoryUsers.Update(ctx, id, func(u *user.User) (*user.User, error) {
		return &user.User{
			Timezone: person.Timezone,
		}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update a user: %w", err)
	}

	_, err = p.postgresUsers.Update(ctx, id, func(u *user.User) (*user.User, error) {
		return &user.User{
			Timezone: person.Timezone,
		}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("cannot update a user: %w", err)
	}

	return u, nil
}

func (p *Standard) FindUserByTelegramID(ctx context.Context, tgid int) (*user.User, error) {
	var u *user.User
	var dbErr error

	u, err := p.memoryUsers.FindByTelegramID(ctx, tgid)
	if err != nil {
		u, dbErr = p.postgresUsers.FindByTelegramID(ctx, tgid)
		if dbErr != nil {
			if errors.Is(dbErr, user.ErrNotFound) {
				return nil, fmt.Errorf("cannot find a user (%d): %w", tgid, ErrUserNotFound)
			}

			return nil, fmt.Errorf("cannot find a user (%d): %w", tgid, err)
		}

	}

	return u, nil
}
