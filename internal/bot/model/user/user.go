package user

import (
	"errors"
)

type User struct {
	ID       int
	TGID     int
	Timezone Timezone
}

func (p *User) HasTimezone() bool {
	return !p.Timezone.IsUnknown()
}

type Timezone struct {
	Name string
	Lon  float64
	Lat  float64
}

func (o Timezone) IsUnknown() bool {
	return o.Name == "" && o.Lon == 0 && o.Lat == 0
}

var ErrNotFound = errors.New("not found")

// type Repo interface {
// 	Add(ctx context.Context, tgid int) (*User, error)
// 	Get(ctx context.Context, id int) (*User, error)
// 	Update(ctx context.Context, id int, updFun func(*User) (*User, error)) (*User, error)
// 	FindByTelegramID(ctx context.Context, tgid int) (*User, error)
// }
