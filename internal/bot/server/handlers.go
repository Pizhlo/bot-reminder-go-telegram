package server

import (
	"context"
	"errors"
	"strings"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/commands"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	logger "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func (s *Server) setupHandlers(ctx context.Context) {
	s.bot.Use(logger.Logging(ctx), middleware.AutoRespond())

	s.bot.Handle(commands.HelpCommand, func(telectx tele.Context) error {
		// if _, ok := s.fsm[telectx.Chat().ID]; !ok {
		// 	s.RegisterUserInFSM(telectx.Chat().ID)
		// }

		//return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		err := s.controller.HelpCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// геолокация
	s.bot.Handle(tele.OnLocation, func(telectx tele.Context) error {
		// err := s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		//s.RegisterUserInFSM(telectx.Chat().ID)

		logrus.Debugf("location")

		err := s.controller.AcceptTimezone(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// restricted: only known users

	restricted := s.bot.Group()
	restricted.Use(s.CheckUser(ctx), logger.Logging(ctx), middleware.AutoRespond())

	// часовой пояс
	restricted.Handle(&view.BtnTimezone, func(telectx tele.Context) error {
		logrus.Debugf("Timezone btn")
		err := s.controller.Timezone(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// изменить часовой пояс
	restricted.Handle(&view.BtnEditTimezone, func(telectx tele.Context) error {
		logrus.Debugf("Edit timezone btn")
		err := s.controller.RequestLocation(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// меню заметок
	restricted.Handle(&view.BtnNotes, func(telectx tele.Context) error {
		logrus.Debugf("Notes btn")
		s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].ListNote)
		// return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		err := s.controller.ListNotes(ctx, telectx)
		if err != nil {
			switch t := err.(type) {
			case *tele.Error:
				if strings.Contains(t.Description, "message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message") {
					return nil
				}
			default:
				s.HandleError(telectx, err)
				return err
			}
		}

		return nil
	})

	// меню напоминаний
	restricted.Handle(&view.BtnReminders, func(telectx tele.Context) error {
		logrus.Debugf("Reminders btn")
		s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].ListReminder)
		err := s.controller.ListReminders(ctx, telectx)
		if err != nil {
			switch t := err.(type) {
			case *tele.Error:
				if strings.Contains(t.Description, "message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message") {
					return nil
				}
			default:
				s.HandleError(telectx, err)
				return err
			}
		}

		return nil
	})

	// сообщить о баге
	restricted.Handle(&view.BtnBugReport, func(telectx tele.Context) error {
		logrus.Debugf("Bug report btn")
		s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].BugReportState)
		err := telectx.EditOrSend(messages.BugReportUserMessage, view.BackToMenuBtn())
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// кнопка чтобы скрыть клавиатуру у сработавшего напоминания
	restricted.Handle(&view.BtnCheckReminder, func(ctx tele.Context) error {
		// отправляем сообщение без клавиатуры
		err := ctx.Edit(ctx.Message().Text)
		if err != nil {
			s.HandleError(ctx, err)
			return err
		}

		return nil
	})

	// назад в меню
	restricted.Handle(&view.BtnBackToMenu, func(telectx tele.Context) error {
		logrus.Debugf("Menu btn")
		// s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].Start)
		// return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)

		s.fsm[telectx.Chat().ID].SetToDefault()

		err := s.controller.MenuCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// /start command
	restricted.Handle(commands.StartCommand, func(telectx tele.Context) error {
		// if _, ok := s.fsm[telectx.Chat().ID]; !ok {
		// 	s.RegisterUserInFSM(telectx.Chat().ID)
		// }

		//return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		err := s.controller.StartCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// /menu command
	restricted.Handle(commands.MenuCommand, func(telectx tele.Context) error {
		// if _, ok := s.fsm[telectx.Chat().ID]; !ok {
		// 	s.RegisterUserInFSM(telectx.Chat().ID)
		// }

		//return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		err := s.controller.MenuCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	restricted.Handle(tele.OnText, func(telectx tele.Context) error {
		logrus.Debugf("on text")
		//return s.controller.CreateNote(ctx, telectx)
		err := s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// notes

	// следующая страница заметок
	restricted.Handle(&view.BtnNextPgNotes, func(c tele.Context) error {
		err := s.controller.NextPageNotes(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// предыдущая страница заметок
	restricted.Handle(&view.BtnPrevPgNotes, func(c tele.Context) error {
		err := s.controller.PrevPageNotes(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// последняя страница заметок
	restricted.Handle(&view.BtnLastPgNotes, func(c tele.Context) error {
		err := s.controller.LastPageNotes(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// первая страница заметок
	restricted.Handle(&view.BtnFirstPgNotes, func(c tele.Context) error {
		err := s.controller.FirstPageNotes(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// поиск заметок по тексту
	restricted.Handle(&view.BtnSearchNotesByText, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].SearchNoteByText)

		err := c.EditOrSend(messages.SearchNotesByTextMessage, view.BackToMenuAndNotesBtn())
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// поиск заметок по дате
	restricted.Handle(&view.BtnSearchNotesByDate, func(c tele.Context) error {
		err := c.EditOrSend(messages.SearchNotesByDateChooseMessage, &tele.SendOptions{
			ReplyMarkup: view.SearchByDateBtn(),
			ParseMode:   "html",
		})
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// поиск заметок по одной дате
	restricted.Handle(&view.BtnSearchByOneDate, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].SearchNoteOneDate)

		s.controller.SetupNoteCalendar(ctx, c)
		s.controller.SetNoteCalendar(c.Chat().ID)

		err := s.controller.SearchNoteByOnedate(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			restricted.Handle(&btn, func(c tele.Context) error {
				s.fsm[c.Chat().ID].SetNext()
				err := s.fsm[c.Chat().ID].Handle(ctx, c)
				if err != nil {
					s.HandleError(c, err)
					return err
				}

				s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ListNote)
				s.controller.ResetCalendars(c.Chat().ID)
				return nil
			})
		}

		return nil
	})

	// поиск заметок по двум датам
	restricted.Handle(&view.BtnSearchByTwoDate, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].SearchNoteTwoDates)

		s.controller.SetupNoteCalendar(ctx, c)
		s.controller.SetNoteCalendar(c.Chat().ID)

		err := s.controller.SearchNoteByTwoDates(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		s.fsm[c.Chat().ID].SetNext()

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			restricted.Handle(&btn, func(c tele.Context) error {
				s.fsm[c.Chat().ID].SetNext()
				err := s.fsm[c.Chat().ID].Handle(ctx, c)
				if err != nil {
					if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
						return s.controller.SecondDateBeforeFirst(ctx, c)
					}

					if errors.Is(err, api_errors.ErrSecondDateFuture) {
						return s.controller.SecondDateInFuture(ctx, c)
					}

					if errors.Is(err, api_errors.ErrFirstDayFuture) {
						return s.controller.FirstDateInFuture(ctx, c)
					}

					s.HandleError(c, err)
					return err
				}

				return nil
			})
		}

		return nil
	})

	// удалить все заметки - спросить а точно ли
	restricted.Handle(&view.BtnDeleteAllNotes, func(c tele.Context) error {
		err := s.controller.ConfirmDeleteAllNotes(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// согласие удалить все заметки
	restricted.Handle(&controller.BtnDeleteAllNotes, func(c tele.Context) error {
		err := s.controller.DeleteAllNotes(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// отказ удалить все заметки
	restricted.Handle(&controller.BtnNotDeleteAllNotes, func(c tele.Context) error {
		err := c.Edit(messages.NotDeleteMessage, view.BackToMenuBtn())
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// reminders

	// навигация по страницам

	// предыдущая страница
	restricted.Handle(&view.BtnPrevPgReminders, func(c tele.Context) error {
		err := s.controller.PrevPageReminders(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// следующая страница
	restricted.Handle(&view.BtnNextPgReminders, func(c tele.Context) error {
		err := s.controller.NextPageReminders(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// первая страница
	restricted.Handle(&view.BtnFirstPgReminders, func(c tele.Context) error {
		err := s.controller.FirstPageReminders(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// последняя страница
	restricted.Handle(&view.BtnLastPgReminders, func(c tele.Context) error {
		err := s.controller.LastPageReminders(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// назад к выбору
	restricted.Handle(&view.BtnBackToReminderType, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderName)
		err := s.controller.ReminderName(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// удалить все напоминания
	restricted.Handle(&view.BtnDeleteAllReminders, func(c tele.Context) error {
		err := s.controller.ConfirmDeleteAllReminders(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// удалить все напоминания - подтверждение
	restricted.Handle(&controller.BtnDeleteAllReminders, func(c tele.Context) error {
		err := s.controller.DeleteAllReminders(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// название напоминания
	restricted.Handle(&view.BtnCreateReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderName)

		err := c.EditOrSend(messages.ReminderNameMessage, view.BackToMenuBtn())
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// reminder types

	// today
	restricted.Handle(&view.BtnToday, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.Today(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// tomorrow
	restricted.Handle(&view.BtnTomorrow, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.Tomorrow(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// everyday
	restricted.Handle(&view.BtnEveryDayReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.EverydayReminder(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// several times a day (once in N minutes, once in N hours)
	restricted.Handle(&view.BtnSeveralTimesDayReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].SeveralTimesDay)

		err := s.controller.SeveralTimesADayReminder(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// once in N minutes
	restricted.Handle(&view.BtnMinutesReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].MinutesDuration)

		err := s.controller.OnceInMinutes(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// once in N hours
	restricted.Handle(&view.BtnHoursReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].HoursDuration)

		err := s.controller.OnceInHours(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// times reminder
	restricted.Handle(&view.BtnTimesReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].Times)

		err := s.controller.TimesReminder(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// every week
	restricted.Handle(&view.BtnEveryWeekReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].EveryWeek)

		// err := s.controller.EveryWeek(ctx, c)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		// return nil

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// once in several days (e.g. once in 10 days)
	restricted.Handle(&view.BtnSeveralDaysReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].DaysDuration)

		err := s.controller.SeveralDays(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// once in month
	restricted.Handle(&view.BtnOnceMonthReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].Month)

		err := s.controller.Month(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		return nil
	})

	// every year
	restricted.Handle(&view.BtnOnceYear, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].Year)

		s.controller.SetupReminderCalendar(ctx, c)

		s.controller.SetReminderCalendar(c.Chat().ID)

		err := s.controller.Year(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			restricted.Handle(&btn, func(c tele.Context) error {
				// s.fsm[c.Chat().ID].SetNext()
				// err := s.fsm[c.Chat().ID].Handle(ctx, c)
				// if err != nil {
				// 	s.HandleError(telectx, err)
				// 	return err
				// }

				// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)
				// s.controller.ResetCalendars(c.Chat().ID)
				// return nil
				return s.fsm[c.Chat().ID].Handle(ctx, c)
			})
		}

		return nil
	})

	// date
	restricted.Handle(&view.BtnOnce, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].Date)

		s.controller.SetupReminderCalendar(ctx, c)
		s.controller.SetReminderCalendar(c.Chat().ID)

		err := s.controller.Date(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			restricted.Handle(&btn, func(c tele.Context) error {
				return s.fsm[c.Chat().ID].Handle(ctx, c)
			})
		}

		return nil
	})

	// calendar

	// prev month
	restricted.Handle(&view.BtnPrevMonth, func(c tele.Context) error {
		err := s.controller.PrevMonth(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			restricted.Handle(&btn, func(c tele.Context) error {
				err := s.fsm[c.Chat().ID].Handle(ctx, c)
				if err != nil {
					if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
						return s.controller.SecondDateBeforeFirst(ctx, c)
					}

					if errors.Is(err, api_errors.ErrSecondDateFuture) {
						return s.controller.SecondDateInFuture(ctx, c)
					}

					if errors.Is(err, api_errors.ErrFirstDayFuture) {
						return s.controller.FirstDateInFuture(ctx, c)
					}

					s.HandleError(c, err)
					return err
				}

				return nil

			})
		}

		return nil
	})

	// next month
	restricted.Handle(&view.BtnNextMonth, func(c tele.Context) error {
		err := s.controller.NextMonth(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			restricted.Handle(&btn, func(c tele.Context) error {
				err := s.fsm[c.Chat().ID].Handle(ctx, c)
				if err != nil {
					if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
						return s.controller.SecondDateBeforeFirst(ctx, c)
					}

					if errors.Is(err, api_errors.ErrSecondDateFuture) {
						return s.controller.SecondDateInFuture(ctx, c)
					}

					if errors.Is(err, api_errors.ErrFirstDayFuture) {
						return s.controller.FirstDateInFuture(ctx, c)
					}

					s.HandleError(c, err)
					return err
				}

				return nil
			})
		}

		return nil
	})

	// prev year
	restricted.Handle(&view.BtnPrevYear, func(c tele.Context) error {
		err := s.controller.PrevYear(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			restricted.Handle(&btn, func(c tele.Context) error {
				err := s.fsm[c.Chat().ID].Handle(ctx, c)
				if err != nil {
					if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
						return s.controller.SecondDateBeforeFirst(ctx, c)
					}

					if errors.Is(err, api_errors.ErrSecondDateFuture) {
						return s.controller.SecondDateInFuture(ctx, c)
					}

					if errors.Is(err, api_errors.ErrFirstDayFuture) {
						return s.controller.FirstDateInFuture(ctx, c)
					}

					s.HandleError(c, err)
					return err
				}

				return nil
			})
		}

		return nil
	})

	// next year
	restricted.Handle(&view.BtnNextYear, func(c tele.Context) error {
		err := s.controller.NextYear(ctx, c)
		if err != nil {
			s.HandleError(c, err)
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			restricted.Handle(&btn, func(c tele.Context) error {
				err := s.fsm[c.Chat().ID].Handle(ctx, c)
				if err != nil {
					if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
						return s.controller.SecondDateBeforeFirst(ctx, c)
					}

					if errors.Is(err, api_errors.ErrSecondDateFuture) {
						return s.controller.SecondDateInFuture(ctx, c)
					}

					if errors.Is(err, api_errors.ErrFirstDayFuture) {
						return s.controller.FirstDateInFuture(ctx, c)
					}

					s.HandleError(c, err)
					return err
				}

				return nil
			})
		}

		return nil
	})

	// week days

	// Monday
	restricted.Handle(&view.MondayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		// return nil

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Tuesday
	restricted.Handle(&view.TuesdayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Wednesday
	restricted.Handle(&view.WednesdayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Thursday
	restricted.Handle(&view.ThursdayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Friday
	restricted.Handle(&view.FridayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Saturday
	restricted.Handle(&view.SaturdayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Sunday
	restricted.Handle(&view.SundayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})
}
