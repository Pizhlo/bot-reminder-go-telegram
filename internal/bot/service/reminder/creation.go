package reminder

import (
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

// SaveReminderName сохраняет название напоминания при создании
func (n *ReminderService) SaveReminderName(userID int64, name string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	reminder := model.Reminder{
		TgID: userID,
		Text: name,
	}

	n.reminderMap[userID] = reminder
}

func (n *ReminderService) SaveType(userID int64, typeMsg string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	r := n.reminderMap[userID]
	r.Type = typeMsg

	n.reminderMap[userID] = r
}

// ProcessTime обрабатывает время, которое прислал пользователь: валидирует и сохраняет в случае успеха
func (n *ReminderService) ProcessTime(userID int64, timeMsg string) error {
	layout := "15:04"

	_, err := time.Parse(layout, timeMsg)
	if err != nil {
		return err
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	r := n.reminderMap[userID]
	r.Time = timeMsg

	n.reminderMap[userID] = r

	return nil
}

func (n *ReminderService) GetFromMemory(userID int64) *model.Reminder {
	n.mu.Lock()
	defer n.mu.Unlock()

	r := n.reminderMap[userID]

	return &r
}
