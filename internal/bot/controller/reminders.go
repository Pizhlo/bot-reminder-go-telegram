package controller

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) Reminders(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend("Напоминания: Кнопка в разработке", view.BackToMenuBtn())
}
