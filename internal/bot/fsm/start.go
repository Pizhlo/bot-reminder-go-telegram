package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки команды старт (кнопки меню) от пользователя
type start struct {
	controller   *controller.Controller
	fsm          *FSM
	location     state
	defaultState state
	logger       *logrus.Logger
	name         stateName
	next         state
}

func newStartState(FSM *FSM, controller *controller.Controller, location state, defaultState state) *start {
	return &start{controller, FSM, location, defaultState, logger.New(), startStateName, nil}
}

// Отправляем пользователю запрос геолокации
func (n *start) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	// если пользователь неизвестен - следующим шагом будет геолокация
	if !n.controller.CheckUser(ctx, telectx.Chat().ID) {
		n.logger.Debugf("Start state: user is unknown. Setting state to location.\n")
		n.fsm.SetState(n.location)
	} else {
		n.fsm.SetState(n.defaultState)
	}

	n.logger.Debugf("Start state: calling controller.Start()\n")
	return n.controller.StartCmd(ctx, telectx)
}

func (n *start) Name() string {
	return string(n.name)
}

func (n *start) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
