package server

import (
	"context"
	"errors"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/commands"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

func (s *Server) setupBot(ctx context.Context) {
	s.bot.Use(logger.Logging(ctx, s.logger), s.CheckUser(ctx), middleware.AutoRespond())

	// геолокация
	s.bot.Handle(tele.OnLocation, func(telectx tele.Context) error {
		// err := s.fsm[telectx.Chat().ID].Handle(ctx, telectx)
		// if err != nil {
		// 	s.HandleError(telectx, err)
		// 	return err
		// }

		//s.RegisterUserInFSM(telectx.Chat().ID)

		err := s.controller.AcceptTimezone(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// часовой пояс
	s.bot.Handle(&view.BtnTimezone, func(telectx tele.Context) error {
		s.logger.Debugf("Timezone btn")
		err := s.controller.Timezone(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// изменить часовой пояс
	s.bot.Handle(&view.BtnEditTimezone, func(telectx tele.Context) error {
		s.logger.Debugf("Edit timezone btn")
		err := s.controller.RequestLocation(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
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
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// меню напоминаний
	s.bot.Handle(&view.BtnReminders, func(telectx tele.Context) error {
		s.logger.Debugf("Reminders btn")
		s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].ListReminder)
		err := s.controller.ListReminders(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// назад в меню
	s.bot.Handle(&view.BtnBackToMenu, func(telectx tele.Context) error {
		s.logger.Debugf("Menu btn")
		// s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].Start)
		// return s.fsm[telectx.Chat().ID].Handle(ctx, telectx)

		s.fsm[telectx.Chat().ID].SetToDefault()

		err := s.controller.StartCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// кнопка чтобы скрыть клавиатуру у сработавшего напоминания
	s.bot.Handle(&view.BtnCheckReminder, func(ctx tele.Context) error {
		// отправляем сообщение без клавиатуры
		err := ctx.Edit(ctx.Message().Text)
		if err != nil {
			s.HandleError(ctx, err)
			return err
		}

		return nil
	})

	// restricted: only known users

	restricted := s.bot.Group()
	restricted.Use(s.CheckUser(ctx), logger.Logging(ctx, s.logger), middleware.AutoRespond())

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

	restricted.Handle(commands.HelpCommand, func(telectx tele.Context) error {
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

	restricted.Handle(tele.OnText, func(telectx tele.Context) error {
		s.logger.Debugf("on text")
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

		err := c.EditOrSend(messages.SearchNotesByTextMessage, view.BackToMenuAndNotesBtn())
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// поиск заметок по дате
	s.bot.Handle(&view.BtnSearchNotesByDate, func(c tele.Context) error {
		err := c.EditOrSend(messages.SearchNotesByDateChooseMessage, &tele.SendOptions{
			ReplyMarkup: view.SearchByDateBtn(),
			ParseMode:   "html",
		})
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// поиск заметок по одной дате
	s.bot.Handle(&view.BtnSearchByOneDate, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].SearchNoteOneDate)

		s.controller.SetupNoteCalendar(ctx, c)
		s.controller.SetNoteCalendar(c.Chat().ID)

		err := s.controller.SearchNoteByOnedate(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				s.fsm[c.Chat().ID].SetNext()
				err := s.fsm[c.Chat().ID].Handle(ctx, c)
				if err != nil {
					s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
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
	s.bot.Handle(&view.BtnSearchByTwoDate, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].SearchNoteTwoDates)

		s.controller.SetupNoteCalendar(ctx, c)
		s.controller.SetNoteCalendar(c.Chat().ID)

		err := s.controller.SearchNoteByTwoDates(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				s.fsm[c.Chat().ID].SetNext()
				err := s.fsm[c.Chat().ID].Handle(ctx, c)
				if err != nil {
					if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
						return s.controller.SecondDateBeforeFirst(ctx, c)
					}

					s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
					return err
				}

				return nil
			})
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

	// удалить все напоминания
	s.bot.Handle(&view.BtnDeleteAllReminders, func(c tele.Context) error {
		err := s.controller.ConfirmDeleteAllReminders(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// удалить все напоминания - подтверждение
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

	// today
	s.bot.Handle(&view.BtnToday, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.Today(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

	// tomorrow
	s.bot.Handle(&view.BtnTomorrow, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		err := s.controller.Tomorrow(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		return nil
	})

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

		// err := s.controller.EveryWeek(ctx, c)
		// if err != nil {
		// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
		// 	return err
		// }

		// return nil

		return s.fsm[c.Chat().ID].Handle(ctx, c)
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

		s.controller.SetupReminderCalendar(ctx, c)

		s.controller.SetReminderCalendar(c.Chat().ID)

		err := s.controller.Year(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				// s.fsm[c.Chat().ID].SetNext()
				// err := s.fsm[c.Chat().ID].Handle(ctx, c)
				// if err != nil {
				// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
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
	s.bot.Handle(&view.BtnOnce, func(c tele.Context) error {
		s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].Once)

		s.controller.SetupReminderCalendar(ctx, c)
		s.controller.SetReminderCalendar(c.Chat().ID)

		err := s.controller.Date(ctx, c)
		if err != nil {
			s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
			return err
		}

		btns := s.controller.DaysBtns(ctx, c)

		for _, btn := range btns {
			s.bot.Handle(&btn, func(c tele.Context) error {
				return s.fsm[c.Chat().ID].Handle(ctx, c)
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
				// s.fsm[c.Chat().ID].SetNext()
				// err := s.fsm[c.Chat().ID].Handle(ctx, c)
				// if err != nil {
				// 	// если режим - дата (одноразовое напоминание) и проверка не прошла
				// 	if errors.Is(err, api_errors.ErrInvalidDate) {
				// 		return s.controller.InvalidDate(ctx, c)
				// 	}
				// 	if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
				// 		return s.controller.SecondDateBeforeFirst(ctx, c)
				// 	}

				// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
				// 	return err
				// }

				// s.fsm[c.Chat().ID].SetNext()
				// s.controller.ResetCalendars(c.Chat().ID)
				// return nil

				return s.fsm[c.Chat().ID].Handle(ctx, c)
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
				// s.fsm[c.Chat().ID].SetNext()
				// err := s.fsm[c.Chat().ID].Handle(ctx, c)
				// if err != nil {
				// 	if errors.Is(err, api_errors.ErrInvalidDate) {
				// 		return s.controller.InvalidDate(ctx, c)
				// 	}
				// 	if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
				// 		return s.controller.SecondDateBeforeFirst(ctx, c)
				// 	}

				// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
				// 	return err
				// }

				// s.fsm[c.Chat().ID].SetNext()
				// s.controller.ResetCalendars(c.Chat().ID)
				return s.fsm[c.Chat().ID].Handle(ctx, c)
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
				// s.fsm[c.Chat().ID].SetNext()
				// err := s.fsm[c.Chat().ID].Handle(ctx, c)
				// if err != nil {
				// 	if errors.Is(err, api_errors.ErrInvalidDate) {
				// 		return s.controller.InvalidDate(ctx, c)
				// 	}
				// 	if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
				// 		return s.controller.SecondDateBeforeFirst(ctx, c)
				// 	}

				// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
				// 	return err
				// }

				// s.fsm[c.Chat().ID].SetNext()
				// s.controller.ResetCalendars(c.Chat().ID)
				return s.fsm[c.Chat().ID].Handle(ctx, c)
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
				// s.fsm[c.Chat().ID].SetNext()
				// err := s.fsm[c.Chat().ID].Handle(ctx, c)
				// if err != nil {
				// 	if errors.Is(err, api_errors.ErrSecondDateBeforeFirst) {
				// 		return s.controller.SecondDateBeforeFirst(ctx, c)
				// 	}

				// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
				// 	return err
				// }

				// s.fsm[c.Chat().ID].SetNext()
				// s.controller.ResetCalendars(c.Chat().ID)
				return s.fsm[c.Chat().ID].Handle(ctx, c)
			})
		}

		return nil
	})

	// week days

	// Monday
	s.bot.Handle(&view.MondayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
		// 	return err
		// }

		// return nil

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Tuesday
	s.bot.Handle(&view.TuesdayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Wednesday
	s.bot.Handle(&view.WednesdayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Thursday
	s.bot.Handle(&view.ThursdayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Friday
	s.bot.Handle(&view.FridayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Saturday
	s.bot.Handle(&view.SaturdayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Sunday
	s.bot.Handle(&view.SundayBtn, func(c tele.Context) error {
		// s.fsm[c.Chat().ID].SetState(s.fsm[c.Chat().ID].ReminderTime)

		// err := s.controller.WeekDay(ctx, c)
		// if err != nil {
		// 	s.controller.HandleError(c, err, s.fsm[c.Chat().ID].Name())
		// 	return err
		// }

		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})
}
