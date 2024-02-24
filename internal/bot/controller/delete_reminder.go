package controller

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

// DeleteReminder удаляет напоминание
func (c *Controller) DeleteReminder(ctx context.Context, telectx tele.Context, reminderID int, userID int64) error {
	// получаем айди задачи
	jobID, err := c.reminderSrv.GetJobID(ctx, userID, reminderID)
	if err != nil {
		return fmt.Errorf("error while getting job ID: %w", err)
	}

	// удаляем из шедулера
	err = c.scheduler.DeleteJob(jobID)
	if err != nil {
		return fmt.Errorf("error while deleting job: %w", err)
	}

	// удаляем из базы
	return c.reminderSrv.DeleteReminderByID(ctx, userID, reminderID)
}
