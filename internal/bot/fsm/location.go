package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type location struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newLocationState(FSM *FSM, controller *controller.Controller) *location {
	return &location{controller, FSM, logger.New(), "location"}
}

// Обрабатываем геолокацию пользователя
func (n *location) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.AcceptTimezone(ctx, telectx)
	if err != nil {
		return telectx.Send("Во время обработки произошла ошибка. Повтори попытку позднее")
	}

	n.fsm.SetState(n.fsm.defaultState)

	return nil
}

func (n *location) Name() string {
	return n.name
}
