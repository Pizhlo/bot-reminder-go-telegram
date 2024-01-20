package user

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/ringsaturn/tzf"
)

func (s *UserService) ProcessTimezone(ctx context.Context, userID int64, location model.UserTimezone) (*user.User, error) {
	s.logger.Debugf("Processing user's location to set timezone. Lon: %f. Lat: %f.\n", location.Long, location.Lat)

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

	s.logger.Debugf("Successfully processed user's timezone. User ID: %d. Timezone: %v.\n", userID, u.Timezone)

	return u, nil
}

func (s *UserService) GetTimezone(ctx context.Context, userID int64) (*user.Timezone, error) {
	s.logger.Debugf("Looking for user's timezone. UserID: %d\n", userID)

	var userTimezone *user.Timezone
	var err error

	userTimezone, err = s.timezoneCache.Get(ctx, userID)
	if err != nil {
		s.logger.Errorf("Error while looking for timezone in cache. User ID: %d\n", userID)
		return s.timezoneEditor.Get(ctx, userID)
	}

	return userTimezone, nil
}

func (s *UserService) SaveTimezone(ctx context.Context, userID int64, tz *user.Timezone) error {
	s.logger.Debugf("Saving user's timezone. UserID: %d. Timezone: %+v\n", userID, tz)
	s.logger.Debugf("Saving timezone in cache... Key: %d. Value: %v\n", userID, tz)

	err := s.timezoneCache.Save(ctx, userID, tz)
	if err != nil {
		s.logger.Errorf("Error while saving timezone in cache. User ID: %d. Timezone: %+v. Error: %v\n", userID, tz, err)
		return fmt.Errorf("error while saving timezone in cache. User ID: %d. Timezone: %+v. Error: %v", userID, tz, err)
	}

	s.logger.Debugf("Successfully saved timezone in cache. Key: %d. Value: %v\n", userID, tz)

	s.logger.Debugf("Saving timezone in DB... Value: %v\n", tz)
	return s.timezoneEditor.Save(ctx, userID, tz)
}
