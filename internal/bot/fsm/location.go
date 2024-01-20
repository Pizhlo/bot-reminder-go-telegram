package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	tele "gopkg.in/telebot.v3"
)

type location struct {
	controller *controller.Controller
	fsm        *FSM
}

func newLocationState(FSM *FSM, controller *controller.Controller) *location {
	return &location{controller, FSM}
}

// Обрабатываем геолокацию пользователя
func (n *location) Handle(ctx context.Context, telectx tele.Context) error {

	err := n.controller.AcceptTimezone(ctx, telectx)
	if err != nil {
		return telectx.Send("Во время обработки произошла ошибка. Повтори попытку позднее")
	}

	return nil
}
