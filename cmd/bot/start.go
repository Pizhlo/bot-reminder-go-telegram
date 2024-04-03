package bot

import (
	"context"
	"fmt"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/config"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
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
	tele "gopkg.in/telebot.v3"
)

func Start(confName, path string) {
	logger := logger.New()

	err := func(ctx context.Context) error {
		logger.Info("starting")
		defer logger.Info("stopped")

		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		conf, err := config.LoadConfig(confName, path)
		if err != nil {
			return fmt.Errorf("unable to load config: %w", err)
		}

		_ = func() context.Context {
			c, cancel := context.WithTimeout(ctx, conf.Timeout)
			defer cancel()
			return c
		}

		// tzf, err := tzf.NewDefaultFinder()
		// if err != nil {
		// 	return fmt.Errorf("cannot initialize a time zone finder: %w", err)
		// }

		// db
		dbAddr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", conf.DBUser, conf.DBPass, conf.DBHost, conf.DBPort, conf.DBName)

		userRepo, err := user_db.New(dbAddr)
		if err != nil {
			return fmt.Errorf("cannot create user repo: %w", err)
		}

		tzRepo, err := tz_db.New(dbAddr)
		if err != nil {
			return fmt.Errorf("cannot create user timezone repo: %w", err)
		}

		noteRepo, err := note_db.New(dbAddr)
		if err != nil {
			return fmt.Errorf("cannot create note repo: %w", err)
		}

		reminderRepo, err := reminder_db.New(dbAddr)
		if err != nil {
			return fmt.Errorf("cannot create reminder repo: %w", err)
		}

		// bot
		bot, err := tele.NewBot(tele.Settings{
			URL:       conf.BotURL,
			Token:     conf.Token,
			Poller:    &tele.LongPoller{Timeout: conf.Timeout},
			ParseMode: "html",
		})
		if err != nil {
			return fmt.Errorf("cannot create a bot: %w", err)
		}

		// cache
		tz := tz_cache.New()

		// services
		userSrv := user_srv.New(ctx, userRepo, tz, tzRepo)
		noteSrv := note_srv.New(noteRepo)
		reminderSrv := reminder.New(reminderRepo)

		controller := controller.New(userSrv, noteSrv, bot, reminderSrv)

		logger.Debug("successfully created bot")

		// server
		server := server.New(bot, controller)

		logger.Debug("starting server...")

		server.Start(ctx)

		go func() {
			//defer cancel()
			_, msgErr := bot.Send(&tele.Chat{ID: conf.ChannelID}, messages.StartBotMessage)
			if msgErr != nil {
				logger.Errorf("Error while sending message 'Бот запущен': %v\n", msgErr)
			}

			bot.Start()
		}()

		notifyCtx, notify := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
		defer notify()

		go func() {
			defer cancel()

			<-notifyCtx.Done()

			closer := make(chan struct{})

			go func() {
				bot.Stop()

				closer <- struct{}{}
			}()

			shutdownCtx, shutdown := context.WithTimeout(context.Background(), conf.Timeout)
			defer shutdown()

			runtime.Gosched()

			select {
			case <-closer:
				_, msgErr := bot.Send(&tele.Chat{ID: conf.ChannelID}, messages.ShutDownMessage)
				if msgErr != nil {
					logger.Errorf("Error while sending message 'Бот запущен': %v\n", msgErr)
				}
				logger.Info("gently shutdown")

				server.Shutdown(ctx)

			case <-shutdownCtx.Done():
				logger.Error("forcing shutdown")
			}
		}()

		logger.Info("started")

		<-ctx.Done()

		logger.Info("shutting down")

		cancel()

		return nil
	}(context.Background())

	if err != nil {
		logger.Fatalf("unable to start: %v\n", err)
	}
}
