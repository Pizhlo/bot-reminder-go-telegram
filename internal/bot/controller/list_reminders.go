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

// ListReminders показывает пользователю все его напоминания
func (c *Controller) ListReminders(ctx context.Context, telectx tele.Context) error {
	reminders, err := c.reminderSrv.GetAll(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrRemindersNotFound) {
			return telectx.EditOrSend(messages.NoRemindersMessage, view.CreateReminderAndBackToMenu())
		}
		return err
	}

	msg, err := c.reminderSrv.Message(telectx.Chat().ID, reminders)
	if err != nil {
		return err
	}

	kb := c.reminderSrv.Keyboard(telectx.Chat().ID)

	return telectx.EditOrSend(msg, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

// NextPageReminders обрабатывает кнопку переключения на следующую страницу
func (c *Controller) NextPageReminders(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling next Reminders page command.\n")
	page := c.reminderSrv.NextPage(telectx.Chat().ID)
	kb := c.reminderSrv.Keyboard(telectx.Chat().ID)

	err := telectx.EditOrSend(page, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		if !errors.Is(err, tele.ErrMessageNotModified) {
			return err
		}
	}

	return nil
}

// NextPageReminders обрабатывает кнопку переключения на предыдущую страницу
func (c *Controller) PrevPageReminders(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling previous Reminders page command.\n")
	page := c.reminderSrv.PrevPage(telectx.Chat().ID)
	kb := c.reminderSrv.Keyboard(telectx.Chat().ID)

	err := telectx.EditOrSend(page, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		if !errors.Is(err, tele.ErrMessageNotModified) {
			return err
		}
	}

	return nil
}

// NextPageReminders обрабатывает кнопку переключения на последнюю страницу
func (c *Controller) LastPageReminders(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling last Reminders page command.\n")
	page := c.reminderSrv.LastPage(telectx.Chat().ID)
	kb := c.reminderSrv.Keyboard(telectx.Chat().ID)

	err := telectx.EditOrSend(page, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		if !errors.Is(err, tele.ErrMessageNotModified) {
			return err
		}
	}

	return nil
}

// NextPageReminders обрабатывает кнопку переключения на первую страницу
func (c *Controller) FirstPageReminders(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling first Reminders page command.\n")
	page := c.reminderSrv.FirstPage(telectx.Chat().ID)
	kb := c.reminderSrv.Keyboard(telectx.Chat().ID)

	err := telectx.EditOrSend(page, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей

	if err != nil {
		if !errors.Is(err, tele.ErrMessageNotModified) {
			return err
		}
	}

	return nil
}
