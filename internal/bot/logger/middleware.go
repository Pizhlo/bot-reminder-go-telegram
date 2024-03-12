package logger

import (
	"context"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func Logging(contxt context.Context, logger *logrus.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			if ctx.Message().Text != "" && !ctx.Message().Sender.IsBot {
				logger.Infof("Handling message. Text: %s. User ID: %d.\n", ctx.Message().Text, ctx.Chat().ID)
			} else if ctx.Callback() != nil {
				logger.Infof("Handling message. Button: %s. User ID: %d.\n", ctx.Callback().Unique, ctx.Chat().ID)
			} else {
				logger.Infof("Handling message. Text: %s. User ID: %d.\n", ctx.Message().Text, ctx.Chat().ID)
			}

			return next(ctx)
		}
	}
}
