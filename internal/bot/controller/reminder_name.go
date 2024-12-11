package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// ReminderName обрабатывает название напоминания
func (c *Controller) ReminderName(ctx context.Context, telectx telebot.Context) error {
	if !telectx.Message().Sender.IsBot {
		c.reminderSrv.SaveName(telectx.Chat().ID, telectx.Message().Text)
	}

	// если это делается для совместного пространства, сохраняем информацию в напоминание
	if space := c.sharedSpace.CurrentSpace(telectx.Chat().ID); space.ID != 0 {
		err := c.reminderSrv.SaveSpace(telectx.Chat().ID, &space)
		if err != nil {
			return err
		}

		err = c.reminderSrv.SaveUsers(telectx.Chat().ID, space.Participants)
		if err != nil {
			return err
		}
	}

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return err
	}

	txt := fmt.Sprintf(messages.TypeOfReminderMessage, r.Name)

	return telectx.EditOrSend(txt, &telebot.SendOptions{
		ParseMode:   htmlParseMode,
		ReplyMarkup: view.ReminderTypes(),
	})
}
