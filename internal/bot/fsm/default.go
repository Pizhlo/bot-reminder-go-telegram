package fsm

import (
	"context"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/commands"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Дефолтное состояние бота, в котором он воспринимает любой текст как заметку
type defaultState struct {
	fsm        *FSM
	controller *controller.Controller
	logger     *logrus.Logger
	name       string
	next       state
}

func newDefaultState(controller *controller.Controller, FSM *FSM) *defaultState {
	return &defaultState{fsm: FSM, controller: controller, logger: logger.New(), name: "default", next: nil}
}

func (n *defaultState) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	msg := telectx.Message().Text

	if strings.HasPrefix(msg, "/del") {
		return n.controller.DeleteNoteByID(ctx, telectx)
	}

	// так как бот может в любой момент находиться в дефолтном состоянии, проверяем текст команды
	switch msg {
	// case notesCommand:
	// 	n.logger.Debugf("Default state: got /notes command. Calling controller.ListNotes(). Message: %s\n", msg)
	// 	n.fsm.SetState(n.fsm.ListNote)
	// 	return n.controller.ListNotes(ctx, telectx)

	case commands.StartCommand:
		n.fsm.SetState(n.fsm.Start)
		n.logger.Debugf("Default state: got /start command. Setting state to Start. Message: %s\n", msg)
		return n.fsm.Start.Handle(ctx, telectx)

	// case deleteAllNotesCommand:
	// 	n.logger.Debugf("Default state: got /notes_del command. Calling controller.ConfirmDeleteAllNotes(). Message: %s\n", msg)
	// 	n.fsm.SetState(n.fsm.ListNote)
	// 	return n.controller.ConfirmDeleteAllNotes(ctx, telectx)

	default:
		n.logger.Debugf("Default state: got usual text. Calling controller.CreateNote(). Message: %s\n", msg)
		return n.controller.CreateNote(ctx, telectx)
	}
}

func (n *defaultState) Name() string {
	return n.name
}

func (n *defaultState) Next() state {
	return n.fsm.DefaultState
}
