package controller

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
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
	c.logger.Debugf("Controller: handling /reminders_del command. Sending confirmation...\n")

	// сначала проверяем, есть ли напоминания у пользователя
	_, _, err := c.reminderSrv.GetAll(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrNotesNotFound) {
			c.logger.Errorf("Controller: cannot delete all user's reminders: user doesn't have any reminders yet. User ID: %d.\n", telectx.Chat().ID)
			return telectx.EditOrSend(messages.UserDoesntHaveRemindersMessage, view.BackToMenuBtn())
		}
		c.logger.Errorf("Controller: error while handling /reminders_del command: checking if user has reminders. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)
		return err
	}

	selector.Inline(
		selector.Row(BtnDeleteAllReminders, BtnNotDeleteAllReminders),
	)

	return telectx.EditOrSend(messages.ConfirmDeleteRemindersMessage, selector)
}

// DeleteAllReminders удаляет все напоминания пользователя
func (c *Controller) DeleteAllReminders(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling /reminders_del command.\n")

	// получаем айди всех задач пользователя, чтобы их остановить в шедулере
	jobIDs, err := c.reminderSrv.GetAllJobs(ctx, telectx.Chat().ID)
	if err != nil {
		c.logger.Errorf("Controller: error getting all jobs' IDs. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)
		return err
	}

	// удаляем из базы
	err = c.reminderSrv.DeleteAll(ctx, telectx.Chat().ID)
	if err != nil {
		c.logger.Errorf("Controller: error while deleting all reminders. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)
		return err
	}

	sch, err := c.getScheduler(telectx.Chat().ID)
	if err != nil {
		return err
	}

	// удаляем из шедулера
	for _, id := range jobIDs {
		sch.DeleteJob(id)
	}

	c.logger.Debugf("Controller: successfully delete all user's reminders. Sending message to user...\n")
	return telectx.Edit(messages.AllRemindersDeletedMessage, view.BackToMenuBtn())
}
