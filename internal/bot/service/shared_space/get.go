package sharedaccess

import (
	"context"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

func (s *SharedSpace) GetAllByUserID(ctx context.Context, userID int64) (string, *telebot.ReplyMarkup, error) {
	spaces, err := s.storage.GetAllByUserID(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	msg := s.viewsMap[userID].Message(spaces)

	return msg, view.BackToMenuBtn(), api_errors.ErrSharedSpacesNotFound
}
