package controller

import (
	"context"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

var (
	// inline кнопка для подтверждения удаления всех напоминаний
	BtnDeleteAllReminders = selector.Data("Да, удалить все напоминания", "delete_all_reminders")
	// inline кнопка для отмены удаления всех напоминаний
	BtnNotDeleteAllReminders = selector.Data("Подожди! Я случайно", "not_delete_all_reminders")
)

// ConfirmDeleteAllReminders отправляет пользователю уточняющее сообщение о том, действительно ли
// он хочет удалить все напоминания. Также отправляет клавиатуру с кнопками подтверждения и отмены
func (c *Controller) ConfirmDeleteAllReminders(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling /reminders_del command. Sending confirmation...\n")

	// сначала проверяем, есть ли напоминания у пользователя
	// _, _, err := c.reminderSrv.GetAll(ctx, telectx.Chat().ID)
	// if err != nil {
	// 	if errors.Is(err, api_errors.ErrNotesNotFound) {
	// 		logrus.Errorf("Controller: cannot delete all user's reminders: user doesn't have any reminders yet. User ID: %d.\n", telectx.Chat().ID)
	// 		return telectx.EditOrSend(messages.UserDoesntHaveRemindersMessage, view.BackToMenuBtn())
	// 	}
	// 	logrus.Errorf("Controller: error while handling /reminders_del command: checking if user has reminders. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)
	// 	return err
	// }

	selector.Inline(
		selector.Row(BtnDeleteAllReminders, BtnNotDeleteAllReminders),
	)

	return telectx.EditOrSend(messages.ConfirmDeleteRemindersMessage, selector)
}

// DeleteAllReminders удаляет все напоминания пользователя
func (c *Controller) DeleteAllReminders(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Controller: handling /reminders_del command.\n")

	err := c.reminderSrv.DeleteAll(ctx, telectx.Chat().ID)
	if err != nil {
		logrus.Errorf("Controller: error while deleting all reminders. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)
		return err
	}

	logrus.Debugf("Controller: successfully delete all user's reminders. Sending message to user...\n")
	return telectx.Edit(messages.AllRemindersDeletedMessage, view.BackToMenuBtn())
}
