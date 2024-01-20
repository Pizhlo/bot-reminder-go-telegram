package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	tele "gopkg.in/telebot.v3"
)

type start struct {
	controller *controller.Controller
	fsm        *FSM
	next       state
}

func newStartState(FSM *FSM, controller *controller.Controller, next state) *start {
	return &start{controller, FSM, next}
}

// Отправляем пользователю запрос геолокации
func (n *start) Handle(ctx context.Context, telectx tele.Context) error {

	// если пользователь неизвестен - следующим шагом будет геолокация
	if !n.controller.CheckUser(ctx, telectx.Chat().ID) {
		n.fsm.setState(n.next)
	}

	return n.controller.Start(ctx, telectx)
}
