package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

var (
	selector = &tele.ReplyMarkup{OneTimeKeyboard: true}

	// inline кнопка для подтверждения удаления всех заметок
	BtnDeleteAllNotes = selector.Data("Да, удалить все заметки", "delete_all_notes")
	// inline кнопка для отмены удаления всех заметок
	BtnNotDeleteAllNotes = selector.Data("Подожди! Я случайно", "not_delete_all_notes")
)

// ConfirmDeleteAllNotes отправляет пользователю уточняющее сообщение о том, действительно ли
// он хочет удалить все заметки. Также отправляет клавиатуру с кнопками подтверждения и отмены
func (c *Controller) ConfirmDeleteAllNotes(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling /notes_del command. Sending confirmation...\n")

	selector.Inline(
		selector.Row(BtnDeleteAllNotes, BtnNotDeleteAllNotes),
	)

	return telectx.EditOrSend(messages.ConfirmDeleteNotesMessage, selector)
}

// DeleteAllNotes удаляет все заметки пользователя
func (c *Controller) DeleteAllNotes(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling /notes_del command.\n")

	err := c.noteSrv.DeleteAll(ctx, telectx.Chat().ID)
	if err != nil {
		logrus.Errorf("Controller: error while handling /notes_del command: deleting all notes. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)
		return err
	}

	logrus.Debugf("Controller: successfully delete all user's notes. Sending message to user...\n")
	return telectx.Edit(messages.AllNotesDeletedMessage, view.BackToMenuBtn())
}
