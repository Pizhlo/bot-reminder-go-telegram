package server

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/commands"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func (s *Server) setupBot(ctx context.Context) {
	s.bot.Use(logger.Logging(ctx, s.logger))
	s.bot.Use(middleware.AutoRespond())

	// геолокация
	s.bot.Handle(tele.OnLocation, func(telectx tele.Context) error {
		err := s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err, s.fsm[telectx.Chat().ID].Name())
			return err
		}

		return nil
	})

	// главное меню
	s.bot.Handle(&view.BtnProfile, func(telectx tele.Context) error {
		s.logger.Debugf("Profile btn")
		err := s.controller.Profile(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err, s.fsm[telectx.Chat().ID].Name())
			return err
		}

		return nil
	})

	// меню настроек
	s.bot.Handle(&view.BtnSettings, func(telectx tele.Context) error {
		s.logger.Debugf("Settings btn")
		err := s.controller.Settings(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err, s.fsm[telectx.Chat().ID].Name())
			return err
		}

		return nil
	})

	// меню заметок
	s.bot.Handle(&view.BtnNotes, func(telectx tele.Context) error {
		s.logger.Debugf("Notes btn")
		s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].ListNote)
		// return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		err := s.controller.ListNotes(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err, s.fsm[telectx.Chat().ID].Name())
			return err
		}

		return nil
	})

	// меню напоминаний
	s.bot.Handle(&view.BtnReminders, func(telectx tele.Context) error {
		s.logger.Debugf("Reminders btn")
		err := s.controller.ListReminders(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err, s.fsm[telectx.Chat().ID].Name())
			return err
		}

		return nil
	})

	// назад в меню
	s.bot.Handle(&view.BtnBackToMenu, func(telectx tele.Context) error {
		s.logger.Debugf("Menu btn")
		// s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].Start)
		// return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)

		s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].DefaultState)

		err := s.controller.StartCmd(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err, s.fsm[telectx.Chat().ID].Name())
			return err
		}

		return nil
	})

	// restricted: only known users

	restricted := s.bot.Group()
	restricted.Use(s.CheckUser(ctx), logger.Logging(ctx, s.logger), middleware.AutoRespond())

	restricted.Handle(commands.StartCommand, func(telectx tele.Context) error {
		if _, ok := s.fsm[telectx.Chat().ID]; !ok {
			s.RegisterUser(telectx.Chat().ID, false)
		}

		//return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		err := s.controller.StartCmd(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err, s.fsm[telectx.Chat().ID].Name())
			return err
		}

		return nil
	})

	restricted.Handle(tele.OnText, func(telectx tele.Context) error {
		s.logger.Debugf("on text")
		//return s.controller.CreateNote(ctx, telectx)
		err := s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		if err != nil {
			s.controller.HandleError(telectx, err, s.fsm[telectx.Chat().ID].Name())
			return err
		}

		return nil
	})

	// notes

	// следующая страница заметок
	s.bot.Handle(&view.BtnNextPgNotes, func(c tele.Context) error {
		err := s.controller.NextPageNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// предыдущая страница заметок
	s.bot.Handle(&view.BtnPrevPgNotes, func(c tele.Context) error {
		err := s.controller.PrevPageNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// последняя страница заметок
	s.bot.Handle(&view.BtnLastPgNotes, func(c tele.Context) error {
		err := s.controller.LastPageNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// первая страница заметок
	s.bot.Handle(&view.BtnFirstPgNotes, func(c tele.Context) error {
		err := s.controller.FirstPageNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// поиск заметок по тексту
	s.bot.Handle(&view.BtnSearchNotesByText, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].SearchNoteByText)

		err := c.EditOrSend(messages.SearchNotesByTextMessage, view.BackToMenuBtn())
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// удалить все заметки - спросить а точно ли
	s.bot.Handle(&view.BtnDeleteAllNotes, func(c tele.Context) error {
		err := s.controller.ConfirmDeleteAllNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// согласие удалить все заметки
	s.bot.Handle(&controller.BtnDeleteAllNotes, func(c tele.Context) error {
		err := s.controller.DeleteAllNotes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// отказ удалить все заметки
	s.bot.Handle(&controller.BtnNotDeleteAllNotes, func(c tele.Context) error {
		err := c.Edit(messages.NotDeleteMessage, view.BackToMenuBtn())
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// reminders

	// навигация по страницам

	// предыдущая страница
	s.bot.Handle(&view.BtnPrevPgReminders, func(c tele.Context) error {
		err := s.controller.PrevPageReminders(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// следующая страница
	s.bot.Handle(&view.BtnNextPgReminders, func(c tele.Context) error {
		err := s.controller.NextPageReminders(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// первая страница
	s.bot.Handle(&view.BtnFirstPgReminders, func(c tele.Context) error {
		err := s.controller.FirstPageReminders(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// последняя страница
	s.bot.Handle(&view.BtnLastPgReminders, func(c tele.Context) error {
		err := s.controller.LastPageReminders(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// назад к выбору
	s.bot.Handle(&view.BtnBackToReminderType, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderName)
		err := s.controller.ReminderName(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// удалить сработавшее напоминание
	s.bot.Handle(&view.BtnDeleteReminder, func(ctx tele.Context) error {
		uniq := ctx.Callback().Unique
		return ctx.EditOrSend(fmt.Sprintf("Удалено: %s", uniq))
	})

	// удалить все напоминания - подтверждение
	s.bot.Handle(&view.BtnDeleteAllReminders, func(c tele.Context) error {
		err := s.controller.ConfirmDeleteAllReminders(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// удалить все напоминания
	s.bot.Handle(&controller.BtnDeleteAllReminders, func(c tele.Context) error {
		err := s.controller.DeleteAllReminders(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// название напоминания
	s.bot.Handle(&view.BtnCreateReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderName)

		err := c.EditOrSend(messages.ReminderNameMessage, view.BackToMenuBtn())
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// reminder types

	// everyday
	s.bot.Handle(&view.BtnEveryDayReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.EverydayReminder(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// several times a day (once in N minutes, once in N hours)
	s.bot.Handle(&view.BtnSeveralTimesDayReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].SeveralTimesDay)

		err := s.controller.SeveralTimesADayReminder(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// once in N minutes
	s.bot.Handle(&view.BtnMinutesReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].MinutesDuration)

		err := s.controller.OnceInMinutes(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// once in N hours
	s.bot.Handle(&view.BtnHoursReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].HoursDuration)

		err := s.controller.OnceInHours(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// every week
	s.bot.Handle(&view.BtnEveryWeekReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].EveryWeek)

		err := s.controller.EveryWeek(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// once in several days (e.g. once in 10 days)
	s.bot.Handle(&view.BtnSeveralDaysReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].DaysDuration)

		err := s.controller.SeveralDays(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// once in month
	s.bot.Handle(&view.BtnOnceMonthReminder, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].Month)

		err := s.controller.Month(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// every year
	s.bot.Handle(&view.BtnOnceYear, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].Year)

		err := s.controller.Year(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				err := s.controller.ProcessDate(ctx, c)
				if err != nil {
					s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
					return err
				}

				s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)
				return nil
			})
		}

		return nil
	})

	// calendar

	// prev month
	s.bot.Handle(&view.BtnPrevMonth, func(c tele.Context) error {
		err := s.controller.PrevMonth(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				err := s.controller.ProcessDate(ctx, c)
				if err != nil {
					s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
					return err
				}

				s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)
				return nil
			})
		}

		return nil
	})

	// next month
	s.bot.Handle(&view.BtnNextMonth, func(c tele.Context) error {
		err := s.controller.NextMonth(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				err := s.controller.ProcessDate(ctx, c)
				if err != nil {
					s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
					return err
				}

				s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)
				return nil
			})
		}

		return nil
	})

	// prev year
	s.bot.Handle(&view.BtnPrevYear, func(c tele.Context) error {
		err := s.controller.PrevYear(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				err := s.controller.ProcessDate(ctx, c)
				if err != nil {
					s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
					return err
				}

				s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)
				return nil
			})
		}

		return nil
	})

	// next year
	s.bot.Handle(&view.BtnNextYear, func(c tele.Context) error {
		err := s.controller.NextYear(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				err := s.controller.ProcessDate(ctx, c)
				if err != nil {
					s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
					return err
				}

				s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)
				return nil
			})
		}

		return nil
	})

	// week days

	// Monday
	s.bot.Handle(&view.MondayBtn, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.WeekDay(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// Tuesday
	s.bot.Handle(&view.TuesdayBtn, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.WeekDay(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// Wednesday
	s.bot.Handle(&view.WednesdayBtn, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.WeekDay(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// Thursday
	s.bot.Handle(&view.ThursdayBtn, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.WeekDay(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// Friday
	s.bot.Handle(&view.FridayBtn, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.WeekDay(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// Saturday
	s.bot.Handle(&view.SaturdayBtn, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.WeekDay(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// Sunday
	s.bot.Handle(&view.SundayBtn, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.WeekDay(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})
}
