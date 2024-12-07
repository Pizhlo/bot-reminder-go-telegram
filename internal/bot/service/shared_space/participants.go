package sharedaccess

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// SpaceParticipants возвращает участников выбранного совместного пространства
func (s *SharedSpace) SpaceParticipants(userID int64) []model.Participant {
	return s.viewsMap[userID].CurrentSpace().Participants
}

// SharedSpaceParticipants возвращает участников совместного пространства в виде сообщения и меню
func (s *SharedSpace) SharedSpaceParticipants(userID int64) (string, *telebot.ReplyMarkup) {
	msg := s.viewsMap[userID].ParticipantsMessage()

	return msg, view.ParticipantsKeyboard()
}

// AcceptInvitation удаляет приглашение из БД и устанавливает состояние приглашенного участника в added
func (s *SharedSpace) AcceptInvitation(ctx context.Context, from, to model.Participant, spaceID int64) error {
	err := s.storage.DeleteInvitation(ctx, from, to, spaceID)
	if err != nil {
		return err
	}

	return s.storage.SetParticipantState(ctx, to, model.AddedState, spaceID)
}

// DenyInvitation удаляет приглашение и участника пространства из БД
func (s *SharedSpace) DenyInvitation(ctx context.Context, from, to model.Participant, spaceID int64) error {
	err := s.storage.DeleteInvitation(ctx, from, to, spaceID)
	if err != nil {
		return err
	}

	return s.storage.DeleteParticipant(ctx, spaceID, to)
}
