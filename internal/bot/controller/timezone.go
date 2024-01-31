package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"

	tele "gopkg.in/telebot.v3"
)

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

		c.HandleError(telectx, err)

		return err
	}

	msg := fmt.Sprintf(messages.LocationMessage, u.Timezone.Name)
	return telectx.Send(msg, tele.RemoveKeyboard, &tele.SendOptions{
		ParseMode: htmlParseMode,
	})
}
