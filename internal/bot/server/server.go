package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/fsm"
	tele "gopkg.in/telebot.v3"
)

type Server struct {
	bot        *tele.Bot
	fsm        map[int64]*fsm.FSM
	controller *controller.Controller
}

const (
	startCommand = "/start"
)

func New(bot *tele.Bot, controller *controller.Controller) *Server {
	return &Server{bot: bot, fsm: make(map[int64]*fsm.FSM, 0)}
}

func (s *Server) Start(ctx context.Context) {
	s.setupBot(ctx)
}

func (s *Server) setupBot(ctx context.Context) {
	s.bot.Handle(startCommand, func(telectx tele.Context) error {
		s.RegisterUser(telectx.Chat().ID)
		return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
	})
}

func (s *Server) RegisterUser(userID int64) {
	s.fsm[userID] = fsm.NewFSM(s.controller)
}
