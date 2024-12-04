package sharedaccess

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// SpaceParticipants возвращает участников выбранного совместного пространства
func (s *SharedSpace) SpaceParticipants(userID int64) []model.User {
	return s.viewsMap[userID].CurrentSpace().Participants
}

// SharedSpaceParticipants возвращает участников совместного пространства в виде сообщения и меню
func (s *SharedSpace) SharedSpaceParticipants(userID int64) (string, *telebot.ReplyMarkup) {
	msg := s.viewsMap[userID].ParticipantsMessage()

	return msg, view.ParticipantsKeyboard()
}
