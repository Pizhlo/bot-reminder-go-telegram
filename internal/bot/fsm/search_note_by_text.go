package fsm

import (
	"context"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type searchNoteByTextState struct {
	controller *controller.Controller
	fsm        *FSM
	logger     *logrus.Logger
	name       string
}

func newSearchNoteByTextState(controller *controller.Controller, FSM *FSM) *searchNoteByTextState {
	return &searchNoteByTextState{controller, FSM, logger.New(), "search note by text"}
}

func (n *searchNoteByTextState) Handle(ctx context.Context, telectx tele.Context) error {
	n.logger.Debugf("Handling search note by text. State: %s. Message: %s\n", n.Name(), telectx.Message().Text)

	if strings.HasPrefix(telectx.Message().Text, "/del") {
		return n.controller.DeleteNoteByID(ctx, telectx)
	}

	return n.controller.SearchNoteByText(ctx, telectx)
}

func (n *searchNoteByTextState) Name() string {
	return n.name
}
