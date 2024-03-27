package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"gopkg.in/telebot.v3"
)

// saveReminder сохраняет напоминание
func (c *Controller) saveReminder(ctx context.Context, telectx telebot.Context) error {
	// достаем часовой пояс пользователя, чтобы установить поле created
	loc, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	if err != nil {
		return fmt.Errorf("error while loading user's timezone: %w", err)
	}

	// сохраняем поле created
	err = c.reminderSrv.SaveCreatedField(telectx.Chat().ID, loc)
	if err != nil {
		return err
	}

	nextRun, err := c.reminderSrv.SaveAndStartReminder(ctx, telectx.Chat().ID, loc, c.SendReminder)
	if err != nil {
		return err
	}

	layout := "02.01.2006 15:04:05"

	r, err := c.reminderSrv.GetFromMemory(telectx.Chat().ID)
	if err != nil {
		return err
	}

	nextRunMsg, err := view.ProcessTypeAndDate(r.Type, r.Date, r.Time)
	if err != nil {
		return err
	}

	var verb string
	// если срабатывает один раз (определенную дату)
	if r.Type == model.DateType {
		verb = "сработает"
	} else { // в остальных случаях срабатывает больше одного раза
		verb = "будет срабатывать"
	}

	c.reminderSrv.Clear(telectx.Chat().ID)

	msg := fmt.Sprintf(messages.SuccessCreationMessage, r.Name, verb, nextRunMsg, nextRun.NextRun.Format(layout))

	return telectx.EditOrSend(msg, &telebot.SendOptions{
		ReplyMarkup: view.BackToMenuAndCreateOneElse(),
		ParseMode:   htmlParseMode,
	})
}
