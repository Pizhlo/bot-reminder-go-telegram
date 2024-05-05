package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/fsm"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type Server struct {
	bot        *tele.Bot
	fsm        map[int64]*fsm.FSM
	controller *controller.Controller
}

func New(bot *tele.Bot, controller *controller.Controller) *Server {
	return &Server{
		bot:        bot,
		fsm:        make(map[int64]*fsm.FSM, 0),
		controller: controller,
	}
}

func (s *Server) Start(ctx context.Context) {
	s.loadUsers(ctx)
	s.setupBot(ctx)

	logrus.Info("server started")
}

func (s *Server) RegisterUserInFSM(userID int64) {
	s.fsm[userID] = fsm.NewFSM(s.controller)
}

func (s *Server) loadUsers(ctx context.Context) {
	users, err := s.controller.GetAllUsers(ctx)
	if err != nil {
		logrus.Fatalf("error while loading all users: %v", err)
	}

	s.controller.SaveUsers(ctx, users)

	for _, user := range users {
		s.RegisterUserInFSM(user.TGID)
		err := s.setupState(ctx, user.TGID)
		if err != nil {
			logrus.Errorf("error while setting state for user on startup: %v", err)
		}
	}

	err = s.controller.LoadMemoryReminders(ctx)
	if err != nil {
		logrus.Errorf("error while loading memory reminders: %v", err)
	}

}

// setupState подготавливает бота к диалогу с пользователем.
// Запрашивает последнее сохраненное состояние и устанавливает в него FSM
func (s *Server) setupState(ctx context.Context, userID int64) error {
	lastState, err := s.controller.GetState(ctx, userID)
	if err != nil {
		return err
	}

	return s.fsm[userID].SetFromString(lastState)
}

// HandleError обрабатывает ошибку: устанавливает состояние в дефолтное, передает контроллеру
func (s *Server) HandleError(ctx tele.Context, err error) {
	// обрабатываем ошибку
	s.controller.HandleError(ctx, err, s.fsm[ctx.Chat().ID].Name())

	// устанавливаем состояние в дефолтное
	s.fsm[ctx.Chat().ID].SetToDefault()
}

// Shutdown сохраняет состояния бота в БД
func (s *Server) Shutdown(ctx context.Context) {
	users, err := s.controller.GetAllUsers(ctx)
	if err != nil {
		logrus.Fatalf("error while loading all users: %v", err)
	}

	for _, u := range users {
		current := s.fsm[u.TGID].Current()
		err := s.controller.SaveState(ctx, u.TGID, current.Name())
		if err != nil {
			logrus.Errorf("error while saving user's state on shutdown: %v", err)
		}
	}

	err = s.controller.SaveMemoryReminder(ctx)
	if err != nil {
		logrus.Errorf("error while saving memory reminder on shutdown: %v", err)
	}
}
