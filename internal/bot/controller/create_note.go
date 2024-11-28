package controller

import (
	"context"
	"fmt"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// CreateNote создает новую заметку пользователя
func (c *Controller) CreateNote(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: saving note. Sender: %d\n", telectx.Chat().ID)

	logrus.Debugf("Controller: getting user's timezone. User ID: %d\n", telectx.Chat().ID)

	loc, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	if err != nil {
		return fmt.Errorf("error while loading user's location: %w", err)
	}

	note := model.Note{
		Creator: model.User{
			TGID: telectx.Chat().ID,
		},
		Text:    telectx.Text(),
		Created: time.Now().In(loc),
	}

	err = c.noteSrv.Save(ctx, note)
	if err != nil {
		return fmt.Errorf("error while saving note. User ID: %d. Error: %v", telectx.Chat().ID, err)
	}

	kb := view.NotesAndMenuBtns()

	return telectx.EditOrSend(messages.SuccessfullyCreatedNoteMessage, kb)
}
