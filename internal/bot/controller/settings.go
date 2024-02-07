package controller

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) Settings(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend("Настройки: Кнопка в разработке", view.BackToMenuBtn())
}
