package reminder

import (
	"context"
)

// DeleteAll удаляет все напоминания пользователя из базы
func (n *ReminderService) DeleteAll(ctx context.Context, userID int64) error {
	// удаляем из базы
	return n.reminderEditor.DeleteAllByUserID(ctx, userID)
}

// DeleteReminderByID удаляет одно напоминание. Для удаления необходим ID напоминания и пользователя
func (n *ReminderService) DeleteReminderByID(ctx context.Context, userID int64, reminderID int) error {
	return n.reminderEditor.DeleteReminderByID(ctx, userID, reminderID)
}
