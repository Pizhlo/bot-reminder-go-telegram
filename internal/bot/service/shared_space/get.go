package sharedaccess

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
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

	return msg, s.viewsMap[userID].KeyboardForSpace(), nil
}

// CurrentSharedSpace возвращает информацию о выбранном совместном пространстве
func (s *SharedSpace) CurrentSharedSpace(userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].MessageByCurrentSpace()
	if err != nil {
		return "", nil, err
	}

	return msg, s.viewsMap[userID].KeyboardForSpace(), nil
}

// NotesBySpace возвращает заметки, принадлежащие конкретному пространству, которое уже было выбрано, поэтому
// для запроса нужен только userID
func (s *SharedSpace) NotesBySpace(userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].Notes()
	if err != nil {
		return "", nil, err
	}

	kb := s.viewsMap[userID].KeyboardForNotes()

	return msg, kb, nil
}

// RemindersBySpace возвращает напоминания, принадлежащие конкретному пространству, которое уже было выбрано, поэтому
// для запроса нужен только userID
func (s *SharedSpace) RemindersBySpace(userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].Reminders()
	if err != nil {
		return "", nil, err
	}

	kb := s.viewsMap[userID].KeyboardForReminders()

	return msg, kb, nil
}

// SharedSpaceParticipants возвращает участников совместного пространства в виде сообщения и меню
func (s *SharedSpace) SharedSpaceParticipants(userID int64) (string, *telebot.ReplyMarkup) {
	msg := s.viewsMap[userID].ParticipantsMessage()

	kb := s.viewsMap[userID].ParticipantsKeyboard()

	return msg, kb
}

// InvintationKeyboard возвращает клавиатуру для приглашения пользователя с кнопками "Принять" и "Отклонить"
func (s *SharedSpace) InvintationKeyboard(userID int64) *telebot.ReplyMarkup {
	return s.viewsMap[userID].InvintationKeyboard()
}

// InvintationKeyboard возвращает клавиатуру с кнопкой "Назад в совместное пространство"
func (s *SharedSpace) BackToSharedSpaceMenu(userID int64) *telebot.ReplyMarkup {
	return s.viewsMap[userID].BackToSharedSpaceMenu()
}

// CurrentSpace возвращает название текущего (выбранного) совметного доступа
func (s *SharedSpace) CurrentSpaceName(userID int64) string {
	return s.viewsMap[userID].CurrentSpaceName()
}

// SpaceParticipants возвращает участников выбранного совместного пространства
func (s *SharedSpace) SpaceParticipants(userID int64) []model.User {
	return s.viewsMap[userID].CurrentSpace().Participants
}
