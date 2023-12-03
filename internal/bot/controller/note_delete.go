package controller

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"gopkg.in/telebot.v3"
)

type NoteDelete struct {
	baseContext func() context.Context
}

var ErrNotImplemented = errors.New("not implemented")

func NewNoteDelete(baseContext func() context.Context) *NoteDelete {
	return &NoteDelete{
		baseContext: baseContext,
	}
}

func (p *NoteDelete) Handle(telctx telebot.Context) error {
	tid := tgid{telctx}.Get()

	l := slog.With(
		"command", "delete note",
		"sender id", tid,
	)

	ctx, cancel := context.WithCancel(p.baseContext())
	defer cancel()

	l.InfoContext(ctx, "handling")

	return fmt.Errorf("delete note: %w", ErrNotImplemented)
}
