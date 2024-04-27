package reminder

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	gocron "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/scheduler"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ReminderService отвечает за напоминания: удаление, создание, выдача
type ReminderService struct {
	reminderEditor reminderEditor
	viewsMap       map[int64]*view.ReminderView
	// для сохранения напоминаний во время создания
	reminderMap map[int64]*model.Reminder
	mu          sync.Mutex

	schedulers map[int64]*gocron.Scheduler
}

//go:generate mockgen -source ./service.go -destination=../../mocks/reminder_srv.go -package=mocks
type reminderEditor interface {
	// Save сохраняет напоминание в базе данных. Для сохранения требуется: ID пользователя, содержимое напоминания, дата создания
	Save(ctx context.Context, reminder *model.Reminder) (uuid.UUID, error)

	// SaveMemory сохраняет напоминание из кэша в базе данных. Напоминание может быть не полным, т.к. находится на стадии создания
	// и в памяти хранится промежуточный результат
	SaveMemory(ctx context.Context, reminder *model.Reminder) error

	// GetMemory возвращает все сохраненные напоминания, находящиеся в стадии создания
	GetMemory(ctx context.Context) ([]model.Reminder, error)

	// DeleteMemory удаляет промежуточное сохранение напоминания. Необходимо, когда напоминание создано и сохранено в основной таблице
	DeleteMemory(ctx context.Context, userID int64) error

	// GetAllByUserID достает из базы все напоминания пользователя по ID, возвращает ErrRemindersNotFound
	GetAllByUserID(ctx context.Context, userID int64) ([]model.Reminder, error)

	// SaveJob сохраняет задачу в базе
	SaveJob(ctx context.Context, reminderID uuid.UUID, jobID uuid.UUID) error

	// DeleteAllByUserID удаляет все напоминания пользователя по user ID
	DeleteAllByUserID(ctx context.Context, userID int64) error

	// GetAllJobs возвращает айди всех задач пользователя
	GetAllJobs(ctx context.Context, userID int64) ([]uuid.UUID, error)

	// GetJobID возвращает айди задачи по айди напоминания
	GetJobID(ctx context.Context, reminderID uuid.UUID) (uuid.UUID, error)

	// DeleteReminderByID удаляет одно напоминание. Для удаления необходим ID напоминания
	DeleteReminderByID(ctx context.Context, reminderID uuid.UUID) (uuid.UUID, error)

	// GetReminderID возвращает айди напоминания в базе. Ищет по пользователю и view id
	GetReminderID(ctx context.Context, userID int64, viewID int) (uuid.UUID, error)

	// DeleteJob удаляет таску из базы и связанное напоминание
	DeleteJobAndReminder(ctx context.Context, jobID uuid.UUID) error

	// GetByViewID возвращает напоминание с переданным viewID. Если такого напоминания нет, возвращает ErrRemindersNotFound
	GetByViewID(ctx context.Context, userID int64, viewID int) (*model.Reminder, error)
}

func New(reminderEditor reminderEditor) *ReminderService {
	return &ReminderService{
		reminderEditor: reminderEditor,
		viewsMap:       make(map[int64]*view.ReminderView),
		reminderMap:    make(map[int64]*model.Reminder),
		mu:             sync.Mutex{},
		schedulers:     make(map[int64]*gocron.Scheduler),
	}
}

// SaveUser сохраняет пользователя в мапе view
func (n *ReminderService) SaveUser(userID int64) {
	if _, ok := n.viewsMap[userID]; !ok {
		logrus.Debugf(wrap(fmt.Sprintf("user %d not found in the views map. Saving...\n", userID)))
		n.viewsMap[userID] = view.NewReminder()
	} else {
		logrus.Debugf(wrap(fmt.Sprintf("user %d already saved in the views map.\n", userID)))
	}

}

// LoadMemory загружает из базы все недоделанные напоминания
func (n *ReminderService) LoadMemory(ctx context.Context) error {
	rList, err := n.reminderEditor.GetMemory(ctx)
	if err != nil {
		return err
	}

	for _, r := range rList {
		logrus.Debugf(wrap(fmt.Sprintf("saving user's reminder from db to cache: %+v", r)))
		n.reminderMap[r.TgID] = &r
	}

	return nil
}

// SetupCalendar устанавливает месяц и год в календаре на текущие
func (n *ReminderService) SetupCalendar(userID int64) {
	n.viewsMap[userID].SetCurMonth()
	n.viewsMap[userID].SetCurYear()
}

// SaveMemory сохраняет напоминания, хранящиеся в памяти, в БД
func (n *ReminderService) SaveMemory(ctx context.Context) error {
	var res error
	for _, r := range n.reminderMap {
		// проверяем, чтобы напоминание не было пустым
		if r.TgID > 0 {
			err := n.reminderEditor.SaveMemory(ctx, r)
			res = errors.Join(err)
		}

	}

	return res
}

func wrap(s string) string {
	return fmt.Sprintf("Reminder service: %s", s)
}
