package user

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

func (s *UserService) GetTimezone(ctx context.Context, userID int64) (model.UserTimezone, error) {
	s.logger.Debugf("Looking for user's timezone. UserID: %d\n", userID)

	var userTimezone model.UserTimezone
	var err error

	userTimezone, err = s.timezoneCache.Get(ctx, userID)
	if err != nil {
		s.logger.Errorf("Error while looking for timezone in cache. User ID: %d\n", userID)
		return s.timezoneEditor.Get(ctx, userID)
	}

	return userTimezone, nil
}

func (s *UserService) SaveTimezone(ctx context.Context, userID int64, tz model.UserTimezone) error {
	s.logger.Debugf("Saving user's timezone. UserID: %d. Timezone: %+v\n", userID, tz)

	err := s.timezoneCache.Save(ctx, userID, tz)
	if err != nil {
		s.logger.Errorf("Error while saving timezone in cache. User ID: %d. Timezone: %+v\n", userID, tz)
		return fmt.Errorf("error while saving timezone in cache. User ID: %d. Timezone: %+v", userID, tz)
	}

	return s.timezoneEditor.Save(ctx, userID, tz)
}
