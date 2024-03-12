package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"

	tele "gopkg.in/telebot.v3"
)

// Timezone обрабатывает нажатие кнопки "Часовой пояс"
func (c *Controller) Timezone(ctx context.Context, telectx tele.Context) error {
	loc, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(messages.TimezoneMessage, loc.String())

	return telectx.EditOrSend(msg, &tele.SendOptions{
		ReplyMarkup: view.TimezoneMenu(),
		ParseMode:   htmlParseMode,
	})
}

// RequestLocation запрашивает геолокацию у пользователя
func (c *Controller) RequestLocation(ctx context.Context, telectx tele.Context) error {
	txt := fmt.Sprintf(messages.StartMessageLocation, telectx.Chat().FirstName)

	return telectx.Send(txt, view.LocationMenu())
}

// AcceptTimezone обрабатывает геолокацию пользователя
func (c *Controller) AcceptTimezone(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling location\n")

	loc := model.UserTimezone{
		Lat:  telectx.Message().Location.Lat,
		Long: telectx.Message().Location.Lng,
	}

	// сохраняем пользователя и часовой пояс
	u, err := c.userSrv.ProcessTimezoneAndSave(ctx, telectx.Chat().ID, loc)
	if err != nil {
		c.logger.Errorf("Controller: error while processing user's location: %v\n", err)
		return err
	}

	// сохраняем пользователя в сервисах
	err = c.saveUser(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	location, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	// запускаем таски пользователя с новым часовым поясом
	err = c.reminderSrv.StartAllJobs(ctx, telectx.Chat().ID, location, c.SendReminder)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(messages.LocationMessage, u.Timezone.Name)
	return telectx.EditOrSend(msg, &tele.SendOptions{
		ParseMode: htmlParseMode,
		ReplyMarkup: &tele.ReplyMarkup{
			RemoveKeyboard: true,
		},
	})
}
