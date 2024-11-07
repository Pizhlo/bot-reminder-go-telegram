package sharedaccess

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
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
	// хранилище совместных пространств
	storage editor

	viewsMap map[int64]*view.SharedSpaceView
}

type editor interface {
	GetAllByUserID(ctx context.Context, userID int64) ([]model.SharedSpace, error)
	Save(ctx context.Context, space model.SharedSpace) error
}

func New(userSrv *user.UserService,
	noteSrv *note.NoteService,
	reminderSrv *reminder.ReminderService,
	editor editor) *SharedSpace {

	return &SharedSpace{
		userSrv:     userSrv,
		noteSrv:     noteSrv,
		reminderSrv: reminderSrv,
		storage:     editor,
		viewsMap:    make(map[int64]*view.SharedSpaceView),
	}
}

// SaveUser сохраняет пользователя в мапе view
func (n *SharedSpace) SaveUser(userID int64) {
	if _, ok := n.viewsMap[userID]; !ok {
		logrus.Debugf("SharedSpaceSrv: user %d not found in the views map. Saving...\n", userID)
		n.viewsMap[userID] = view.NewSharedSpaceView()
	} else {
		logrus.Debugf("SharedSpaceSrv: user %d already saved in the views map.\n", userID)
	}

}

func (s *SharedSpace) Buttons(userID int64) []telebot.Btn {
	return s.viewsMap[userID].Buttons()
}
