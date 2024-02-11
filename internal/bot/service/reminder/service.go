package reminder

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
)

// ReminderService отвечает за напоминания: удаление, создание, выдача
type ReminderService struct {
	reminderEditor reminderEditor
	logger         *logrus.Logger
	viewsMap       map[int64]*view.ReminderView
}

//go:generate mockgen -source ./service.go -destination=./mocks/reminder_editor.go
type reminderEditor interface {
	// Save сохраняет напоминание в базе данных. Для сохранения требуется: ID пользователя, содержимое напоминания, дата создания
	Save(ctx context.Context, reminder model.Reminder) error

	// GetAllByUserID достает из базы все напоминания пользователя по ID, возвращает ErrRemindersNotFound
	GetAllByUserID(ctx context.Context, userID int64) ([]model.Reminder, error)

	// // DeleteAllByUserID удаляет все напоминания пользователя по user ID
	// DeleteAllByUserID(ctx context.Context, userID int64) error

	// // DeleteReminderByID удаляет одно напоминание. Для удаления необходим ID напоминания и пользователя
	// DeleteReminderByID(ctx context.Context, userID int64, noteID int) error

	// // GetByID возвращает напоминание с переданным ID. Если такого напоминания нет, возвращает ErrRemindersNotFound
	// GetByID(ctx context.Context, userID int64, noteID int) (*model.Reminder, error)

	// // SearchByText производит поиск по напоминаний по тексту. Если таких напоминаний нет, возвращает ErrRemindersNotFound
	// SearchByText(ctx context.Context, searchNote model.SearchByText) ([]model.Reminder, error)
}

func New(reminderEditor reminderEditor) *ReminderService {
	return &ReminderService{reminderEditor: reminderEditor, logger: logger.New(), viewsMap: make(map[int64]*view.ReminderView)}
}

// SaveUser сохраняет пользователя в мапе view
func (n *ReminderService) SaveUser(userID int64) {
	n.logger.Debugf("Reminder service: checking if user saved in the views map...\n")
	if _, ok := n.viewsMap[userID]; !ok {
		n.logger.Debugf("Reminder service: user not found in the views map. Saving...\n")
		n.viewsMap[userID] = view.NewReminder()
	} else {
		n.logger.Debugf("Reminder service: user already saved in the views map.\n")
	}

	n.logger.Debugf("Reminder service: successfully saved user in the views map.\n")
}
