package reminder

import (
	"context"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
)

// Save сохраняет напоминание. Tz - часовой пояс пользователя (чтобы установить поле created)
func (s *ReminderService) Save(ctx context.Context, userID int64, tz *user.Timezone) error {
	r := s.GetFromMemory(userID)

	t := time.Now()
	utc, err := time.LoadLocation(tz.Name)
	if err != nil {
		return fmt.Errorf("error while setting timezone (for setting 'created' field): %w", err)
	}

	r.Created = t.In(utc)

	s.logger.Debugf("Reminder service: saving user's reminder. Model: %+v\n", r)

	return s.reminderEditor.Save(ctx, r)
}
