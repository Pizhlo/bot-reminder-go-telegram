package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/commands"
	tele "gopkg.in/telebot.v3"
)

// setupCommands настраивает хендлеры для текстовых команд: /start, /help, /menu
func (s *Server) setupCommands(ctx context.Context, restricted *tele.Group) {
	// /help command
	s.bot.Handle(commands.HelpCommand, func(telectx tele.Context) error {
		err := s.controller.HelpCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// /start command
	restricted.Handle(commands.StartCommand, func(telectx tele.Context) error {
		err := s.controller.StartCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// /menu command
	restricted.Handle(commands.MenuCommand, func(telectx tele.Context) error {
		err := s.controller.MenuCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})
}
