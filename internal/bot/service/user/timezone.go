package user

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/ringsaturn/tzf"
)

func (s *UserService) ProcessTimezone(ctx context.Context, userID int64, location model.UserTimezone) (*user.User, error) {
	finder, err := tzf.NewDefaultFinder()
	if err != nil {
		return nil, fmt.Errorf("error creating default finder: %w", err)
	}

	tz := finder.GetTimezoneName(float64(location.Long), float64(location.Lat))

	u := &user.User{
		TGID: userID,
		Timezone: user.Timezone{
			Name: tz,
			Lon:  float64(location.Long),
			Lat:  float64(location.Lat),
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

func (s *UserService) GetTimezone(ctx context.Context, userID int64) (*user.Timezone, error) {
	var userTimezone *user.Timezone
	var err error

	userTimezone, err = s.timezoneCache.Get(ctx, userID)
	if err != nil {
		s.logger.Errorf("Error while looking for timezone in cache. User ID: %d. Error: %v\n", userID, err)
		return s.timezoneEditor.Get(ctx, userID)
	}

	return userTimezone, nil
}

func (s *UserService) SaveTimezone(ctx context.Context, userID int64, tz *user.Timezone) error {
	err := s.timezoneCache.Save(ctx, userID, tz)
	if err != nil {
		s.logger.Errorf("Error while saving timezone in cache. User ID: %d. Timezone: %+v. Error: %v\n", userID, tz, err)
		return fmt.Errorf("error while saving timezone in cache. User ID: %d. Timezone: %+v. Error: %v", userID, tz, err)
	}

	return s.timezoneEditor.Save(ctx, userID, tz)
}
