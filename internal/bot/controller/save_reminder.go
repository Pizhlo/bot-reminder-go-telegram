package controller

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

func (c *Controller) SaveReminder(ctx context.Context, telectx tele.Context) error {
	return c.reminderSrv.Save(ctx, telectx.Chat().ID)
}
