package controller

import (
	"context"
	"fmt"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"gopkg.in/telebot.v3"
)

func (c *Controller) CreateNote(ctx context.Context, telectx telebot.Context) error {
	c.logger.Debugf("Saving note. Sender: %d\n", telectx.Chat().ID)

	c.logger.Debugf("Getting user's timezone. User ID: %d\n", telectx.Chat().ID)

	tz, err := c.userSrv.GetTimezone(ctx, telectx.Chat().ID)
	if err != nil {
		c.logger.Errorf("Error while getting user timezone. User ID: %d. Error: %v\n", telectx.Chat().ID, err)
		return fmt.Errorf("error while getting user timezone. User ID: %d. Error: %v", telectx.Chat().ID, err)
	}

	c.logger.Debug("Successfully got user's timezone")

	loc, err := time.LoadLocation(tz.Location)
	if err != nil {
		c.logger.Errorf("Error while loading location. Location: %s. Error: %v\n", tz.Location, err)
		return fmt.Errorf("error while loading location. Location: %s. Error: %v", tz.Location, err)
	}

	note := model.Note{
		TgID:    telectx.Chat().ID,
		Text:    telectx.Text(),
		Created: time.Now().In(loc),
	}

	err = c.noteSrv.Save(ctx, note)
	if err != nil {
		c.logger.Errorf("Error while saving note. User ID: %d. Error: %v\n", telectx.Chat().ID, err)
		return fmt.Errorf("error while saving note. User ID: %d. Error: %v", telectx.Chat().ID, err)
	}

	return telectx.Send(messages.SuccessfullyCreatedNoteMessage)
}
