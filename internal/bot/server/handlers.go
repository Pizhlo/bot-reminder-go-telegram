package server

import (
	"context"

	logger "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"gopkg.in/telebot.v3/middleware"
)

func (s *Server) setupHandlers(ctx context.Context) {
	s.bot.Use(logger.Logging(ctx), middleware.AutoRespond())

	restricted := s.bot.Group()
	restricted.Use(s.CheckUser(ctx), logger.Logging(ctx), middleware.AutoRespond())

	// настройка хендлеров по типам
	s.setupMainHandlers(ctx, restricted)
	s.setupCommands(ctx, restricted)
	s.setupNotesHandlers(ctx, restricted)
	s.setupRemindersHandlers(ctx, restricted)
	s.setupCalendarHandlers(ctx, restricted)
}
