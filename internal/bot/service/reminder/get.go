package reminder

import (
	"context"

	"github.com/google/uuid"
	tele "gopkg.in/telebot.v3"
)

// GetAll обрабатывает запрос всех напоминаний пользователя. Возвращает: первую страницу напоминаний (в виде целого сообщения),
// клавиатуру и ошибку
func (s *ReminderService) GetAll(ctx context.Context, userID int64) (string, *tele.ReplyMarkup, error) {
	s.logger.Debugf("Reminder service: getting all user's reminders. User ID: %d\n", userID)

	reminders, err := s.reminderEditor.GetAllByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Reminder service: error while getting all reminders by user ID %d: %v\n", userID, err)
		return "", nil, err
	}

	s.logger.Debugf("Reminder service: got %d user's reminders\n", len(reminders))

	return s.viewsMap[userID].Message(reminders), s.viewsMap[userID].Keyboard(), nil
}

// GetAllJobs возвращает id всех задач, созданных пользователем
func (s *ReminderService) GetAllJobs(ctx context.Context, userID int64) ([]uuid.UUID, error) {
	return s.reminderEditor.GetAllJobs(ctx, userID)
}

// GetJobID возвращает айди задачи по айди пользователя и напоминания
func (s *ReminderService) GetJobID(ctx context.Context, userID int64, reminderID int) (uuid.UUID, error) {
	return s.reminderEditor.GetJobID(ctx, userID, reminderID)
}
