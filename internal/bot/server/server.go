package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/fsm"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type Server struct {
	bot        *tele.Bot
	fsm        map[int64]*fsm.FSM
	controller *controller.Controller
	logger     *logrus.Logger
}

const (
	startCommand = "/start"
)

func New(bot *tele.Bot, controller *controller.Controller) *Server {
	return &Server{bot: bot, fsm: make(map[int64]*fsm.FSM, 0), controller: controller, logger: logger.New()}
}

func (s *Server) Start(ctx context.Context) {
	s.setupBot(ctx)
}

func (s *Server) setupBot(ctx context.Context) {
	s.bot.Use(logger.Logging(ctx, s.logger))

	s.bot.Handle(tele.OnLocation, func(telectx tele.Context) error {
		return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
	})

	s.bot.Handle(startCommand, func(telectx tele.Context) error {
		s.RegisterUser(telectx.Chat().ID)
		return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
	})

	restricted := s.bot.Group()
	restricted.Use(s.Middleware(ctx), logger.Logging(ctx, s.logger))
}

func (s *Server) RegisterUser(userID int64) {
	s.fsm[userID] = fsm.NewFSM(s.controller)
}
