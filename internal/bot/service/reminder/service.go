package reminder

import (
	"context"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	gocron "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/scheduler"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ReminderService отвечает за напоминания: удаление, создание, выдача
type ReminderService struct {
	reminderEditor reminderEditor
	logger         *logrus.Logger
	viewsMap       map[int64]*view.ReminderView
	// для сохранения напоминаний во время создания
	reminderMap map[int64]model.Reminder
	mu          sync.Mutex

	schedulers map[int64]*gocron.Scheduler
}

//go:generate mockgen -source ./service.go -destination=./mocks/reminder_editor.go
type reminderEditor interface {
	// Save сохраняет напоминание в базе данных. Для сохранения требуется: ID пользователя, содержимое напоминания, дата создания
	Save(ctx context.Context, reminder *model.Reminder) (int64, error)

	// GetAllByUserID достает из базы все напоминания пользователя по ID, возвращает ErrRemindersNotFound
	GetAllByUserID(ctx context.Context, userID int64) ([]model.Reminder, error)

	// SaveJob сохраняет задачу в базе
	SaveJob(ctx context.Context, userID, reminderID int64, jobID uuid.UUID) error

	// DeleteAllByUserID удаляет все напоминания пользователя по user ID
	DeleteAllByUserID(ctx context.Context, userID int64) error

	// GetAllJobs возвращает айди всех задач пользователя
	GetAllJobs(ctx context.Context, userID int64) ([]uuid.UUID, error)

	// GetJobID возвращает айди задачи по айди пользователя и напоминания
	GetJobID(ctx context.Context, userID int64, reminderID int) (uuid.UUID, error)

	// DeleteReminderByID удаляет одно напоминание. Для удаления необходим ID напоминания и пользователя
	DeleteReminderByID(ctx context.Context, userID int64, reminderID int) error

	// // GetByID возвращает напоминание с переданным ID. Если такого напоминания нет, возвращает ErrRemindersNotFound
	// GetByID(ctx context.Context, userID int64, noteID int) (*model.Reminder, error)

	// // SearchByText производит поиск по напоминаний по тексту. Если таких напоминаний нет, возвращает ErrRemindersNotFound
	// SearchByText(ctx context.Context, searchNote model.SearchByText) ([]model.Reminder, error)
}

func New(reminderEditor reminderEditor) *ReminderService {
	return &ReminderService{reminderEditor: reminderEditor, logger: logger.New(), viewsMap: make(map[int64]*view.ReminderView),
		reminderMap: make(map[int64]model.Reminder), mu: sync.Mutex{}, schedulers: make(map[int64]*gocron.Scheduler)}
}

// SaveUser сохраняет пользователя в мапе view
func (n *ReminderService) SaveUser(userID int64) {
	if _, ok := n.viewsMap[userID]; !ok {
		n.logger.Debugf("Reminder service: user %d not found in the views map. Saving...\n", userID)
		n.viewsMap[userID] = view.NewReminder()
	} else {
		n.logger.Debugf("Reminder service: user %d already saved in the views map.\n", userID)
	}

}

// SetupCalendar устанавливает месяц и год в календаре на текущие
func (n *ReminderService) SetupCalendar(userID int64) {
	n.viewsMap[userID].SetCurMonth()
	n.viewsMap[userID].SetCurYear()
}
