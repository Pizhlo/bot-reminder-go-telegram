package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// дефолтное состояние бота
type defaultState struct {
	fsm        *FSM
	controller *controller.Controller
	start      state
	logger     *logrus.Logger
	name       string
}

const (
	startCommand          = "/start"
	notesCommand          = "/notes"
	deleteAllNotesCommand = "/notes_del"
)

func newDefaultState(controller *controller.Controller, FSM *FSM, start state) *defaultState {
	return &defaultState{fsm: FSM, controller: controller, start: start, logger: logger.New(), name: "default"}
}

func (n *defaultState) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	msg := telectx.Message().Text

	// так как бот может в любой момент находиться в дефолтном состоянии, проверяем текст команды
	switch msg {
	case notesCommand:
		n.logger.Debugf("Default state: got /notes command. Calling controller.ListNotes(). Message: %s\n", msg)
		n.fsm.SetState(n.fsm.ListNote)
		return n.controller.ListNotes(ctx, telectx)
	case startCommand:
		//n.fsm.SetState(n.start)
		n.logger.Debugf("Default state: got /start command. Calling controller.Start(). Message: %s\n", msg)
		//return n.fsm.Start.Handle(ctx, telectx)
		return nil
	case deleteAllNotesCommand:
		n.logger.Debugf("Default state: got /notes_del command. Calling controller.ConfirmDeleteAllNotes(). Message: %s\n", msg)
		//n.fsm.SetState(n.fsm.listNote)
		return n.controller.ConfirmDeleteAllNotes(ctx, telectx)
	default:
		// n.logger.Debugf("Default state: got usual text. Calling controller.CreateNote(). Message: %s\n", msg)
		// return n.controller.CreateNote(ctx, telectx)

		//return n.fsm.Start.Handle(ctx, telectx)
		return nil
	}
}

func (n *defaultState) Name() string {
	return n.name
}
