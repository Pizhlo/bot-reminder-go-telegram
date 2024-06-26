package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние создания заметки
type createNote struct {
	controller *controller.Controller
	fsm        *FSM
	name       stateName
	next       state
}

func newCreateNoteState(controller *controller.Controller, FSM *FSM) *createNote {
	return &createNote{controller, FSM, createNoteName, nil}
}

func (n *createNote) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	return n.controller.CreateNote(ctx, telectx)
}

func (n *createNote) Name() string {
	return string(n.name)
}

func (n *createNote) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
