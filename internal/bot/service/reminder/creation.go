package reminder

import (
	"errors"
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// SaveName сохраняет название напоминания при создании
func (n *ReminderService) SaveName(userID int64, name string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	reminder, ok := n.reminderMap[userID]
	if !ok {
		reminder = model.Reminder{
			TgID: userID,
			Name: name,
		}
	} else {
		reminder.Name = name
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

func (n *ReminderService) SaveTime(userID int64, timeMsg string) error {
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

// SaveCalendarDate сохраняет дату, которая хранится в календаре
func (n *ReminderService) SaveCalendarDate(userID int64, dayOfMonth string) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	month := n.viewsMap[userID].Month()

	var date string

	monthStr := fixMonth(month)

	if r.Type == model.OnceYearType {
		date = fmt.Sprintf("%s.%s", dayOfMonth, monthStr)
	} else if r.Type == model.DateType {

		year := n.viewsMap[userID].Year()
		date = fmt.Sprintf("%s.%s.%d", dayOfMonth, monthStr, year)
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

// checkFields проверяет, заполнены ли все поля в напоминании
func (n *ReminderService) checkFields(userID int64) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	if r.TgID == 0 {
		return errors.New("field TgID is not filled")
	}

	if r.Name == "" {
		return errors.New("field Name is not filled")
	}

	if r.Type == "" {
		return errors.New("field Type is not filled")
	}

	if r.Date == "" {
		return errors.New("field Date is not filled")
	}

	if r.Time == "" {
		return errors.New("field Time is not filled")
	}

	if r.Created.IsZero() {
		return errors.New("field Created is not filled")
	}

	return nil
}
