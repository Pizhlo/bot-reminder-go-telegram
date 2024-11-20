package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние обработки username для добавления в совместное пространство в качестве участника
type addParticipant struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
	spaceName  string
}

func NewAddParticipantState(controller *controller.Controller, FSM *FSM, spaceName string) *addParticipant {
	return &addParticipant{controller, FSM, addParticipantState, nil, spaceName}
}

func (n *addParticipant) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	err := n.controller.HandleParticipant(ctx, telectx, n.spaceName)
	if err != nil {
		return err
	}

	n.fsm.SetNext()

	return nil
}

func (n *addParticipant) Name() string {
	return string(n.name)
}

func (n *addParticipant) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
