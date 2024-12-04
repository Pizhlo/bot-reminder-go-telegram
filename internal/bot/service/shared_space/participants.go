package sharedaccess

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"gopkg.in/telebot.v3"
)

// SpaceParticipants возвращает участников выбранного совместного пространства
func (s *SharedSpace) SpaceParticipants(userID int64) []model.User {
	return s.viewsMap[userID].CurrentSpace().Participants
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
