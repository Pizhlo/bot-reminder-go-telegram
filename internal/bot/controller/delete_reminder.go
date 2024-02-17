package controller

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	tele "gopkg.in/telebot.v3"
)

// DeleteReminder удаляет сработавшее напоминание
func (c *Controller) DeleteReminder(ctx context.Context, telectx tele.Context) error {
	// userID_reminderID - для удаления
	reminderAndUser := strings.Split(telectx.Callback().Unique, "_")

	reminderID, userID := reminderAndUser[0], reminderAndUser[1]

	reminderInt, err := strconv.Atoi(reminderID)
	if err != nil {
		return fmt.Errorf("error while converting string reminder ID to int: %w", err)
	}

	userInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return fmt.Errorf("error while converting string user ID to int64: %w", err)
	}

	// получаем айди задачи
	jobID, err := c.reminderSrv.GetJobID(ctx, userInt, reminderInt)
	if err != nil {
		return fmt.Errorf("error while getting job ID: %w", err)
	}

	// удаляем из шедулера
	c.scheduler.DeleteJob(jobID)

	// удаляем из базы
	return c.reminderSrv.DeleteReminderByID(ctx, userInt, reminderInt)
}
