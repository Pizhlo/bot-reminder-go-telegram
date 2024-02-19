package controller

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// ListReminders показывает пользователю все его напоминания
func (c *Controller) ListReminders(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.reminderSrv.GetAll(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrRemindersNotFound) {
			return telectx.EditOrSend(messages.NoRemindersMessage, view.CreateReminderAndBackToMenu())
		}
		return err
	}

	return telectx.EditOrSend(msg, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

// NextPageReminders обрабатывает кнопку переключения на следующую страницу
func (c *Controller) NextPageReminders(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling next Reminders page command.\n")
	page, kb := c.reminderSrv.NextPage(telectx.Chat().ID)

	err := telectx.EditOrSend(page, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		switch t := err.(type) {
		case *tele.Error:
			if t.Description == "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message (400)" {
				break
			}
		default:
			return err
		}
	}

	return nil
}

// NextPageReminders обрабатывает кнопку переключения на предыдущую страницу
func (c *Controller) PrevPageReminders(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling previous Reminders page command.\n")
	page, kb := c.reminderSrv.PrevPage(telectx.Chat().ID)

	err := telectx.EditOrSend(page, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		switch t := err.(type) {
		case *tele.Error:
			if t.Description == "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message (400)" {
				break
			}
		default:
			return err
		}
	}

	return nil
}

// NextPageReminders обрабатывает кнопку переключения на последнюю страницу
func (c *Controller) LastPageReminders(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling last Reminders page command.\n")
	page, kb := c.reminderSrv.LastPage(telectx.Chat().ID)

	err := telectx.EditOrSend(page, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		switch t := err.(type) {
		case *tele.Error:
			if t.Description == "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message (400)" {
				break
			}
		default:
			return err
		}
	}

	return nil
}

// NextPageReminders обрабатывает кнопку переключения на первую страницу
func (c *Controller) FirstPageReminders(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling first Reminders page command.\n")
	page, kb := c.reminderSrv.FirstPage(telectx.Chat().ID)

	err := telectx.EditOrSend(page, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей

	if err != nil {
		switch t := err.(type) {
		case *tele.Error:
			if t.Description == "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message (400)" {
				break
			}
		default:
			return err
		}
	}

	return nil
}
