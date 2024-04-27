package reminder

import (
	"context"
	"errors"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
)

// DeleteAll удаляет все напоминания пользователя из базы, останавливает таски
func (n *ReminderService) DeleteAll(ctx context.Context, userID int64) error {
	// получаем айди всех задач пользователя, чтобы их остановить в шедулере
	jobIDs, err := n.reminderEditor.GetAllJobs(ctx, userID)
	if err != nil {
		logrus.Errorf(wrap(fmt.Sprintf("error getting all jobs' IDs. User ID: %d. Error: %+v\n", userID, err)))
		return err
	}

	var resultErr error

	for _, id := range jobIDs {
		// удаляем из БД
		err = n.reminderEditor.DeleteJobAndReminder(ctx, id)
		if err != nil {
			logrus.Errorf(wrap(fmt.Sprintf("error while deleting job from DB while deleting all jobs: %v\n", err)))
			resultErr = errors.Join(err)
			continue
		}

		// удаляем из шедулера
		err := n.DeleteJob(userID, id)
		if err != nil {
			resultErr = errors.Join(err)
			logrus.Errorf(wrap(fmt.Sprintf("error while deleting all jobs from scheduler: deleting job: %v\n", err)))
			continue
		}
	}

	//return n.reminderEditor.DeleteAllByUserID(ctx, userID)

	return resultErr
}

// deleteReminderByID удаляет одно напоминание. Для удаления необходим ID напоминания и пользователя
func (n *ReminderService) deleteReminderByID(ctx context.Context, reminder *model.Reminder) error {
	// удаляем из шедулера
	err := n.DeleteJob(reminder.TgID, reminder.Job.ID)
	if err != nil {
		return fmt.Errorf(wrap(fmt.Sprintf("error while deleting job: %v", err)))
	}

	_, err = n.reminderEditor.DeleteReminderByID(ctx, reminder.ID)
	// удаляем из базы
	return err
}

// DeleteByViewID удаляет напоминание по айди, которое видит пользователь (viewID)
func (n *ReminderService) DeleteByViewID(ctx context.Context, userID int64, viewID int) (string, error) {
	// получаем уникальный айди напоминания
	r, err := n.reminderEditor.GetByViewID(ctx, userID, viewID)
	if err != nil {
		return "", err
	}

	// обрабатываем
	return r.Name, n.deleteReminderByID(ctx, r)
}

// DeleteReminder удаляет переданное напоминание из БД.
// Используется для удаления одноразовых напоминаний
func (n *ReminderService) DeleteReminder(ctx context.Context, reminder *model.Reminder) error {
	job, err := n.reminderEditor.DeleteReminderByID(ctx, reminder.ID)
	if err != nil {
		return fmt.Errorf(wrap(fmt.Sprintf("error while deleting reminder: %v", err)))
	}

	logrus.Errorf("got job ID: %v", job)
	return n.DeleteJob(reminder.TgID, job)
}
