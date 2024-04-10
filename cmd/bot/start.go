package bot

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/config"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	note_srv "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/reminder"
	user_srv "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	tz_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	note_db "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/note"
	reminder_db "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/reminder"
	tz_db "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/timezone"
	user_db "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/user"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func Start(confName, path string) {
	logrus.Info("starting")
	defer logrus.Info("stopped")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.LoadConfig(confName, path)
	if err != nil {
		logrus.Fatalf("unable to load config: %v", err)
	}

	switch conf.LogLvl {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "trace":
		logrus.SetLevel(logrus.TraceLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.PanicLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// db
	dbAddr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName)

	userRepo, err := user_db.New(dbAddr)
	if err != nil {
		logrus.Fatalf("cannot create user repo: %v", err)
	}

	tzRepo, err := tz_db.New(dbAddr)
	if err != nil {
		logrus.Fatalf("cannot create user timezone repo: %v", err)
	}

	noteRepo, err := note_db.New(dbAddr)
	if err != nil {
		logrus.Fatalf("cannot create note repo: %v", err)
	}

	reminderRepo, err := reminder_db.New(dbAddr)
	if err != nil {
		logrus.Fatalf("cannot create reminder repo: %v", err)
	}

	// bot
	bot, err := tele.NewBot(tele.Settings{
		URL:       conf.BotURL,
		Token:     conf.Token,
		Poller:    &tele.LongPoller{Timeout: conf.Timeout},
		ParseMode: "html",
	})
	if err != nil {
		logrus.Fatalf("cannot create a bot: %v", err)
	}

	logrus.Info("successfully created bot")

	// cache
	tz := tz_cache.New()

	// services
	userSrv := user_srv.New(ctx, userRepo, tz, tzRepo)
	noteSrv := note_srv.New(noteRepo)
	reminderSrv := reminder.New(reminderRepo)

	controller := controller.New(userSrv, noteSrv, bot, reminderSrv)

	// server
	server := server.New(bot, controller)

	logrus.Info("starting server...")

	server.Start(ctx)

	go func() {
		_, msgErr := bot.Send(&tele.Chat{ID: conf.ChannelID}, messages.StartBotMessage)
		if msgErr != nil {
			logrus.Errorf("Error while sending message 'Бот запущен': %v\n", msgErr)
		}

		bot.Start()
	}()

	notifyCtx, notify := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer notify()

	<-notifyCtx.Done()
	logrus.Info("shutdown")

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		_, msgErr := bot.Send(&tele.Chat{ID: conf.ChannelID}, messages.ShutDownMessage)
		if msgErr != nil {
			logrus.Errorf("Error while sending message 'Бот выключается': %v\n", msgErr)
		}
		logrus.Info("gently shutdown")

		server.Shutdown(ctx)

		bot.Stop()

		userRepo.Close()
		noteRepo.Close()
		tzRepo.Close()
		reminderRepo.Close()
	}(&wg)

	wg.Wait()

	notify()

}
