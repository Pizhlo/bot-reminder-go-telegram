package bot

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"os"
// 	"os/signal"
// 	"runtime"
// 	"syscall"

// 	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/config"
// 	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
// 	"github.com/ringsaturn/tzf"
// 	"github.com/sirupsen/logrus"
// 	tele "gopkg.in/telebot.v3"
// )

// const _BotPrefix = "TELENOTE"

// const helpText = `
// /start - зарегистрироваться;
// /help - вывести это сообщение;
// /note мандарины среди нас - добавить заметку с текстом "мандарины среди нас";
// /notes - вывести все заметки списком;
// /notes telegram bot - вывести все заметки, например, которые содержать слова: telegram и/или bot.
// `

// const locDlgMsg = "Ну, что ж!\nДля начала мне надо узнать ваш часовой пояс, а потом повторите свою команду."
// const LocDlgNotice = "Хорошо, как-нибудь в следующий раз!"

// func Start(confName, path string) {
// 	logger := logrus.New()

// 	err := func(ctx context.Context) error {
// 		logger.Info("starting")
// 		defer logger.Info("stopped")

// 		ctx, cancel := context.WithCancel(ctx)
// 		defer cancel()

// 		conf, err := config.LoadConfig(confName, path)
// 		if err != nil {
// 			return fmt.Errorf("unable to load config: %w", err)
// 		}

// 		baseContext := func() context.Context {
// 			c, _ := context.WithTimeout(ctx, conf.Timeout)
// 			return c
// 		}

// 		tzf, err := tzf.NewDefaultFinder()
// 		if err != nil {
// 			return fmt.Errorf("cannot initialize a time zone finder: %w", err)
// 		}

// 		var users = user_cache.New()

// 		logger.Info("created user inmemory store")

// 		notes, err := postgres_note.New(conf.DBAddress)
// 		if err != nil {
// 			return fmt.Errorf("unable to create notes repo: %w", err)
// 		}

// 		logger.Info("created note postgres store")

// 		userRepo, err := postgres_user.New(conf.DBAddress)
// 		if err != nil {
// 			return fmt.Errorf("unable to user repo: %w", err)
// 		}

// 		logger.Info("created user postgres store")

// 		locDlg := NewLocationDialog(locDlgMsg, LocDlgNotice, "Ну, ладно, разрешаю следить за собой ...", "Помогите, зрения лишают!")

// 		noteCtrl := controller.NewNote(baseContext, noting, accounting)
// 		noteDelCtrl := controller.NewNoteDelete(baseContext)

// 		dynamicCtrl := controller.NewRegex(noteCtrl.Handle)
// 		dynamicCtrl.Handle(`^\/dn\d+$`, noteDelCtrl.Handle)

// 		bot, err := tele.NewBot(tele.Settings{
// 			Token:  conf.Token,
// 			Poller: &tele.LongPoller{Timeout: conf.Timeout},
// 		})

// 		if err != nil {
// 			return fmt.Errorf("cannot create a bot: %w", err)
// 		}

// 		locDlg.OnClose(bot.NewContext(tele.Update{}), func(ctx tele.Context) error {
// 			logger.Info("handling", "command", "location dlg close btn", "sender id", ctx.Sender().ID)

// 			return locDlg.Close(ctx)
// 		})

// 		bot.Handle("/start", func(telctx tele.Context) error {
// 			l := logger.With(
// 				"command", "start",
// 				"sender id", telctx.Sender().ID,
// 			)

// 			ctx, cancel := context.WithCancel(baseContext())
// 			defer cancel()

// 			_, err := accounting.FindUserByTelegramID(ctx, int(telctx.Sender().ID))
// 			if err != nil {
// 				err = errors.Unwrap(err)
// 				if errors.Is(err, user.ErrNotFound) {
// 					acc, err := accounting.AddUser(ctx, int(telctx.Sender().ID))
// 					if err != nil {
// 						l.Error("error adding a user", "error", err)

// 						return fmt.Errorf("cannot add a user on start: %w", err)
// 					}

// 					logger.Info("registered a user", "user id", acc.ID, "telegram id", int(telctx.Sender().ID))

// 					err = telctx.Send(fmt.Sprintf("Привет, %s!", telctx.Sender().FirstName))
// 					if err != nil {
// 						return err
// 					}

// 					err = locDlg.Open(telctx)
// 					if err != nil {
// 						return err
// 					}

// 					return nil
// 				}

// 				l.Error("error finding a user", "error", err)

// 				return fmt.Errorf("cannot find a user on start: %w", err)
// 			}

// 			return telctx.Send("Чего желаешь в этот раз?")
// 		})

// 		bot.Handle("/help", func(ctx tele.Context) error {
// 			logger.Info("handling", "command", "help", "chat id", ctx.Chat().ID)

// 			return ctx.Send(fmt.Sprintf("Привет, %s!\n%s", ctx.Sender().FirstName, helpText))
// 		})

// 		bot.Handle(tele.OnLocation, func(telctx tele.Context) error {
// 			l := logger.With(
// 				"command", "start",
// 				"sender id", telctx.Sender().ID,
// 			)

