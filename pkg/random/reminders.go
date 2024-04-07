package random

import (
	"math/rand"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/google/uuid"
)

// Reminders генерирует указанное количество напоминаний
func Reminders(n int) []model.Reminder {
	var reminders []model.Reminder

	for i := 0; i < n; i++ {
		r := Reminder()
		reminders = append(reminders, r)
	}

	return reminders
}

// Reminder генерирует одно напоминание, заполненное рандомными данными
func Reminder() model.Reminder {
	reminderTypes := []model.ReminderType{model.EverydayType, model.EveryWeekType}
	weekDays := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

	randomType := reminderTypes[rand.Intn(len(reminderTypes))]

	r := model.Reminder{
		ID:      uuid.New(),
		ViewID:  int64(Int(1, 10)),
		TgID:    1,
		Name:    String(10),
		Date:    String(10),
		Type:    randomType,
		Time:    "10:10",
		Created: time.Now(),
	}

	if r.Type == model.EveryWeekType {
		randomWd := weekDays[rand.Intn(len(weekDays))]
		r.Date = randomWd
	}

	return r
}
