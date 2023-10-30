package controller

import (
	default_handler "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/handler/default"
	note_handler "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/handler/note"
	tz_handler "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/handler/timezone"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	pkg_err "github.com/pkg/errors"
	tele "gopkg.in/telebot.v3"
)

type Controller struct {
	commandHandlerMap map[int64]commandHandler
	textHanderMap     map[int64]textHander
	bot               *tele.Bot
	logger            *logger.Logger
	srv               *server.Server
}

type commandHandler interface {
	Handle(ctx tele.Context) error
}

type textHander interface {
	Handle(ctx tele.Context) error
}

const (
	startCmd             = `/start`
	addReminderCmd       = `/add_reminder`
	searchNotesByTextCmd = `/notes_search_text`
	notesCmd             = `/notes`
)

func New(bot *tele.Bot, logger *logger.Logger, srv *server.Server) *Controller {
	handlerMap := make(map[int64]commandHandler)
	textHander := make(map[int64]textHander)

	return &Controller{handlerMap, textHander, bot, logger, srv}
}
func (c *Controller) SetupBot() error {
	//c.srv.UserCacheEditor.SaveUser(1, 297850814) // для целей тестирования
	//c.srv.UserCacheEditor.SaveUser(1, 297850814) // для целей тестирования

	// commands

	c.bot.Handle(startCmd, func(ctx tele.Context) error {
		c.commandHandlerMap[ctx.Chat().ID] = default_handler.New(c.srv)
		c.textHanderMap[ctx.Chat().ID] = note_handler.NewSaveNoteHandler(c.srv) // поведение на текст по умолчанию: сохранить заметку
		return c.commandHandlerMap[ctx.Chat().ID].Handle(ctx)
	})

	c.bot.Handle(addReminderCmd, func(ctx tele.Context) error {
		return nil
	})

	c.bot.Handle(searchNotesByTextCmd, func(ctx tele.Context) error {
		return nil
	})

	c.bot.Handle(notesCmd, func(ctx tele.Context) error {
		c.commandHandlerMap[ctx.Chat().ID] = note_handler.NewNotesHandler(c.srv)
		return c.commandHandlerMap[ctx.Chat().ID].Handle(ctx)
	})

	// types

	c.bot.Handle(tele.OnText, func(ctx tele.Context) error {
		return c.textHanderMap[ctx.Chat().ID].Handle(ctx)
	})

	// buttons

	c.bot.Handle(&default_handler.RejectBtn, func(ctx tele.Context) error {
		return ctx.Send(`text`, tele.RemoveKeyboard)
	})

	c.bot.Handle(&default_handler.LocationBtn, func(ctx tele.Context) error {
		tz := c.parseTimezone(ctx.Chat().ID, ctx.Query().Location)
		c.commandHandlerMap[ctx.Chat().ID] = tz_handler.New(tz, c.srv)

		return c.commandHandlerMap[ctx.Chat().ID].Handle(ctx)
	})

	if err := c.sendStartupMsg(); err != nil {
		return pkg_err.Wrap(err, `unable to send startup message`)
	}

	c.startBot()

	return nil
}

func (c *Controller) sendStartupMsg() error {
	_, err := c.bot.Send(tele.Recipient(&tele.Chat{ID: -1001890622926}), "#запуск\nБот запущен")
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) startBot() {
	c.logger.Info().Msg(`successfully loaded app`)
	c.bot.Start()
	c.logger.Info().Msg(`successfully loaded app`)
	c.bot.Start()
}
