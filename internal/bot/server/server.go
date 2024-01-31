package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/fsm"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Server struct {
	bot        *tele.Bot
	fsm        map[int64]*fsm.FSM
	controller *controller.Controller
	logger     *logrus.Logger
}

const (
	startCommand = "/start"
	notesCommand = "/notes"
)

func New(bot *tele.Bot, controller *controller.Controller) *Server {
	return &Server{bot: bot, fsm: make(map[int64]*fsm.FSM, 0),
		controller: controller,
		logger:     logger.New()}
}

func (s *Server) Start(ctx context.Context) {
	s.loadFSM(ctx)
	s.setupBot(ctx)
}

func (s *Server) setupBot(ctx context.Context) {
	s.bot.Use(logger.Logging(ctx, s.logger))
	s.bot.Use(middleware.AutoRespond())

	s.bot.Handle(tele.OnLocation, func(telectx tele.Context) error {
		return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
	})

	s.bot.Handle(startCommand, func(telectx tele.Context) error {
		if _, ok := s.fsm[telectx.Chat().ID]; !ok {
			s.RegisterUser(telectx.Chat().ID, false)
		}

		return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
	})

	restricted := s.bot.Group()
	restricted.Use(s.Middleware(ctx), logger.Logging(ctx, s.logger), middleware.AutoRespond())

	restricted.Handle(tele.OnText, func(telectx tele.Context) error {
		return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
	})

	// restricted.Handle(notesCommand, func(telectx tele.Context) error {
	// 	return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
	// })

	// inline

	s.bot.Handle(&view.BtnNextPgNotes, func(c tele.Context) error {
		err := c.Respond()
		if err != nil {
			s.controller.HandleError(c, err)
			return err
		}

		err = s.controller.NextPageNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err)
			return err
		}

		return nil
	})

	s.bot.Handle(&view.BtnPrevPgNotes, func(c tele.Context) error {
		err := c.Respond()
		if err != nil {
			s.controller.HandleError(c, err)
			return err
		}

		err = s.controller.PrevPageNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err)
			return err
		}

		return nil
	})

	s.bot.Handle(&view.BtnLastPgNotes, func(c tele.Context) error {
		err := c.Respond()
		if err != nil {
			s.controller.HandleError(c, err)
			return err
		}

		err = s.controller.LastPageNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err)
			return err
		}

		return nil
	})

	s.bot.Handle(&view.BtnFirstPgNotes, func(c tele.Context) error {
		err := c.Respond()
		if err != nil {
			s.controller.HandleError(c, err)
			return err
		}

		err = s.controller.FirstPageNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err)
			return err
		}

		return nil
	})
}

func (s *Server) RegisterUser(userID int64, known bool) {
	s.fsm[userID] = fsm.NewFSM(s.controller, known)
}

func (s *Server) loadFSM(ctx context.Context) {
	users := s.controller.GetAllUsers(ctx)

	for _, user := range users {
		s.RegisterUser(user.TGID, true)
	}
}
