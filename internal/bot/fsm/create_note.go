package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	tele "gopkg.in/telebot.v3"
)

type createNote struct {
	controller *controller.Controller
	fsm        *FSM
}

func newCreateNoteState(controller *controller.Controller, FSM *FSM) *createNote {
	return &createNote{controller, FSM}
}

func (n *createNote) Handle(ctx context.Context, telectx tele.Context) error {
	return n.controller.CreateNote(ctx, telectx)
}
