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
		Date:    randomDate().Format("02.01.2006"),
		Type:    model.DateType,
		Time:    "10:10",
		Created: time.Now(),
		Job: model.Job{
			ID:      uuid.Nil,
			NextRun: time.Time{}, // в реальности в базе такие вещи не хранятся и из базы поле приходит пустым
		},
	}

	return r
}

func randomDate() time.Time {
	day := Int(1, 10)
	month := time.Now().Add(30 * 24 * time.Hour).Month() // прибавляем 30 дней
	year := time.Now().Year()

	date := time.Date(year, month, day, 0, 0, 0, 0, time.Local)

	return date
}
