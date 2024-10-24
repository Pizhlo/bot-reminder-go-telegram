package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	"github.com/sirupsen/logrus"
)

// UserService отвечает за информацию о пользователях: айди, часовой пояс и т.п.
type UserService struct {
	//userCache      userEditor
	userEditor     userEditor
	timezoneCache  *cache.TimezoneCache // in-memory cache
	timezoneEditor timezoneEditor       // db
}

//go:generate mockgen -source ./user.go -destination=../../mocks/user_srv.go -package=mocks
type userEditor interface {
	Get(ctx context.Context, userID int64) (*user.User, error)
	Save(ctx context.Context, id int64, u *user.User) error
	GetAll(ctx context.Context) ([]*user.User, error) // для восстановления кэша на старте
	SaveState(ctx context.Context, id int64, state string) error
	GetState(ctx context.Context, id int64) (string, error)
}

//go:generate mockgen -source ./user.go -destination=../../mocks/user_srv.go -package=mocks
type timezoneEditor interface {
	Get(ctx context.Context, userID int64) (*user.Timezone, error)
	Save(ctx context.Context, id int64, tz *user.Timezone) error
	GetAll(ctx context.Context) ([]*user.User, error) // для восстановления кэша на старте
}

func New(ctx context.Context, userEditor userEditor, timezoneCache *cache.TimezoneCache, timezoneEditor timezoneEditor) *UserService {
	srv := &UserService{
		userEditor:     userEditor,
		timezoneCache:  timezoneCache,
		timezoneEditor: timezoneEditor,
	}

	srv.loadAll(ctx)

	return srv
}

func (s *UserService) loadAll(ctx context.Context) {
	tzs, err := s.timezoneEditor.GetAll(ctx)
	if err != nil {
		logrus.Fatalf(wrap("unable to load all timezones from DB on start: %v\n", err))
	}

	for _, tz := range tzs {
		loc, err := time.LoadLocation(tz.Timezone.Name)
		if err != nil {
			logrus.Fatalf(wrap("unable to load user's location on start. Location: %s. Error: %v", tz.Timezone.Name, err))
		}

		logrus.Debugf(wrap("saving timezone: name: %s", loc.String()))

		if loc == nil {
			logrus.Error(wrap("cannot save user's timezone: timezone is empty"))
			continue
		}

		s.timezoneCache.Save(ctx, tz.TGID, loc)
	}

	logrus.Debugf(wrap("successfully saved %d users' timezone(s) to cache\n", len(tzs)))
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
		logrus.Errorf(wrap("error while checking user in cache: %v\n", err))
	} else {
		logrus.Debugf(wrap("found user in cache: %+v\n", u))
	}

	if u == nil {
		logrus.Debugf(wrap("User not found in cache: %d\n", tgID))
	}

	return u != nil
}

func (s *UserService) checkInRepo(ctx context.Context, tgID int64) bool {
	u, err := s.userEditor.Get(ctx, tgID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			logrus.Errorf(wrap("error while checking user in DB: %v\n", err))
		}
		logrus.Debugf(wrap("user not found in DB: %d\n", tgID))
	} else {
		logrus.Debugf(wrap("found user in DB: %+v\n", u))
	}

	return u != nil
}

func (s *UserService) SaveState(ctx context.Context, tgID int64, state string) error {
	return s.userEditor.SaveState(ctx, tgID, state)
}

func (s *UserService) GetState(ctx context.Context, tgID int64) (string, error) {
	return s.userEditor.GetState(ctx, tgID)
}

func wrap(s string, args ...any) string {
	str := fmt.Sprintf(s, args...)
	return fmt.Sprintf("User service: %s", str)
}
