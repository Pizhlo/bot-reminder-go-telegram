package sharedaccess

import (
	"context"

	"gopkg.in/telebot.v3"
)

func (s *SharedSpace) GetAllByUserID(ctx context.Context, userID int64) (string, *telebot.ReplyMarkup, error) {
	spaces, err := s.storage.GetAllByUserID(ctx, userID)
	if err != nil {
		return "", nil, err
	}

	msg := s.viewsMap[userID].Message(spaces)

	return msg, s.viewsMap[userID].Keyboard(), nil
}
