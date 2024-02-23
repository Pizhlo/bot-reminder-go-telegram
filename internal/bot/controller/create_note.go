package controller

import (
	"context"
	"fmt"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// CreateNote создает новую заметку пользователя
func (c *Controller) CreateNote(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: saving note. Sender: %d\n", telectx.Chat().ID)

	c.logger.Debugf("Controller: getting user's timezone. User ID: %d\n", telectx.Chat().ID)

	tz, err := c.userSrv.GetTimezone(ctx, telectx.Chat().ID)
	if err != nil {
		c.logger.Errorf("Controller: error while getting user timezone. User ID: %d. Error: %v\n", telectx.Chat().ID, err)

		return fmt.Errorf("error while getting user timezone. User ID: %d. Error: %v", telectx.Chat().ID, err)
	}

	c.logger.Debug("Controller: successfully got user's timezone")

	loc, err := time.LoadLocation(tz.Name)
	if err != nil {
		c.logger.Errorf("Controller: error while loading location. Location: %s. Error: %v\n", tz.Name, err)

		return fmt.Errorf("error while loading location. Location: %s. Error: %v", tz.Name, err)
	}

	note := model.Note{
		TgID:    telectx.Chat().ID,
		Text:    telectx.Text(),
		Created: time.Now().In(loc),
	}

	err = c.noteSrv.Save(ctx, note)
	if err != nil {
		c.logger.Errorf("Controller: error while saving note. User ID: %d. Error: %v\n", telectx.Chat().ID, err)

		return fmt.Errorf("error while saving note. User ID: %d. Error: %v", telectx.Chat().ID, err)
	}

	kb := view.NotesAndMenuBtns()

	return telectx.EditOrSend(messages.SuccessfullyCreatedNoteMessage, kb)
}
