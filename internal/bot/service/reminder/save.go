package reminder

import (
	"context"

	"github.com/google/uuid"
)

// Save сохраняет напоминание
func (s *ReminderService) Save(ctx context.Context, userID int64) error {
	r, err := s.GetFromMemory(userID)
	if err != nil {
		return err
	}

	s.logger.Debugf("Reminder service: saving user's reminder. Model: %+v\n", r)

	id, err := s.reminderEditor.Save(ctx, r)
	if err != nil {
		return err
	}

	s.SaveID(userID, id)

	return nil
}

// SaveJobID сохраняет в базе ID задачи, связанной с напоминанием
func (s *ReminderService) SaveJobID(ctx context.Context, jobID uuid.UUID, userID int64) error {
	r, err := s.GetFromMemory(userID)
	if err != nil {
		return err
	}

	s.logger.Debugf("Reminder service: saving user's job. Model: %+v\n", r)

	return s.reminderEditor.SaveJob(ctx, userID, r.ID, jobID)
}
