package server

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

// CheckUser проверяет, зарегистрирован ли пользователь. Если нет - запрашивает геолокацию.
// Если да - обрабатывает запрос
func (s *Server) CheckUser(contxt context.Context) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(ctx tele.Context) error {
			if !s.controller.CheckUser(contxt, ctx.Chat().ID) {
				s.RegisterUserInFSM(ctx.Chat().ID, false)
				return s.controller.RequestLocation(contxt, ctx)
			}
			return next(ctx)
		}
	}
}
