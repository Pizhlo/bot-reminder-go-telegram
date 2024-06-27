package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние для редактирования заметки. В этом состоянии бот сохраняет присланный текст вместо существующего текста
type editNoteState struct {
	fsm        *FSM
	controller *controller.Controller
	name       stateName
	next       state
	noteNumber int
}

func newEditNoteState(controller *controller.Controller, FSM *FSM, noteNumber int) *editNoteState {
	return &editNoteState{fsm: FSM, controller: controller, name: editNoteName, next: nil, noteNumber: noteNumber}
}

func (n *editNoteState) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	n.fsm.SetToDefault()

	return n.controller.UpdateNote(ctx, telectx, n.noteNumber)
}

func (n *editNoteState) Name() string {
	return string(n.name)
}

func (n *editNoteState) Next() state {
	return n.fsm.defaultState
}
