package user

import (
	"context"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/ringsaturn/tzf"
)

func (s *UserService) ProcessTimezoneAndSave(ctx context.Context, userID int64, location model.UserTimezone) (*user.User, error) {
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		return nil, fmt.Errorf("error creating default finder: %w", err)
	}

	tz := finder.GetTimezoneName(float64(location.Long), float64(location.Lat))

	u := &user.User{
		TGID: userID,
		Timezone: user.Timezone{
			Name: tz,
		},
	}

	err = s.SaveUser(ctx, userID, u)
	if err != nil {
		return nil, fmt.Errorf("error saving new user: %w", err)
	}

	err = s.SaveTimezone(ctx, userID, &u.Timezone)
	if err != nil {
		return nil, fmt.Errorf("error saving timezone: %w", err)
	}

	return u, nil
}

func (s *UserService) GetLocation(ctx context.Context, userID int64) (*time.Location, error) {
	return s.timezoneCache.Get(ctx, userID)
}

func (s *UserService) SaveTimezone(ctx context.Context, userID int64, tz *user.Timezone) error {
	loc, err := time.LoadLocation(tz.Name)
	if err != nil {
		return err
	}

	s.timezoneCache.Save(ctx, userID, loc)

	return s.timezoneEditor.Save(ctx, userID, tz)
}
