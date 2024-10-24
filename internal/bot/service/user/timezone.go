package user

import (
	"context"
	"errors"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/ringsaturn/tzf"

	_ "time/tzdata"
)

func (s *UserService) ProcessTimezoneAndSave(ctx context.Context, userID int64, location model.UserTimezone) (*user.User, error) {
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		return nil, errors.New(wrap("error creating default finder: %v", err))
	}

	tz := finder.GetTimezoneName(float64(location.Long), float64(location.Lat))

	u := &user.User{
		TGID: userID,
		Timezone: user.Timezone{
			Name: tz,
		},
	}

	loc, err := time.LoadLocation(u.Timezone.Name)
	if err != nil {
		return nil, err
	}

	err = s.SaveUser(ctx, userID, u)
	if err != nil {
		return nil, errors.New(wrap("error saving new user: %v", err))
	}

	err = s.SaveTimezone(ctx, userID, &u.Timezone, loc)
	if err != nil {
		return nil, errors.New(wrap("error saving timezone: %v", err))
	}

	return u, nil
}

func (s *UserService) GetLocation(ctx context.Context, userID int64) (*time.Location, error) {
	return s.timezoneCache.Get(ctx, userID)
}

func (s *UserService) SaveTimezone(ctx context.Context, userID int64, tz *user.Timezone, loc *time.Location) error {
	s.timezoneCache.Save(ctx, userID, loc)

	return s.timezoneEditor.Save(ctx, userID, tz)
}
