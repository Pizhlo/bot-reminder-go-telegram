package fsm

import (
	"context"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type listNote struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newListNoteState(FSM *FSM, controller *controller.Controller) *listNote {
	return &listNote{controller, FSM, logger.New(), "list note"}
}

func (n *listNote) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	msg := telectx.Message().Text

	if !strings.HasPrefix(msg, "/del") {
		n.fsm.SetState(n.fsm.DefaultState)
		return n.fsm.Handle(ctx, telectx)
	} else {
		return n.controller.DeleteNoteByID(ctx, telectx)
	}

	//return n.controller.ListNotes(ctx, telectx)
}

func (n *listNote) Name() string {
	return n.name
}
