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

func New(bot *tele.Bot, controller *controller.Controller) *Server {
	return &Server{bot: bot, fsm: make(map[int64]*fsm.FSM, 0),
		controller: controller,
		logger:     logger.New()}
}

func (s *Server) Start(ctx context.Context) {
	s.loadFSM(ctx)
	s.setupBot(ctx)
}

func (s *Server) RegisterUser(userID int64, known bool) {
	s.fsm[userID] = fsm.NewFSM(s.controller, known)
}

func (s *Server) loadFSM(ctx context.Context) {
	users, err := s.controller.GetAllUsers(ctx)
	if err != nil {
		s.logger.Fatalf("error while loading all users: %v", err)
	}

	for _, user := range users {
		s.RegisterUser(user.TGID, true)

	}

	s.controller.SaveUsers(ctx, users)
}

// HandleError обрабатывает ошибку: устанавливает состояние в дефолтное, передает контроллеру
func (s *Server) HandleError(ctx tele.Context, err error) {
	// обрабатываем ошибку
	s.controller.HandleError(ctx, err, s.fsm[ctx.Chat().ID].Name())

	// устанавливаем состояние в дефолтное
	s.fsm[ctx.Chat().ID].SetState(s.fsm[ctx.Chat().ID].DefaultState)
}
