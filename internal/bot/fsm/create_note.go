package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type createNote struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newCreateNoteState(controller *controller.Controller, FSM *FSM) *createNote {
	return &createNote{controller, FSM, logger.New(), "create note"}
}

func (n *createNote) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling request. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)
	return n.controller.CreateNote(ctx, telectx)
}

func (n *createNote) Name() string {
	return n.name
}
