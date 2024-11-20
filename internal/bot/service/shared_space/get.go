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

// GetSharedSpace возвращает информацию по конкретному совместному пространству по spaceID
func (s *SharedSpace) GetSharedSpace(spaceID int, userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].MessageBySpace(spaceID)
	if err != nil {
		return "", nil, err
	}

	return msg, s.viewsMap[userID].KeyboardForSpace(), nil
}

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
	msg := s.viewsMap[userID].Notes()
	kb := s.viewsMap[userID].KeyboardForNotes()

	return msg, kb, nil
}

func (s *SharedSpace) RemindersBySpace(userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].Reminders()
	if err != nil {
		return "", nil, err
	}

	kb := s.viewsMap[userID].KeyboardForReminders()

	return msg, kb, nil
}

func (s *SharedSpace) SharedSpaceParticipants(userID int64) (string, *telebot.ReplyMarkup) {
	msg := s.viewsMap[userID].ParticipantsMessage()

	kb := s.viewsMap[userID].ParticipantsKeyboard()

	return msg, kb
}

func (s *SharedSpace) InvintationKeyboard(userID int64) *telebot.ReplyMarkup {
	return s.viewsMap[userID].InvintationKeyboard()
}
