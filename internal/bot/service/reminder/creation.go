package reminder

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
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

// ProcessTime обрабатывает время, которое прислал пользователь: валидирует и сохраняет в случае успеха
func (n *ReminderService) ProcessTime(userID int64, timeMsg string) error {
	// формат, в котором пользователь должен прислать время
	layout := "15:04"

	_, err := time.Parse(layout, timeMsg)
	if err != nil {
		return err
	}

	return n.saveTime(userID, timeMsg)
}

func (n *ReminderService) saveTime(userID int64, timeMsg string) error {
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

func (n *ReminderService) ValidateDate(userID int64, dayOfMonth string, timezone *time.Location) error {
	month := n.viewsMap[userID].Month()
	year := n.viewsMap[userID].Year()

	dayInt, err := n.checkIfInt(dayOfMonth)
	if err != nil {
		return err
	}

	now := time.Now().In(timezone)

	userDate := time.Date(year, month, dayInt, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), timezone)

	// если дата уже прошла - не можем создать уведомление на эту дату, возвращаем ошибку
	if now.After(userDate) {
		return api_errors.ErrInvalidDate
	}

	return nil
}

func fixMonth(month time.Month) string {
	if month < 10 {
		return "0" + strconv.Itoa(int(month))
	}

	return strconv.Itoa(int(month))
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

// checkFields проверяет, заполнены ли все поля в напоминании
func (n *ReminderService) checkFields(userID int64) error {
	n.mu.Lock()
	defer n.mu.Unlock()

	r, ok := n.reminderMap[userID]
	if !ok {
		return fmt.Errorf("error while getting reminder by user ID: reminder not found")
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
