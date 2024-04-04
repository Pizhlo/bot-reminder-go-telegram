package reminder

import (
	"context"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Save проверяет заполненость полей сохраняет напоминание в БД
func (s *ReminderService) Save(ctx context.Context, userID int64) error {
	r, err := s.GetFromMemory(userID)
	if err != nil {
		return err
	}

	// проверяем, заполнены ли все поля в напоминании
	if err := s.checkFields(r); err != nil {
		return err
	}

	id, err := s.reminderEditor.Save(ctx, r)
	if err != nil {
		return err
	}

	return s.SaveID(userID, id)
}

// SaveJobID сохраняет в базе ID задачи, связанной с напоминанием
func (s *ReminderService) SaveJobID(ctx context.Context, jobID uuid.UUID, userID int64, reminderID uuid.UUID) error {
	// r, err := s.GetFromMemory(userID)
	// if err != nil {
	// 	return err
	// }

	logrus.Debugf(wrap(fmt.Sprintf("saving user's job. UUID: %+v. Reminder ID: %v\n", jobID, reminderID)))

	return s.reminderEditor.SaveJob(ctx, reminderID, jobID)
}

// Clear очищает память после успешного сохранения
func (s *ReminderService) Clear(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.reminderMap[userID] = model.Reminder{}
}
