package controller

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) HandleSharedAccess(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.GetAllByUserID(telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrSharedSpacesNotFound) {
			return telectx.Edit(messages.SharedSpacesNotFoundMessage, view.SharedAccessMenu())
		}

		logrus.Errorf("Error while handling shared spaces command. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)

		return err
	}

	logrus.Debugf("Controller: successfully got all user's shared spaces. Sending message to user...\n")
	return telectx.Edit(msg, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

func (c *Controller) CreateSharedSpace(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend("совместное пространство создано")
}
