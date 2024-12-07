package sharedaccess

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// GetAllByUserID возвращает все совместные доступы пользователя
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

// GetSharedSpace возвращает информацию по конкретному совместному пространству по spaceID
func (s *SharedSpace) GetSharedSpace(spaceID int, userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].MessageBySpace(spaceID)
	if err != nil {
		return "", nil, err
	}

	return msg, view.KeyboardForSpace(), nil
}

// CurrentSharedSpace возвращает информацию о выбранном совместном пространстве
func (s *SharedSpace) CurrentSharedSpace(userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].MessageByCurrentSpace()
	if err != nil {
		return "", nil, err
	}

	return msg, view.KeyboardForSpace(), nil
}

// RemindersBySpace возвращает напоминания, принадлежащие конкретному пространству, которое уже было выбрано, поэтому
// для запроса нужен только userID
func (s *SharedSpace) RemindersBySpace(userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].Reminders()
	if err != nil {
		return "", nil, err
	}

	kb := view.KeyboardForReminders()

	return msg, kb, nil
}

// CurrentSpace возвращает название текущего (выбранного) совметного доступа
func (s *SharedSpace) CurrentSpaceName(userID int64) string {
	return s.viewsMap[userID].CurrentSpaceName()
}

// CurrentSpace возвращает ID текущего (выбранного) совметного доступа
func (s *SharedSpace) CurrentSpaceID(userID int64) int {
	return s.viewsMap[userID].CurrentSpaceID()
}

// CurrentSpace возвращает выбранное совметное пространство
func (s *SharedSpace) CurrentSpace(userID int64) model.SharedSpace {
	return s.viewsMap[userID].CurrentSpace()
}

// SpaceCreator возвращает создателя совместного пространства
func (s *SharedSpace) SpaceCreator(userID int64) model.Participant {
	return s.viewsMap[userID].CurrentSpace().Creator
}
