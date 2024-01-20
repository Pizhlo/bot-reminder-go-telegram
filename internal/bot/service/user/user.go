package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	logger         *logrus.Logger
	userCache      userEditor
	userEditor     userEditor
	timezoneCache  timezoneEditor // in-memory cache
	timezoneEditor timezoneEditor // db
}

type userEditor interface {
	Get(ctx context.Context, userID int64) (*user.User, error)
	Save(ctx context.Context, id int64, u *user.User) error
	GetAll(ctx context.Context) ([]*user.User, error) // для восстановления кэша на старте
}

type timezoneEditor interface {
	Get(ctx context.Context, userID int64) (*user.Timezone, error)
	Save(ctx context.Context, id int64, tz *user.Timezone) error
	GetAll(ctx context.Context) ([]*user.User, error) // для восстановления кэша на старте
}

func New(userEditor userEditor, userCache userEditor, timezoneCache timezoneEditor, timezoneEditor timezoneEditor) *UserService {
	srv := &UserService{userEditor: userEditor, userCache: userCache, logger: logger.New(),
		timezoneCache: timezoneCache, timezoneEditor: timezoneEditor}

	srv.logger.Debugf("Loading all users from DB to cache...\n")

	users, err := userEditor.GetAll(context.Background())
	if err != nil {
		srv.logger.Fatalf("unable to load all users from DB on start: %v\n", err)
	}

	for _, user := range users {
		srv.userCache.Save(context.Background(), user.TGID, user)
	}

	srv.logger.Debugf("Successfully saved %d user(s) to cache\n", len(users))

	srv.logger.Debugf("Loading all users' timezones from DB to cache...\n")

	tzs, err := srv.timezoneEditor.GetAll(context.Background())
	if err != nil {
		srv.logger.Fatalf("unable to load all timezones from DB on start: %v\n", err)
	}

	for _, tz := range tzs {
		srv.timezoneCache.Save(context.Background(), tz.TGID, &tz.Timezone)
	}

	srv.logger.Debugf("Successfully saved %d users' timezone(s) to cache\n", len(users))

	return srv
}

func (s *UserService) SaveUser(ctx context.Context, userID int64, u *user.User) error {
	s.logger.Debugf("Saving user in cache... Key: %d. Value: %v\n", userID, u)
	s.userCache.Save(ctx, userID, u)

	s.logger.Debugf("Successfully saved user in cache. Key: %d. Value: %v\n", userID, u)

	s.logger.Debugf("Saving user in DB... Value: %v\n", u)
	return s.userEditor.Save(ctx, userID, u)
}

func (s *UserService) CheckUser(ctx context.Context, tgID int64) bool {
	if !s.checkInCache(ctx, tgID) {
		return s.checkInRepo(ctx, tgID)
	}
	return true
}

func (s *UserService) checkInCache(ctx context.Context, tgID int64) bool {
	s.logger.Debugf("Checking user in cache. ID: %d\n", tgID)

	u, err := s.userCache.Get(ctx, tgID)
	if err != nil {
		s.logger.Errorf("Error while checking user in cache: %v\n", err)
	} else {
		s.logger.Debugf("Found user in cache: %+v\n", u)
	}

	if u == nil {
		s.logger.Debugf("User not found in cache: %d\n", tgID)
	}

	return u != nil
}

func (s *UserService) checkInRepo(ctx context.Context, tgID int64) bool {
	s.logger.Debugf("Checking user in DB. ID: %d\n", tgID)

	u, err := s.userEditor.Get(ctx, tgID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			s.logger.Errorf("Error while checking user in DB: %v\n", err)
		}
		s.logger.Debugf("User not found in DB: %d\n", tgID)
	} else {
		s.logger.Debugf("Found user in DB: %+v\n", u)
	}

	return u != nil
}
