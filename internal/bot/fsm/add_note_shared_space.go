package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние обработки username для добавления в совместное пространство в качестве участника
type addNote struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
	spaceName  string
}

func NewAddNoteState(controller *controller.Controller, FSM *FSM, spaceName string) *addNote {
	return &addNote{controller, FSM, addParticipantState, nil, spaceName}
}

func (n *addNote) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	err := n.controller.AddNoteToSharedSpace(ctx, telectx)
	if err != nil {
		return err
	}

	n.fsm.SetNext()

	return nil
}

func (n *addNote) Name() string {
	return string(n.name)
}

func (n *addNote) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
