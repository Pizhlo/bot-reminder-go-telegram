package bot

import (
	"context"
	"fmt"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/config"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/server"
	note_srv "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	user_srv "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	tz_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	note_db "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/note"
	tz_db "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/timezone"
	user_db "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/postgres/user"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

func Start(confName, path string) {
	logger := logrus.New()

	err := func(ctx context.Context) error {
		logger.Info("starting")
		defer logger.Info("stopped")

		ctx, cancel := context.WithCancel(ctx)
		//defer cancel()

		conf, err := config.LoadConfig(confName, path)
		if err != nil {
			return fmt.Errorf("unable to load config: %w", err)
		}

		baseContext := func() context.Context {
			c, _ := context.WithTimeout(ctx, conf.Timeout)
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
			return fmt.Errorf("cannot create user timezone: %w", err)
		}

		noteRepo, err := note_db.New(dbAddr)
		if err != nil {
			return fmt.Errorf("cannot create note timezone: %w", err)
		}

		// cache
		//u_cache := user_cache.New()
		tz := tz_cache.New()

		// services
		userSrv := user_srv.New(userRepo, tz, tzRepo)
		noteSrv := note_srv.New(noteRepo)

		controller := controller.NewMyController(userSrv, noteSrv)

		// bot
		bot, err := tele.NewBot(tele.Settings{
			Token:  conf.Token,
			Poller: &tele.LongPoller{Timeout: conf.Timeout},
		})
		if err != nil {
			return fmt.Errorf("cannot create a bot: %w", err)
		}

		logger.Debug("successfully created bot")

		// server
		server := server.New(bot, controller)

		logger.Debug("starting server...")

		server.Start(baseContext())

		go func() {
			//defer cancel()

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
				logger.Info("gently shutdown")

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
