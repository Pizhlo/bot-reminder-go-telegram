package controller

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

// inline кнопка для удаления сработавшего напоминания
//var DeleteBtn = telebot.Btn{Text: "❌Удалить"}

// SendReminder отправляет пользователю напоминание в указанное время
func (c *Controller) SendReminder(ctx context.Context, reminder *model.Reminder) error {
	logrus.Debugf("Sending reminder to: %d\n", reminder.TgID)

	msg, err := view.ReminderMessage(reminder)
	if err != nil {
		return err
	}

	// передаем только айди напоминания, потому что айди пользователя сможем выяснить из контекста
	unique := fmt.Sprintf("%d", reminder.ViewID)

	deleteBtn := telebot.Btn{Text: "❌Удалить", Unique: unique}

	kb := &telebot.ReplyMarkup{}

	kb.Inline(
		kb.Row(view.BtnCheckReminder),
		kb.Row(deleteBtn),
		kb.Row(view.BtnBackToMenu),
	)

	c.bot.Handle(&deleteBtn, func(telectx telebot.Context) error {
		return c.ProcessDeleteReminder(ctx, telectx)
	})

	_, err = c.bot.Send(&telebot.Chat{ID: reminder.TgID}, msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: kb,
	})

	if err != nil {
		logrus.Errorf("error while sending reminder to user: %v", err)
		return err
	}

	return nil
}

// ProcessDeleteReminder обрабатывает кнопку "Удалить" у сработавшего напоминания
func (c *Controller) ProcessDeleteReminder(ctx context.Context, telectx telebot.Context) error {
	reminderID := telectx.Callback().Unique

	reminderInt, err := strconv.Atoi(reminderID)
	if err != nil {
		return fmt.Errorf("error while converting string reminder ID to int: %w", err)
	}

	reminderName, err := c.reminderSrv.DeleteByViewID(ctx, telectx.Chat().ID, reminderInt)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			// значит, что напоминание уже автоматически удалилось
			return err
		}
	}

	msg := fmt.Sprintf(messages.ReminderDeletedMessage, reminderName)

	return telectx.Edit(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToRemindersAndMenu(),
	})
}
