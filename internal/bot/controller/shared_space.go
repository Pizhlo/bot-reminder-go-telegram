package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// AddParticipant обрабатывает нажатие на кнопку "Добавить участника"
func (c *Controller) AddParticipant(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.AddParticipantMessage, view.BackToMenuBtn())
}

// HandleParticipant обрабатывает ссылку на пользователя, которого надо добавить в совместное пространство
func (c *Controller) HandleParticipant(ctx context.Context, telectx tele.Context, spaceName string) error {
	user := telectx.Message().Text
	from := telectx.Chat().Username

	// если пользователь прислал username с собакой
	if strings.HasPrefix(user, "@") {
		username := strings.TrimPrefix(user, "@")
		return c.handleUsername(ctx, telectx, username, from, spaceName)
	}

	// если пользователь прислал ссылку
	if strings.Contains(user, "t.me") {
		return c.handleUserLink(ctx, telectx, user, from, spaceName)
	}

	// если нет собаки, скорее всего это тоже username
	return c.handleUsername(ctx, telectx, user, from, spaceName)
}

// handleUsername обрабатывает username во время добавления пользователя как нового участника пространства
func (c *Controller) handleUsername(ctx context.Context, telectx tele.Context, username, from, spaceName string) error {
	user, err := c.userSrv.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, api_errors.ErrUserNotFound) {
			return telectx.EditOrSend(messages.UserNotRegisteredMessage, view.BackToMenuBtn())
		}

		return err
	}

	msg := fmt.Sprintf(messages.InvintationMessage, from, spaceName)

	_, err = c.bot.Send(&tele.Chat{ID: user.TGID}, msg, c.sharedSpace.InvintationKeyboard(telectx.Chat().ID))
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.SuccessfullySentInvintationMessage)
}

// handleUserLink обрабатывает ссылку на пользователя во время добавления его как нового участника пространства
func (c *Controller) handleUserLink(ctx context.Context, telectx tele.Context, user, from, spaceName string) error {
	// https://t.me/iloveprogramming

	link := strings.Split(user, "/")

	if len(link) != 4 {
		return telectx.EditOrSend(messages.InvalidUserLinkMessage, view.BackToMenuBtn())
	}

	username := link[3]

	return c.handleUsername(ctx, telectx, username, from, spaceName)
}

func (c *Controller) SpaceName(ctx context.Context, telectx tele.Context) string {
	return c.sharedSpace.CurrentSpace(telectx.Chat().ID)
}
