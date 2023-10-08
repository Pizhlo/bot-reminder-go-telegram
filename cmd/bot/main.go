package main

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/calendar"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/config"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/reminder"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/user"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	setup()
}

func setup() {
	conf, err := config.LoadConfig(`../..`)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	logger := logger.New()
	srv := server.New(setupDB(conf.DBAddress))
	calendar := calendar.New()

	controller := controller.New(srv, calendar, logger)

	logger.Log().Msg(`successfully loaded app`)
	fmt.Println(controller)
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
