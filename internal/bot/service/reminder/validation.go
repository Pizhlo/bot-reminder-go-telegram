package reminder

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
)

// ParseTime обрабатывает время, которое прислал пользователь: валидирует и сохраняет в случае успеха
func (n *ReminderService) ParseTime(userID int64, timeMsg string) error {
	// формат, в котором пользователь должен прислать время
	layout := "15:04"

	_, err := time.Parse(layout, timeMsg)
	if err != nil {
		return api_errors.ErrInvalidTime
	}

	return nil
}

// ValidateTime проверет, что время не прошло
func (n *ReminderService) ValidateTime(loc *time.Location, userTime time.Time) error {
	now := time.Now().In(loc)

	if userTime.Before(now) {
		return api_errors.ErrTimeInPast
	}

	return nil
}

func (n *ReminderService) ValidateDate(userID int64, dayOfMonth string, timezone *time.Location) error {
	month := n.viewsMap[userID].Month()
	year := n.viewsMap[userID].Year()

	dayInt, err := n.checkIfInt(dayOfMonth)
	if err != nil {
		return api_errors.ErrInvalidDate
	}

	now := time.Now().In(timezone)

	userDate := time.Date(year, month, dayInt, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), timezone)

	// если дата уже прошла - не можем создать уведомление на эту дату, возвращаем ошибку
	if now.After(userDate) {
		return api_errors.ErrInvalidDate
	}

	return nil
}

func fixInt(i int) string {
	if i < 10 {
		return "0" + strconv.Itoa(int(i))
	}

	return strconv.Itoa(int(i))
}

// ProcessMinutes обрабатывает количество минут: валидирует (число должно быть от 1 до 59) и сохраняет
func (n *ReminderService) ProcessMinutes(userID int64, minutes string) error {
	// проверяем, является ли пользовательский ввод числом
	minuesInt, err := n.checkIfInt(minutes)
	if err != nil {
		return err
	}

	// проверяем на соответствие требованиям
	if minuesInt < 1 || minuesInt > 59 {
		return errors.New("must be in within the range from 1 to 59")
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	// сохраняем изменения

	r.Time = minutes

	n.reminderMap[userID] = r

	return nil
}

// ProcessMinutes обрабатывает количество часов: валидирует (число должно быть от 1 до 59) и сохраняет
func (n *ReminderService) ProcessHours(userID int64, hours string) error {
	// проверяем, является ли пользовательский ввод числом
	hoursInt, err := n.checkIfInt(hours)
	if err != nil {
		return err
	}

	// проверяем на соответствие требованиям
	if hoursInt < 1 || hoursInt > 24 {
		return errors.New("must be in within the range from 1 to 24")
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	// сохраняем изменения

	r.Time = hours

	n.reminderMap[userID] = r

	return nil
}

// checkIfInt проверяет, является ли строка числом
func (n *ReminderService) checkIfInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// ProcessDaysInMonth валидирует количество дней, введенных пользователем.
// Количество должно быть в диапазоне [1, 31].
// Используется для напоминаний раз в месяц
func (n *ReminderService) ProcessDaysInMonth(userID int64, days string) error {
	daysInt, err := n.checkIfInt(days)
	if err != nil {
		return api_errors.ErrInvalidDays
	}

	if daysInt < 1 || daysInt > 31 {
		return api_errors.ErrInvalidDays
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	// сохраняем изменения

	r.Date = days

	n.reminderMap[userID] = r

	return nil
}

// ProcessDaysInMoth валидирует количество дней, введенных пользователем.
// Используется для напоминаний раз в несколько дней
func (n *ReminderService) ProcessDaysDuration(userID int64, days string) error {
	daysInt, err := n.checkIfInt(days)
	if err != nil {
		return api_errors.ErrInvalidDays
	}

	if daysInt < 1 || daysInt > 180 {
		return api_errors.ErrInvalidDays
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
	}

	// сохраняем изменения

	r.Date = days

	n.reminderMap[userID] = r

	return nil
}
