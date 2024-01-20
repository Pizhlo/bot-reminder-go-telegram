package fsm

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	tele "gopkg.in/telebot.v3"
)

type listNote struct {
	controller *controller.Controller
	fsm        *FSM
}

func newListNoteState(FSM *FSM, controller *controller.Controller) *listNote {
	return &listNote{controller, FSM}
}

func (n *listNote) Handle(ctx context.Context, telectx tele.Context) error {
	return n.controller.ListNotes(ctx, telectx)
}
