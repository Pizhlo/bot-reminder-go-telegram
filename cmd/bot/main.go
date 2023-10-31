package main

import (
	"context"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/config"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	timezone_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	user_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/user"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
	tele "gopkg.in/telebot.v3"
)

func main() {
	setup()
}

func setup() {
	// conf
	conf, err := config.LoadConfig(`../..`)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	// bot
	pref := tele.Settings{
		Token:  conf.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create bot")
	}

	// server
	logger := logger.New()
	noteEditor, remiderEditor, userEditor := setupDB(conf.DBAddress)

	logger.Info().Msg(`successfully connected db`)

	tzCache, userCache := setupCache()

	srv := server.New(noteEditor, remiderEditor, userEditor, tzCache, userCache)

	// contoller
	controller := controller.New(b, logger, srv)

	if err := controller.SetupBot(); err != nil {
		logger.Fatal().Err(err).Msg("failed to start bot")
	}
}

func setupDB(dbAddr string) (*note.NoteRepo, *reminder.ReminderRepo, *user.UserRepo) {
	conn, err := pgxpool.Connect(context.Background(), dbAddr)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect db")
	}

	noteEditor := note.New(conn)
	reminderEditor := reminder.New(conn)
	userEditor := user.New(conn)

	return noteEditor, reminderEditor, userEditor
}

func setupCache() (*timezone_cache.TimezoneCache, *user_cache.UserCache) {
	tzCache := timezone_cache.New()
	userCache := user_cache.New()

	return tzCache, userCache
}
