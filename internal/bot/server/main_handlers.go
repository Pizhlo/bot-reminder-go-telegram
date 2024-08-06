package server

import (
	"context"
	"strings"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// setupMainHandlers настраивает хендлеры на кнопки главного меню: Заметки, Напоминания, Часовой пояс, Сообщить о баге,
// а также кнопку Назад в меню и кнопку сработавшего напоминания (галочку)
func (s *Server) setupMainHandlers(ctx context.Context, restricted *tele.Group) {
	// text
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

	// геолокация
	s.bot.Handle(tele.OnLocation, func(telectx tele.Context) error {
		logrus.Debugf("location")

		err := s.controller.AcceptTimezone(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// restricted: only known users

	// часовой пояс
	restricted.Handle(&view.BtnTimezone, func(telectx tele.Context) error {
		err := s.controller.Timezone(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// изменить часовой пояс
	restricted.Handle(&view.BtnEditTimezone, func(telectx tele.Context) error {
		err := s.controller.RequestLocation(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})

	// меню заметок
	restricted.Handle(&view.BtnNotes, func(telectx tele.Context) error {
		s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].ListNote)
		err := s.controller.ListNotes(ctx, telectx)
		if err != nil {
			switch t := err.(type) {
			case *tele.Error:
				if strings.Contains(t.Description, "message is not modified") {
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
		s.fsm[telectx.Chat().ID].SetState(s.fsm[telectx.Chat().ID].ListReminder)
		err := s.controller.ListReminders(ctx, telectx)
		if err != nil {
			switch t := err.(type) {
			case *tele.Error:
				if strings.Contains(t.Description, "message is not modified") {
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
		s.fsm[telectx.Chat().ID].SetToDefault()

		err := s.controller.MenuCmd(ctx, telectx)
		if err != nil {
			s.HandleError(telectx, err)
			return err
		}

		return nil
	})
}
