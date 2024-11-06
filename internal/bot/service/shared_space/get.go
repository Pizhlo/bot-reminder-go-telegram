package sharedaccess

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"gopkg.in/telebot.v3"
)

func (s *SharedSpace) GetAllByUserID(ctx context.Context, userID int64) (string, *telebot.ReplyMarkup, error) {
	spaces, err := s.storage.GetAllByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, api_errors.ErrSharedSpacesNotFound) {
			return messages.SharedSpacesNotFoundMessage, s.viewsMap[userID].Keyboard(), nil
		}

		return "", nil, err
	}

	msg := s.viewsMap[userID].Message(spaces)

	return msg, s.viewsMap[userID].Keyboard(), nil
}

func (s *SharedSpace) GetSharedSpace(spaceID int, userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].MessageBySpace(spaceID)
	if err != nil {
		return "", nil, err
	}

	return msg, s.viewsMap[userID].KeyboardForSpace(), nil
}