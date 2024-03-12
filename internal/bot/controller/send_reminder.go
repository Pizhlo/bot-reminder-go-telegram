package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// SendReminder отправляет пользователю напоминание в указанное время
func (c *Controller) SendReminder(ctx context.Context, reminder model.Reminder) error {
	c.logger.Debugf("Sending reminder to: %d\n", reminder.TgID)

	msg, err := view.ReminderMessage(reminder)
	if err != nil {
		return err
	}

	//unique := fmt.Sprintf("%d.%d", reminder.ID, reminder.TgID)

	// получаем кнопку для удаления сработавшего напоминания
	deleteBtn := telebot.Btn{Text: "❌Удалить", Unique: fmt.Sprintf("%d-%d", reminder.ID, reminder.TgID), Data: fmt.Sprintf("%d-%d", reminder.ID, reminder.TgID)}

	kb := &telebot.ReplyMarkup{}

	kb.Inline(
		kb.Row(deleteBtn),
		kb.Row(view.BtnBackToMenu),
	)

	c.bot.Handle(&deleteBtn, func(telectx telebot.Context) error {
		return c.processDeleteReminder(ctx, telectx, reminder.Name)
	})

	_, err = c.bot.Send(&telebot.Chat{ID: reminder.TgID}, msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: kb,
	})

	if err != nil {
		return err
	}

	return nil
}

// processDeleteReminder обрабатывает кнопку "Удалить" у сработавшего напоминания
func (c *Controller) processDeleteReminder(ctx context.Context, telectx telebot.Context, reminderName string) error {
	// userID_reminderID - для удаления
	reminderAndUser := strings.Split(telectx.Callback().Unique, "-")

	reminderID, userID := reminderAndUser[0], reminderAndUser[1]

	// конвертируем айди из string в int
	reminderInt, err := strconv.Atoi(reminderID)
	if err != nil {
		return fmt.Errorf("error while converting string reminder ID to int: %w", err)
	}

	userInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return fmt.Errorf("error while converting string user ID to int64: %w", err)
	}

	// удаляем напоминание из базы и scheduler
	err = c.DeleteReminder(ctx, telectx, reminderInt, userInt)
	if err != nil {
		return err
	}

	// в случае успеха - отправляем сообщение
	msg := fmt.Sprintf(messages.ReminderDeletedMessage, reminderName)

	return telectx.Edit(msg, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.BackToMenuBtn(),
	})
}
