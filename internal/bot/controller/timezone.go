package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"

	tele "gopkg.in/telebot.v3"
)

func (c *Controller) AcceptTimezone(ctx context.Context, telectx tele.Context) error {
	loc := model.UserTimezone{
		Lat:  telectx.Message().Location.Lat,
		Long: telectx.Message().Location.Lng,
	}

	u, err := c.userSrv.ProcessTimezone(ctx, telectx.Chat().ID, loc)
	if err != nil {
		c.logger.Errorf("Error while processing user's location: %v\n", err)
		return err
	}

	msg := fmt.Sprintf(messages.LocationMessage, u.Timezone.Name)
	return telectx.Send(msg)
}
