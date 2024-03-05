package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"

	tele "gopkg.in/telebot.v3"
)

// RequestLocation запрашивает геолокацию у пользователя
func (c *Controller) RequestLocation(ctx context.Context, telectx tele.Context) error {
	txt := fmt.Sprintf(messages.StartMessageLocation, telectx.Chat().FirstName)

	return telectx.EditOrSend(txt, view.LocationMenu())
}

// AcceptTimezone обрабатывает геолокацию пользователя
func (c *Controller) AcceptTimezone(ctx context.Context, telectx tele.Context) error {
	c.logger.Debugf("Controller: handling location\n")

	loc := model.UserTimezone{
		Lat:  telectx.Message().Location.Lat,
		Long: telectx.Message().Location.Lng,
	}

	u, err := c.userSrv.ProcessTimezone(ctx, telectx.Chat().ID, loc)
	if err != nil {
		c.logger.Errorf("Controller: error while processing user's location: %v\n", err)
		return err
	}

	msg := fmt.Sprintf(messages.LocationMessage, u.Timezone.Name)
	return telectx.EditOrSend(msg, tele.RemoveKeyboard, &tele.SendOptions{
		ParseMode: htmlParseMode,
	})
}
