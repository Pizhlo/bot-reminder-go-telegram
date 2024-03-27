package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/sirupsen/logrus"
)

// UserService отвечает за информацию о пользователях: айди, часовой пояс и т.п.
type UserService struct {
	logger *logrus.Logger
	//userCache      userEditor
	userEditor     userEditor
	timezoneCache  timezoneCache  // in-memory cache
	timezoneEditor timezoneEditor // db
}

type userEditor interface {
	Get(ctx context.Context, userID int64) (*user.User, error)
	Save(ctx context.Context, id int64, u *user.User) error
	GetAll(ctx context.Context) ([]*user.User, error) // для восстановления кэша на старте
	SaveState(ctx context.Context, id int64, state string) error
}

type timezoneEditor interface {
	Get(ctx context.Context, userID int64) (*user.Timezone, error)
	Save(ctx context.Context, id int64, tz *user.Timezone) error
	GetAll(ctx context.Context) ([]*user.User, error) // для восстановления кэша на старте
}

type timezoneCache interface {
	Get(ctx context.Context, userID int64) (*time.Location, error)
	Save(ctx context.Context, id int64, tz *time.Location)
}

func New(ctx context.Context, userEditor userEditor, timezoneCache timezoneCache, timezoneEditor timezoneEditor) *UserService {
	srv := &UserService{userEditor: userEditor, logger: logger.New(),
		timezoneCache: timezoneCache, timezoneEditor: timezoneEditor}

	srv.loadAll(ctx)

	return srv
}

func (s *UserService) loadAll(ctx context.Context) {
	tzs, err := s.timezoneEditor.GetAll(ctx)
	if err != nil {
		s.logger.Fatalf("unable to load all timezones from DB on start: %v\n", err)
	}

	for _, tz := range tzs {
		loc, err := time.LoadLocation(tz.Timezone.Name)
		if err != nil {
			s.logger.Fatalf("unable to load user's location on start. Location: %s. Error: %v", tz.Timezone.Name, err)
		}

		s.timezoneCache.Save(ctx, tz.TGID, loc)
	}

	s.logger.Debugf("Successfully saved %d users' timezone(s) to cache\n", len(tzs))
}

func (s *UserService) GetAll(ctx context.Context) ([]*user.User, error) {
	return s.userEditor.GetAll(ctx)
}

func (s *UserService) SaveUser(ctx context.Context, userID int64, u *user.User) error {
	return s.userEditor.Save(ctx, userID, u)
}

func (s *UserService) CheckUser(ctx context.Context, tgID int64) bool {
	if !s.checkInCache(ctx, tgID) {
		return s.checkInRepo(ctx, tgID)
	}
	return true
}

func (s *UserService) checkInCache(ctx context.Context, tgID int64) bool {
	u, err := s.timezoneCache.Get(ctx, tgID)
	if err != nil {
		s.logger.Errorf("User service: Error while checking user in cache: %v\n", err)
	} else {
		s.logger.Debugf("User service: Found user in cache: %+v\n", u)
	}

	if u == nil {
		s.logger.Debugf("User service: User not found in cache: %d\n", tgID)
	}

	return u != nil
}

func (s *UserService) checkInRepo(ctx context.Context, tgID int64) bool {
	u, err := s.userEditor.Get(ctx, tgID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			s.logger.Errorf("User service: Error while checking user in DB: %v\n", err)
		}
		s.logger.Debugf("User service: User not found in DB: %d\n", tgID)
	} else {
		s.logger.Debugf("User service: Found user in DB: %+v\n", u)
	}

	return u != nil
}

func (s *UserService) SaveState(ctx context.Context, tgID int64, state string) error {
	return s.userEditor.SaveState(ctx, tgID, state)
}
