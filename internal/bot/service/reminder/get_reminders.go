package reminder

import (
	"context"
	"errors"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

// GetAll обрабатывает запрос всех напоминаний пользователя. Возвращает: первую страницу напоминаний (в виде целого сообщения),
// клавиатуру и ошибку
func (s *ReminderService) GetAll(ctx context.Context, userID int64) ([]model.Reminder, error) {
	reminders, err := s.reminderEditor.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sch, err := s.getScheduler(userID)
	if err != nil {
		return nil, err
	}

	jobs := sch.Jobs()
	jobsMap := map[uuid.UUID]gocron.Job{}

	for _, j := range jobs {
		jobsMap[j.ID()] = j
	}

	for i := 0; i < len(reminders); i++ {
		j, ok := jobsMap[reminders[i].Job.ID]
		if !ok {
			return nil, errors.New("not found job in JobsMap by ID")
		}

		nextRun, err := j.NextRun()
		if err != nil {
			return nil, err
		}

		reminders[i].Job.NextRun = nextRun
	}

	return reminders, nil
}

func (s *ReminderService) Message(userID int64, reminders []model.Reminder) (string, error) {
	return s.viewsMap[userID].Message(reminders)
}

func (s *ReminderService) Keyboard(userID int64) *tele.ReplyMarkup {
	return s.viewsMap[userID].Keyboard()
}

// NextPage обрабатывает кнопку переключения на следующую страницу
func (s *ReminderService) NextPage(userID int64) string {
	return s.viewsMap[userID].Next()
}

// PrevPage обрабатывает кнопку переключения на предыдущую страницу
func (s *ReminderService) PrevPage(userID int64) string {
	return s.viewsMap[userID].Previous()
}

// LastPage обрабатывает кнопку переключения на последнюю страницу
func (s *ReminderService) LastPage(userID int64) string {
	return s.viewsMap[userID].Last()
}

// FirstPage обрабатывает кнопку переключения на первую страницу
func (s *ReminderService) FirstPage(userID int64) string {
	return s.viewsMap[userID].First()
}
