package reminder

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// Save сохраняет напоминание
func (s *ReminderService) Save(ctx context.Context, reminder model.Reminder) error {
	s.logger.Debugf("Reminder service: saving user's reminder. Model: %+v\n", reminder)

	return s.reminderEditor.Save(ctx, reminder)
}
