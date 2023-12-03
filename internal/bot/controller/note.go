package controller

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/account"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"

	"gopkg.in/telebot.v3"
)

type Note struct {
	baseContext func() context.Context
	noting      note.Service
	accounting  account.Service
}

func NewNote(
	baseContext func() context.Context,
	noting note.Service,
	accounting account.Service,
) *Note {
	return &Note{
		baseContext: baseContext,
		noting:      noting,
		accounting:  accounting,
	}
}

func (p *Note) Handle(telctx telebot.Context) error {
	tid := tgid{telctx}.Get()

	l := slog.With(
		"command", "note",
		"sender id", tid,
	)

	text := telctx.Message().Payload // /note some text
	if len(text) == 0 {
		text = telctx.Text() // plain text message
	}

	if len(strings.TrimSpace(text)) == 0 {
		l.Info("skip handling", "reason", "no payload")

		return nil
	}

	l.Info("handling", "text", text)

	ctx, cancel := context.WithCancel(p.baseContext())
	defer cancel()

	u, err := p.accounting.FindUserByTelegramID(ctx, tid)
	if err != nil {
		return err
	}

	// escaping special characters for Telegram Markdown V2 markup
	text = strings.Replace(text, `_`, `\_`, -1)
	text = strings.Replace(text, `-`, `\-`, -1)

	n, err := p.noting.AddNote(ctx, u.ID, text)
	if err != nil {
		return err
	}

	l.Info("adding a note", "note id", n.ID)

	return nil
}
