package notelist

import (
	"context"
	"errors"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/note"

	tele "gopkg.in/telebot.v3"
)

type TelegramMarkdownView struct {
	ctx    tele.Context
	userID int
	maxLen int
}

type NoteList interface {
	GetViewModel(ctx context.Context) (*note.SearchParams, error)
	SetNotes(ctx context.Context, notes []*note.Note) error
}

var _ NoteList = (*TelegramMarkdownView)(nil)

func NewTelegramMarkdownView(ctx tele.Context, userID int, msgMaxLen int) *TelegramMarkdownView {
	return &TelegramMarkdownView{
		ctx:    ctx,
		userID: userID,
		maxLen: msgMaxLen,
	}
}

func (p *TelegramMarkdownView) GetViewModel(ctx context.Context) (*note.SearchParams, error) {
	return &note.SearchParams{
		UserID: p.userID,
		Terms:  p.ctx.Args(),
	}, nil
}

const tmvMsgTpl = "\\-\\-\\-\n%s\n\\-\\-\\-\nсоздано: %s\nудалить: /dn%d\n\n"
const tmvTimestampFormat = "02\\.01\\.2006 в 15:04:05"
const tmvNothingTpl = "_ничего нигде нет, такие дела_"

func (p *TelegramMarkdownView) SetNotes(ctx context.Context, notes []*note.Note) error {
	var list []string
	currLen := 0
	item := ""

	for _, n := range notes {
		part := fmt.Sprintf(
			tmvMsgTpl,
			n.Text,
			n.UpdatedAt.Format(tmvTimestampFormat),
			n.ID,
		)

		l := len([]rune(part))

		if currLen+l <= p.maxLen {
			item += part
			currLen += l
		} else {
			list = append(list, item)
			item = part
			currLen = l
		}
	}

	if len(item) > 0 {
		list = append(list, item)
	}

	if len(list) == 0 {
		list = append(list, tmvNothingTpl)
	}

	var err error
	for _, msg := range list {
		serr := p.ctx.Send(msg, tele.ModeMarkdownV2)
		if serr != nil {
			err = errors.Join(err, serr)
		}
	}

	if err != nil {
		return err
	}

	return nil
}
