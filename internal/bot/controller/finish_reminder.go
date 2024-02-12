package controller

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

// FinishReminder обрабатывает завершение создания напоминания
func (c *Controller) FinishReminder(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend("finish reminder")
}
