package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// GetSharedAccess - хендлер для кнопки "Совметный доступ". Выгружает из базы совместные пространства пользователя, либо говорит, что их нет
func (c *Controller) GetSharedAccess(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.GetAllByUserID(ctx, telectx.Chat().ID)
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
	loc, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	space := model.SharedSpace{
		Name:    telectx.Message().Text,
		Created: time.Now().In(loc),
		Creator: user.User{
			TGID:     telectx.Chat().ID,
			Username: telectx.Chat().Username,
		},
	}

	if err := c.sharedSpace.Save(ctx, space); err != nil {
		return err
	}

	msg := fmt.Sprintf(messages.SharedSpaceCreationSuccessMessage, space.Name)
	return telectx.EditOrSend(msg, view.ShowSharedSpacesMenu())
}
