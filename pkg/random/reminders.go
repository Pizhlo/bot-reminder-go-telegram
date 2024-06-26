package random

import (
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/google/uuid"
)

// Reminders генерирует указанное количество напоминаний
func Reminders(n int) []model.Reminder {
	var reminders []model.Reminder

	for i := 0; i < n; i++ {
		r := Reminder()
		reminders = append(reminders, *r)
	}

	return reminders
}

// Reminder генерирует одно напоминание, заполненное рандомными данными
func Reminder() *model.Reminder {
	//weekDays := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

	r := &model.Reminder{
		ID:      uuid.New(),
		ViewID:  int64(Int(1, 10)),
		TgID:    1,
		Name:    String(10),
		Date:    "10.10.2024",
		Type:    model.DateType,
		Time:    "10:10",
		Created: time.Now(),
		Job: model.Job{
			ID:      uuid.Nil,
			NextRun: time.Date(time.Now().Year(), time.October, 10, 10, 10, 0, 0, time.Local),
		},
	}

	return r
}
