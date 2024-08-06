package server

import (
	"context"
	"errors"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

func (s *Server) setupCalendarHandlers(ctx context.Context, restricted *tele.Group) {
	// calendar

	// months buttons
	// restricted.Handle(&view.BtnJanuary, func(c tele.Context) error {
	// 	return s.controller.ListMonths(ctx, c)
	// })

	restricted.Handle(&view.BtnMonth, func(c tele.Context) error {
		return s.controller.ListMonths(ctx, c)
	})

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
		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Tuesday
	restricted.Handle(&view.TuesdayBtn, func(c tele.Context) error {
		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Wednesday
	restricted.Handle(&view.WednesdayBtn, func(c tele.Context) error {
		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Thursday
	restricted.Handle(&view.ThursdayBtn, func(c tele.Context) error {
		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Friday
	restricted.Handle(&view.FridayBtn, func(c tele.Context) error {
		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Saturday
	restricted.Handle(&view.SaturdayBtn, func(c tele.Context) error {
		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})

	// Sunday
	restricted.Handle(&view.SundayBtn, func(c tele.Context) error {
		return s.fsm[c.Chat().ID].Handle(ctx, c)
	})
}
