package server

import (
	"context"
	"errors"

	api_err "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
)

// type Server struct {
// 	NoteEditor          noteEditor
// 	ReminderEditor      reminderEditor
// 	UserEditor          userEditor
// 	TimezoneCacheEditor timezoneCacheEditor
// 	UserCacheEditor     userCacheEditor
// 	Calendar            *calendar.Calendar
// 	State               state
// }

func (s *Server) GetUserID(ctx context.Context, tgID int64) (int, error) {
	id, err := s.getUserFromCache(ctx, tgID)
	if err != nil {
		if errors.Is(err, api_err.ErrUserNotFound) {
			return s.getUserFromStorage(ctx, tgID)
		}
	}

	return id, nil
}

func (s *Server) getUserFromCache(ctx context.Context, tgID int64) (int, error) {
	user, err := s.userCacheEditor.GetUser(ctx, tgID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (s *Server) getUserFromStorage(ctx context.Context, tgID int64) (int, error) {
	user, err := s.userEditor.GetUser(ctx, tgID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
