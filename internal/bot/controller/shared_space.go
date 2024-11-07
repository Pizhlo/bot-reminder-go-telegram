package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

	btns := c.sharedSpace.Buttons(telectx.Chat().ID)

	for _, btn := range btns {
		c.bot.Handle(&btn, func(telectx tele.Context) error {
			spaceID, err := strconv.Atoi(btn.Unique)
			if err != nil {
				return fmt.Errorf("error converting string space ID '%s' to int: %+v", btn.Unique, err)
			}

			return c.GetSharedSpace(ctx, telectx, spaceID)
		})
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

func (c *Controller) GetSharedSpace(ctx context.Context, telectx tele.Context, spaceID int) error {
	msg, kb, err := c.sharedSpace.GetSharedSpace(spaceID, telectx.Chat().ID)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(msg, kb)
}

func (c *Controller) GetCurrentSharedSpace(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.CurrentSharedSpace(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(msg, kb)
}

func (c *Controller) NotesBySharedSpace(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.NotesBySpace(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(msg, kb)
}

func (c *Controller) RemindersBySharedSpace(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.RemindersBySpace(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(msg, kb)
}

func (c *Controller) SharedSpaceParticipants(ctx context.Context, telectx tele.Context) error {
	msg, kb := c.sharedSpace.SharedSpaceParticipants(telectx.Chat().ID)

	return telectx.EditOrSend(msg, kb)
}
