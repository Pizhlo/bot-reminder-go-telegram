package reminder

import (
	"context"
)

// Save сохраняет напоминание
func (s *ReminderService) Save(ctx context.Context, userID int64) error {
	r := s.GetFromMemory(userID)

	s.logger.Debugf("Reminder service: saving user's reminder. Model: %+v\n", r)

	return s.reminderEditor.Save(ctx, r)
}
