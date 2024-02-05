package logger

import (
	"context"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func Logging(contxt context.Context, logger *logrus.Logger) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			if ctx.Message().Text != "" {
				logger.Infof("Handling message. Text: %s. User ID: %d. Commands: %v\n", ctx.Message().Text, ctx.Chat().ID, ctx.Args())
			}

			return next(ctx)
		}
	}
}
