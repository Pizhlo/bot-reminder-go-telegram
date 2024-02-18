package reminder

import (
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// SaveName сохраняет название напоминания при создании
func (n *ReminderService) SaveName(userID int64, name string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	reminder := model.Reminder{
		TgID: userID,
		Name: name,
	}

	n.reminderMap[userID] = reminder
}

// SaveType сохраняет тип напоминания
func (n *ReminderService) SaveType(userID int64, reminderType model.ReminderType) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	r.Type = reminderType

	n.reminderMap[userID] = r

	return nil
}

// SaveCreatedField сохраняет в напоминании поле created в указанном часовом поясе
func (n *ReminderService) SaveCreatedField(userID int64, tz *time.Location) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	r.Created = time.Now().In(tz)

	n.reminderMap[userID] = r

	return nil
}

// ProcessTime обрабатывает время, которое прислал пользователь: валидирует и сохраняет в случае успеха
func (n *ReminderService) ProcessTime(userID int64, timeMsg string) error {
	// формат, в котором пользователь должен прислать время
	layout := "15:04"

	_, err := time.Parse(layout, timeMsg)
	if err != nil {
		return err
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	r.Time = timeMsg

	n.reminderMap[userID] = r

	return nil
}

// SaveDate сохраняет переданную дату напоминания
func (n *ReminderService) SaveDate(userID int64, date string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	r.Date = date

	n.reminderMap[userID] = r

	return nil
}

// GetFromMemory достает из кэша напоминание в текущем состоянии (могут быть не заполнены все поля)
func (n *ReminderService) GetFromMemory(userID int64) (*model.Reminder, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return nil, fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	return &r, nil
}

// SaveID сохраняет ID напоминания, указанное в базе
func (n *ReminderService) SaveID(userID int64, reminderID int64) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	r.ID = reminderID

	n.reminderMap[userID] = r

	return nil
}

// GetID возвращает ID напоминания
func (n *ReminderService) GetID(userID int64) (int64, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return 0, fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	return r.ID, nil
}
