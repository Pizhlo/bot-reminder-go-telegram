package fsm

import (
	"context"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// Состояние перечисления заметок
type listNote struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       stateName
	next       state
}

func newListNoteState(FSM *FSM, controller *controller.Controller) *listNote {
	return &listNote{controller, FSM, logger.New(), listNoteName, nil}
}

const deleteNotePrefix = "/dn"

func (n *listNote) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	msg := telectx.Message().Text

	if !strings.HasPrefix(msg, deleteNotePrefix) {
		n.fsm.SetToDefault()
		return n.fsm.Handle(ctx, telectx)
	} else {
		return n.controller.DeleteNoteByID(ctx, telectx)
	}

	//return n.controller.ListNotes(ctx, telectx)
}

func (n *listNote) Name() string {
	return string(n.name)
}

func (n *listNote) Next() state {
	if n.next != nil {
		return n.next
	}
	return n.fsm.defaultState
}
