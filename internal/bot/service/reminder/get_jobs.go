package reminder

import (
	"context"

	"github.com/google/uuid"
)

// GetAllJobs возвращает id всех задач, созданных пользователем
func (s *ReminderService) GetAllJobs(ctx context.Context, userID int64) ([]uuid.UUID, error) {
	return s.reminderEditor.GetAllJobs(ctx, userID)
}

// GetJobID возвращает айди задачи по айди пользователя и напоминания
func (s *ReminderService) GetJobID(ctx context.Context, userID int64, reminderID int) (uuid.UUID, error) {
	return s.reminderEditor.GetJobID(ctx, userID, reminderID)
}
