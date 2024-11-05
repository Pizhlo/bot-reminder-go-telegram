package sharedaccess

import (
	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	"gopkg.in/telebot.v3"
)

// SharedSpace - структура, управляющая совместными пространствами
type SharedSpace struct {
	// отвечает за информацию о пользователях
	userSrv *user.UserService
	// отвечает за обработку заметок
	noteSrv *note.NoteService
	// отвечает за напоминания
	reminderSrv *reminder.ReminderService
}

func New(userSrv *user.UserService,
	noteSrv *note.NoteService,
	reminderSrv *reminder.ReminderService,
	channelID int64) *SharedSpace {

	return &SharedSpace{
		userSrv:     userSrv,
		noteSrv:     noteSrv,
		reminderSrv: reminderSrv,
	}
}

func (s *SharedSpace) GetAllByUserID(userID int64) (string, *telebot.ReplyMarkup, error) {
	return "", nil, api_errors.ErrSharedSpacesNotFound
}
