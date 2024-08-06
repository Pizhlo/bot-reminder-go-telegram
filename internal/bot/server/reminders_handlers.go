package server

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

func (s *Server) setupRemindersHandlers(ctx context.Context, restricted *tele.Group) {
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
}
