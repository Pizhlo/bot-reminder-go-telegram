package fsm

import (
	"context"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/commands"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Дефолтное состояние бота, в котором он воспринимает любой текст как заметку
type defaultState struct {
	fsm        *FSM
	controller *controller.Controller
	name       stateName
	next       state
}

func newDefaultState(controller *controller.Controller, FSM *FSM) *defaultState {
	return &defaultState{fsm: FSM, controller: controller, name: defaultStateName, next: nil}
}

func (n *defaultState) Handle(ctx context.Context, telectx tele.Context) error {
	logrus.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	msg := telectx.Message().Text

	if strings.HasPrefix(msg, "/dn") {
		return n.controller.DeleteNoteByID(ctx, telectx)
	}

	// так как бот может в любой момент находиться в дефолтном состоянии, проверяем текст команды
	switch msg {
	// case notesCommand:
	// 	logrus.Debugf("Default state: got /notes command. Calling controller.ListNotes(). Message: %s\n", msg)
	// 	n.fsm.SetState(n.fsm.ListNote)
	// 	return n.controller.ListNotes(ctx, telectx)

	case commands.StartCommand:
		n.fsm.SetState(n.fsm.start)
		logrus.Debugf("Default state: got /start command. Setting state to Start. Message: %s\n", msg)
		return n.fsm.start.Handle(ctx, telectx)

	// case deleteAllNotesCommand:
	// 	logrus.Debugf("Default state: got /notes_del command. Calling controller.ConfirmDeleteAllNotes(). Message: %s\n", msg)
	// 	n.fsm.SetState(n.fsm.ListNote)
	// 	return n.controller.ConfirmDeleteAllNotes(ctx, telectx)

	default:
		logrus.Debugf("Default state: got usual text. Calling controller.CreateNote(). Message: %s\n", msg)
		return n.controller.CreateNote(ctx, telectx)
	}
}

func (n *defaultState) Name() string {
	return string(n.name)
}

func (n *defaultState) Next() state {
	return n.fsm.defaultState
}
