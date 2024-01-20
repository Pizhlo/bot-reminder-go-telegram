package server

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

func (s *Server) Middleware(contxt context.Context) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			if !s.controller.CheckUser(contxt, ctx.Chat().ID) {
				s.RegisterUser(ctx.Chat().ID)
				return s.controller.Location(contxt, ctx)
			}
			return next(ctx)
		}
	}
}
