package controller

import (
	"context"
	"errors"
	"strings"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// ListNotes возвращает все заметки пользователя и отправляет ему
func (c *Controller) ListNotes(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling %s command.\n", telectx.Message().Text)

	text := telectx.Message().Text

	// проверяем, не пришла ли команда удалить конкретную заметку
	if strings.HasPrefix(text, "/del") {
		return c.DeleteNoteByID(ctx, telectx)
	}

	message, kb, err := c.noteSrv.GetAll(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			return telectx.Edit(messages.UserDoesntHaveNotesMessage, view.BackToMenuBtn())
		}

		logrus.Errorf("Error while handling /notes command. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)

		return err
	}

	logrus.Debugf("Controller: successfully got all user's notes. Sending message to user...\n")
	return telectx.Edit(message, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

// NextPageNotes обрабатывает кнопку переключения на следующую страницу
func (c *Controller) NextPageNotes(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling next notes page command.\n")
	next, kb := c.noteSrv.NextPage(telectx.Chat().ID)

	return telectx.Edit(next, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

// NextPageNotes обрабатывает кнопку переключения на предыдущую страницу
func (c *Controller) PrevPageNotes(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling previous notes page command.\n")
	next, kb := c.noteSrv.PrevPage(telectx.Chat().ID)

	err := telectx.Edit(next, &tele.SendOptions{
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

// NextPageNotes обрабатывает кнопку переключения на последнюю страницу
func (c *Controller) LastPageNotes(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling last notes page command.\n")
	next, kb := c.noteSrv.LastPage(telectx.Chat().ID)

	err := telectx.Edit(next, &tele.SendOptions{
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

// NextPageNotes обрабатывает кнопку переключения на первую страницу
func (c *Controller) FirstPageNotes(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling first notes page command.\n")
	next, kb := c.noteSrv.FirstPage(telectx.Chat().ID)

	err := telectx.Edit(next, &tele.SendOptions{
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