// 			ctx, cancel := context.WithCancel(baseContext())
// 			defer cancel()

// 			lon, lat := telctx.Message().Location.Lng, telctx.Message().Location.Lat

// 			locName := tzf.GetTimezoneName(float64(lon), float64(lat))

// 			tz := model_user.Timezone{
// 				Name: locName,
// 				Lon:  float64(lon),
// 				Lat:  float64(lat),
// 			}

// 			l.Info("figured out a time zone", "timezone", tz)

// 			acc, err := accounting.FindUserByTelegramID(ctx, int(telctx.Sender().ID))
// 			if err != nil {
// 				l.Error("error updating a timezone", "error", err)

// 				return fmt.Errorf("cannot update a timezone on location: %w", err)
// 			}

// 			accounting.UpdateUser(ctx, acc.ID, &model_user.User{Timezone: tz})
// 			if err != nil {
// 				l.Error("error updating a timezone", "error", err)

// 				return fmt.Errorf("cannot update a timezone on location: %w", err)
// 			}

// 			return telctx.Send(fmt.Sprintf("Спасибо! Хорошего дня в %s!", locName), tele.RemoveKeyboard)
// 		})

// 		restricted := bot.Group()
// 		restricted.Use(NewVerifyTimezone(accounting, locDlg))

// 		restricted.Handle("/note", noteCtrl.Handle)

// 		restricted.Handle("/notes", func(telctx tele.Context) error {
// 			tid := tgid{telctx}.Get()

// 			l := logger.With(
// 				"command", "/notes",
// 				"sender id", tid,
// 			)

// 			l.Info("handling", "args", telctx.Args())

// 			ctx, cancel := context.WithCancel(baseContext())
// 			defer cancel()

// 			u, err := accounting.FindUserByTelegramID(ctx, tid)
// 			if err != nil {
// 				return err
// 			}

// 			var view notelist.NoteList = notelist.NewTelegramMarkdownView(telctx, u.ID, conf.TelegramMaxMsgSize)

// 			var presenter pnl.NoteList = pnl.NewStandard(noting, view)

// 			return presenter.UpdateNoteList(ctx)
// 		})

// 		// restricted.Handle(tele.OnText, func(telectx tele.Context) error {

// 		// })

// 		go func() {
// 			defer cancel()

// 			bot.Start()
// 		}()

// 		notifyCtx, notify := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
// 		defer notify()

// 		go func() {
// 			defer cancel()

// 			<-notifyCtx.Done()

// 			closer := make(chan struct{})

// 			go func() {
// 				bot.Stop()

// 				closer <- struct{}{}
// 			}()

// 			shutdownCtx, shutdown := context.WithTimeout(context.Background(), conf.Timeout)
// 			defer shutdown()

// 			runtime.Gosched()

// 			select {
// 			case <-closer:
// 				logger.Info("gently shutdown")

// 			case <-shutdownCtx.Done():
// 				logger.Error("forcing shutdown")
// 			}
// 		}()

// 		logger.Info("started")

// 		<-ctx.Done()

// 		logger.Info("shutting down")

// 		cancel()

// 		return nil
// 	}(context.Background())
// 	if err != nil {
// 		logger.Error("error running a bot", "error", err)
// 		os.Exit(1)
// 	}
// }

// // func setupDB(dbAddr string) (*note.NoteRepo, *reminder.ReminderRepo, *user.UserRepo) {
// // 	conn, err := pgxpool.Connect(context.Background(), dbAddr)
// // 	if err != nil {
// // 		log.Fatal().Err(err).Msg("failed to connect db")
// // 	}

// // 	noteEditor := note.New(conn)
// // 	reminderEditor := reminder.New(conn)
// // 	userEditor := user.New(conn)

// // 	return noteEditor, reminderEditor, userEditor
// // }

// // func setupCache() (*timezone_cache.TimezoneCache, *user_cache.UserCache) {
// // 	tzCache := timezone_cache.New()
// // 	userCache := user_cache.New()

// // 	return tzCache, userCache
// // }

// func setup() {
// 	// conf
// 	// conf, err := config.LoadConfig(`../..`)
// 	// if err != nil {
// 	// 	log.Fatal().Err(err).Msg("failed to load config")
// 	// }

// 	// // bot
// 	// pref := tele.Settings{
// 	// 	Token:  conf.Token,
// 	// 	Poller: &tele.LongPoller{Timeout: 10 * time.Second},
// 	// }
// 	// b, err := tele.NewBot(pref)
// 	// if err != nil {
// 	// 	log.Fatal().Err(err).Msg("failed to create bot")
// 	// }

// 	// server
// 	// logger := logger.New()
// 	// noteEditor, remiderEditor, userEditor := setupDB(conf.DBAddress)

// 	// logger.Info().Msg(`successfully connected db`)

// 	// tzCache, userCache := setupCache()

// 	// srv := server.New(noteEditor, remiderEditor, userEditor, tzCache, userCache)

// 	// // contoller
// 	// controller := controller.New(b, logger, srv)

// 	// if err := controller.SetupBot(); err != nil {
// 	// 	logger.Fatal().Err(err).Msg("failed to start bot")
// 	// }
// }
