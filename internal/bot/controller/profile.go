package controller

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

func (c *Controller) Profile(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend("Профиль: Кнопка в разработке", view.BackToMenuBtn())
}
