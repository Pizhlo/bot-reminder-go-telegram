package server

import (
	"context"
	"errors"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/controller"
	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	tele "gopkg.in/telebot.v3"
)

func (s *Server) setupNotesHandlers(ctx context.Context, restricted *tele.Group) {
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
}
