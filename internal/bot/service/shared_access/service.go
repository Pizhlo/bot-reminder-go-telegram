package sharedaccess

import (
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
)

// SharedAccess - структура, управляющая совместными пространствами
type SharedAccess struct {
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
	channelID int64) *SharedAccess {

	return &SharedAccess{
		userSrv:     userSrv,
		noteSrv:     noteSrv,
		reminderSrv: reminderSrv,
	}
}
