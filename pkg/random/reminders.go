package random

import (
	"math/rand"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
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
	reminderTypes := []model.ReminderType{model.SeveralTimesDayType, model.EverydayType, model.EveryWeekType, model.SeveralDaysType,
		model.SeveralDaysType, model.OnceMonthType, model.OnceYearType, model.DateType}

	randomType := reminderTypes[rand.Intn(len(reminderTypes))]

	return model.Reminder{
		ID:      int64(randomInt(0, 10)),
		TgID:    1,
		Name:    String(10),
		Date:    String(10),
		Type:    randomType,
		Time:    "10:10",
		Created: time.Now(),
	}
}
