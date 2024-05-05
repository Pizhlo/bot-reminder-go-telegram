package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для обработки геолокации от пользователя
type location struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newLocationState(FSM *FSM, controller *controller.Controller) *location {
	return &location{controller, FSM, locationStateName, nil}
}

// Обрабатываем геолокацию пользователя
func (n *location) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	err := n.controller.AcceptTimezone(ctx, telectx)
	if err != nil {
		return telectx.EditOrSend("Во время обработки произошла ошибка. Повтори попытку позднее", view.BackToMenuBtn())
	}

	n.fsm.SetToDefault()

	return nil
}

func (n *location) Name() string {
	return string(n.name)
}

func (n *location) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
