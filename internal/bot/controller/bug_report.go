package controller

import (
	"context"
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

// BugReport обрабатывает баг-репорт от пользователя
func (c *Controller) BugReport(ctx context.Context, telectx tele.Context) error {
	msg := fmt.Sprintf(messages.BugReportMessage, telectx.Text())
	_, err := c.bot.Send(&tele.Chat{ID: c.channelID}, msg)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(messages.BugReportSucessMessage, view.BackToMenuBtn())
}
