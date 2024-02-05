package bot

import "gopkg.in/telebot.v3"

type LocationDialog struct {
	markup *telebot.ReplyMarkup
	share  *telebot.Btn
	cancel *telebot.Btn
	text   string
	notice string
}

func NewLocationDialog(text, notice, share, cancel string) *LocationDialog {
	dlg := &LocationDialog{
		markup: &telebot.ReplyMarkup{ResizeKeyboard: true, OneTimeKeyboard: true},
		share:  &telebot.Btn{Text: share, Location: true},
		cancel: &telebot.Btn{Text: cancel},
		text:   text,
		notice: notice,
	}

	dlg.markup.Reply(telebot.Row{*dlg.share}, telebot.Row{*dlg.cancel})

	return dlg
}

func (p *LocationDialog) Open(ctx telebot.Context) error {
	return ctx.Send(p.text, p.markup)
}

func (p *LocationDialog) Close(ctx telebot.Context) error {
	return ctx.Send(p.notice, telebot.RemoveKeyboard)
}

func (p *LocationDialog) OnClose(ctx telebot.Context, fn telebot.HandlerFunc) {
	ctx.Bot().Handle(p.cancel, fn)
}
