package controller

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

// DeleteReminderByViewID удаляет напоминание по айди, который видит пользователь
func (c *Controller) DeleteReminderByViewID(ctx context.Context, telectx tele.Context, viewID int) error {
	_, err := c.reminderSrv.DeleteByViewID(ctx, telectx.Chat().ID, viewID)
	if err != nil {
		return err
	}

	return nil
}
