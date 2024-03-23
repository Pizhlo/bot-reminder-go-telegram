package reminder

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveReminderName_New(t *testing.T) {
	n := New(nil)

	userID := int64(1)
	reminderName := random.String(10)

	n.SaveName(userID, reminderName)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, result.TgID)
	assert.Equal(t, reminderName, result.Name)
}

func TestSaveReminderName_Exist(t *testing.T) {
	n := New(nil)

	userID := int64(1)
	reminderName1 := random.String(10)
	reminderName2 := random.String(10)

	n.SaveName(userID, reminderName1)
	n.SaveName(userID, reminderName2)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, result.TgID)
	assert.Equal(t, reminderName2, result.Name)
}

func TestSaveType_NotFound(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	err := n.SaveType(userID, model.DateType)
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")

	_, ok := n.reminderMap[userID]
	assert.Equal(t, false, ok)
}

func TestSaveType(t *testing.T) {
	type test struct {
		userID       int64
		reminderType model.ReminderType
		reminderName string
	}

	tests := []test{
		{
			userID:       int64(1),
			reminderType: model.SeveralTimesDayType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.EverydayType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.EveryWeekType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.SeveralDaysType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.OnceMonthType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.OnceYearType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.DateType,
			reminderName: random.String(10),
		},
	}

	for _, tt := range tests {
		n := New(nil)

		// сначала сохраняем напоминание в мапу
		n.SaveName(tt.userID, tt.reminderName)

		// сохраняем тип
		err := n.SaveType(tt.userID, tt.reminderType)
		assert.NoError(t, err)

		result, ok := n.reminderMap[tt.userID]
		assert.Equal(t, true, ok, "reminder not found in the map")

		assert.Equal(t, tt.userID, result.TgID, fmt.Sprintf("expected userID: %d. Actual: %d", tt.userID, result.ID))
		assert.Equal(t, tt.reminderType, result.Type, fmt.Sprintf("expected type: %s. Actual: %s", tt.reminderType, result.Type))
		assert.Equal(t, tt.reminderName, result.Name)
	}

}

func TestSaveCreatedField_NotFound(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	err := n.SaveCreatedField(userID, time.Local)
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")

	_, ok := n.reminderMap[userID]
	assert.Equal(t, false, ok)

}

func TestSaveCreatedField(t *testing.T) {
	type test struct {
		userID         int64
		reminderName   string
		timezoneString string
	}

	tests := []test{
		{
			userID:         int64(1),
			timezoneString: "Europe/Moscow",
			reminderName:   random.String(10),
		},
		{
			userID:         int64(1),
			timezoneString: "Europe/London",
			reminderName:   random.String(10),
		},
		{
			userID:         int64(1),
			timezoneString: "Europe/Oslo",
			reminderName:   random.String(10),
		},
		{
			userID:         int64(1),
			timezoneString: "Europe/Paris",
			reminderName:   random.String(10),
		},
	}

	for _, tt := range tests {
		n := New(nil)

		n.SaveName(tt.userID, tt.reminderName)

		location, err := time.LoadLocation(tt.timezoneString)
		if err != nil {
			t.Errorf("error while loading location: %v", err)
		}

		err = n.SaveCreatedField(tt.userID, location)
		assert.NoError(t, err)

		result, ok := n.reminderMap[tt.userID]
		assert.Equal(t, true, ok)

		assert.Equal(t, tt.userID, result.TgID)

		created := result.Created

		assert.Equal(t, location, created.Location(), fmt.Sprintf("locations not equal: expected %s. Actual: %s", location, created.Location()))
		assert.Equal(t, tt.reminderName, result.Name)
	}
}

func TestProcessTime_ValidTime(t *testing.T) {
	type test struct {
		userID       int64
		timeMsg      string
		reminderName string
	}

	tests := []test{
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      "12:01",
		},
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      "23:59",
		},
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      "00:01",
		},
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      "10:10",
		},
	}

	for _, tt := range tests {
		n := New(nil)

		n.SaveName(tt.userID, tt.reminderName)

		err := n.ParseTime(tt.userID, tt.timeMsg)
		assert.NoError(t, err)
	}
}

func TestProcessTime_InvalidTime(t *testing.T) {
	type test struct {
		userID       int64
		reminderName string
		timeMsg      string
	}

	tests := []test{
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      random.String(10),
		},
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      "36:79",
		},
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      "ab:cd",
		},
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      "10:10:15",
		},
		{
			userID:       int64(1),
			reminderName: random.String(10),
			timeMsg:      "12345",
		},
	}

	for _, tt := range tests {
		n := New(nil)

		n.SaveName(tt.userID, tt.reminderName)

		err := n.ParseTime(tt.userID, tt.timeMsg)

		assert.Error(t, err)

		result, ok := n.reminderMap[tt.userID]
		assert.Equal(t, true, ok)

		assert.Equal(t, tt.reminderName, result.Name)

	}
}

func TestSaveDate(t *testing.T) {
	n := New(nil)

	userID := int64(1)
	date := random.String(10)
	reminderName := random.String(10)

	n.SaveName(userID, reminderName)

	err := n.SaveDate(userID, date)
	assert.NoError(t, err)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, result.TgID)
	assert.Equal(t, reminderName, result.Name)
	assert.Equal(t, date, result.Date)
}

func TestSaveDate_NotFound(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	err := n.SaveDate(userID, random.String(10))
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")

	_, ok := n.reminderMap[userID]
	assert.Equal(t, false, ok)
}

func TestGetFromMemory_OnlyName(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)

	n := New(nil)
	n.SaveName(userID, reminderName)

	result, err := n.GetFromMemory(userID)
	assert.NoError(t, err)

	assert.Equal(t, userID, result.TgID)
	assert.Equal(t, reminderName, result.Name)
}

func TestGetFromMemory_NameAndType(t *testing.T) {
	type test struct {
		userID       int64
		reminderType model.ReminderType
		reminderName string
	}

	tests := []test{
		{
			userID:       int64(1),
			reminderType: model.SeveralTimesDayType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.EverydayType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.EveryWeekType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.SeveralDaysType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.OnceMonthType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.OnceYearType,
			reminderName: random.String(10),
		},
		{
			userID:       int64(1),
			reminderType: model.DateType,
			reminderName: random.String(10),
		},
	}

	for _, tt := range tests {
		n := New(nil)
		n.SaveName(tt.userID, tt.reminderName)
		err := n.SaveType(tt.userID, tt.reminderType)
		assert.NoError(t, err)

		result, err := n.GetFromMemory(tt.userID)
		assert.NoError(t, err)

		assert.Equal(t, tt.userID, result.TgID)
		assert.Equal(t, tt.reminderName, result.Name)
		assert.Equal(t, tt.reminderType, result.Type)
	}
}

func TestGetFromMemory_NameTypeCreated(t *testing.T) {
	type test struct {
		userID         int64
		reminderType   model.ReminderType
		reminderName   string
		timezoneString string
	}

	tests := []test{
		{
			userID:         int64(1),
			reminderType:   model.SeveralTimesDayType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
		},
		{
			userID:         int64(1),
			reminderType:   model.EverydayType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
		},
		{
			userID:         int64(1),
			reminderType:   model.EveryWeekType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
		},
		{
			userID:         int64(1),
			reminderType:   model.SeveralDaysType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
		},
		{
			userID:         int64(1),
			reminderType:   model.OnceMonthType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
		},
		{
			userID:         int64(1),
			reminderType:   model.OnceYearType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
		},
		{
			userID:         int64(1),
			reminderType:   model.DateType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
		},
	}

	for _, tt := range tests {
		n := New(nil)
		n.SaveName(tt.userID, tt.reminderName)
		err := n.SaveType(tt.userID, tt.reminderType)
		assert.NoError(t, err)

		location, err := time.LoadLocation(tt.timezoneString)
		if err != nil {
			t.Errorf("error while loading location: %v", err)
		}

		err = n.SaveCreatedField(tt.userID, location)
		assert.NoError(t, err)

		result, err := n.GetFromMemory(tt.userID)
		assert.NoError(t, err)

		assert.Equal(t, tt.userID, result.TgID)
		assert.Equal(t, tt.reminderName, result.Name)
		assert.Equal(t, tt.reminderType, result.Type)
		assert.Equal(t, location, result.Created.Location(), fmt.Sprintf("locations not equal: expected %s. Actual: %s", location, result.Created.Location()))
	}
}

func TestGetFromMemory_NameTypeCreatedDate(t *testing.T) {
	type test struct {
		userID         int64
		reminderType   model.ReminderType
		reminderName   string
		timezoneString string
		dateString     string
	}

	tests := []test{
		{
			userID:         int64(1),
			reminderType:   model.SeveralTimesDayType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
			dateString:     random.String(10),
		},
		{
			userID:         int64(1),
			reminderType:   model.EverydayType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
			dateString:     random.String(10),
		},
		{
			userID:         int64(1),
			reminderType:   model.EveryWeekType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
			dateString:     random.String(10),
		},
		{
			userID:         int64(1),
			reminderType:   model.SeveralDaysType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
			dateString:     random.String(10),
		},
		{
			userID:         int64(1),
			reminderType:   model.OnceMonthType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
			dateString:     random.String(10),
		},
		{
			userID:         int64(1),
			reminderType:   model.OnceYearType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
			dateString:     random.String(10),
		},
		{
			userID:         int64(1),
			reminderType:   model.DateType,
			reminderName:   random.String(10),
			timezoneString: "Europe/London",
			dateString:     random.String(10),
		},
	}

	for _, tt := range tests {
		n := New(nil)
		n.SaveName(tt.userID, tt.reminderName)
		err := n.SaveType(tt.userID, tt.reminderType)
		assert.NoError(t, err)

		location, err := time.LoadLocation(tt.timezoneString)
		if err != nil {
			t.Errorf("error while loading location: %v", err)
		}

		err = n.SaveCreatedField(tt.userID, location)
		assert.NoError(t, err)

		err = n.SaveDate(tt.userID, tt.dateString)
		assert.NoError(t, err)

		result, err := n.GetFromMemory(tt.userID)
		assert.NoError(t, err)

		assert.Equal(t, tt.userID, result.TgID)
		assert.Equal(t, tt.reminderName, result.Name)
		assert.Equal(t, tt.reminderType, result.Type)
		assert.Equal(t, location, result.Created.Location(), fmt.Sprintf("locations not equal: expected %s. Actual: %s", location, result.Created.Location()))
		assert.Equal(t, tt.dateString, result.Date)
	}
}

func TestGetFromMemory_NotFound(t *testing.T) {
	n := New(nil)

	_, err := n.GetFromMemory(int64(1))
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")
}

func TestSaveID(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)
	reminderID := uuid.New()

	n := New(nil)

	n.SaveName(userID, reminderName)

	err := n.SaveID(userID, reminderID)
	assert.NoError(t, err)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, result.TgID)
	assert.Equal(t, reminderName, result.Name)
	assert.Equal(t, reminderID, result.ID)
}

func TestSaveID_NotFound(t *testing.T) {
	userID := int64(1)
	reminderID := uuid.New()

	n := New(nil)

	err := n.SaveID(userID, reminderID)
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")

	_, ok := n.reminderMap[userID]
	assert.Equal(t, false, ok)

}

func TestGetID(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)
	reminderID := uuid.New()

	n := New(nil)

	n.SaveName(userID, reminderName)

	err := n.SaveID(userID, reminderID)
	assert.NoError(t, err)

	id, err := n.GetID(userID)
	assert.NoError(t, err)

	assert.Equal(t, reminderID, id)
}

func TestGetID_NotFound(t *testing.T) {
	userID := int64(1)

	n := New(nil)

	_, err := n.GetID(userID)
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")
}

func TestProcessMinutes(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)

	n := New(nil)

	n.SaveName(userID, reminderName)

	for i := 1; i < 60; i++ {
		minutesStr := strconv.Itoa(i)
		err := n.ProcessMinutes(userID, minutesStr)
		assert.NoError(t, err, fmt.Sprintf("case: %d", i))

		result, ok := n.reminderMap[userID]
		assert.Equal(t, true, ok)

		assert.Equal(t, userID, result.TgID)
		assert.Equal(t, minutesStr, result.Time)
	}
}

func TestProcessMinutes_UserNotFound(t *testing.T) {
	userID := int64(1)

	n := New(nil)

	err := n.ProcessMinutes(userID, "10")
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")

	_, ok := n.reminderMap[userID]
	assert.Equal(t, false, ok)
}

func TestProcessMinutes_RandomString(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)

	n := New(nil)

	n.SaveName(userID, reminderName)

	minutes := random.String(5)
	err := n.ProcessMinutes(userID, minutes)
	assert.Error(t, err)
}

func TestProcessMinutes_OutOfRange(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)

	n := New(nil)

	n.SaveName(userID, reminderName)

	minutes := []string{"-1", "0", "60"}

	for _, min := range minutes {
		err := n.ProcessMinutes(userID, min)
		assert.Error(t, err)
	}
}

func TestProcessHours(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)

	n := New(nil)

	n.SaveName(userID, reminderName)

	for i := 1; i < 25; i++ {
		hour := strconv.Itoa(i)
		err := n.ProcessHours(userID, hour)
		assert.NoError(t, err, fmt.Sprintf("case: %d", i))

		result, ok := n.reminderMap[userID]
		assert.Equal(t, true, ok)

		assert.Equal(t, userID, result.TgID)
		assert.Equal(t, hour, result.Time)
	}
}

func TestProcessHours_UserNotFound(t *testing.T) {
	userID := int64(1)

	n := New(nil)

	err := n.ProcessHours(userID, "10")
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")

	_, ok := n.reminderMap[userID]
	assert.Equal(t, false, ok)
}

func TestProcessHours_RandomString(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)

	n := New(nil)

	n.SaveName(userID, reminderName)

	hour := random.String(5)
	err := n.ProcessHours(userID, hour)
	assert.Error(t, err)
}

func TestProcessHours_OutOfRange(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)

	n := New(nil)

	n.SaveName(userID, reminderName)

	hours := []string{"-1", "0", "25"}

	for _, h := range hours {
		err := n.ProcessHours(userID, h)
		assert.Error(t, err, fmt.Sprintf("case: %s", h))
	}
}

func TestProcessDaysInMonth_OutOfRange(t *testing.T) {
	type test struct {
		name   string
		days   string
		userID int64
	}

	tests := []test{
		{
			name:   "1",
			days:   "-1",
			userID: 1,
		},
		{
			name:   "2",
			days:   "0",
			userID: 1,
		},
		{
			name:   "3",
			days:   "-10",
			userID: 1,
		},
		{
			name:   "4",
			days:   "100",
			userID: 1,
		},
		{
			name:   "5",
			days:   "32",
			userID: 1,
		},
	}

	for _, tt := range tests {
		n := New(nil)

		n.SaveName(tt.userID, random.String(5))

		err := n.ProcessDaysInMonth(tt.userID, tt.days)
		assert.EqualError(t, err, api_errors.ErrInvalidDays.Error())
	}
}

func TestProcessDaysInMonth_NotInteger(t *testing.T) {
	n := New(nil)

	n.SaveName(1, random.String(5))

	err := n.ProcessDaysInMonth(1, random.String(5))
	assert.Error(t, err)
}

func TestProcessDaysInMonth_Valid(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	n.SaveName(userID, random.String(5))

	daysInt := random.Int(1, 31)
	days := strconv.Itoa(daysInt)

	err := n.ProcessDaysInMonth(userID, days)
	assert.NoError(t, err)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, result.TgID)
	assert.Equal(t, days, result.Date)
}

func TestProcessDaysDuration_OutOfRange(t *testing.T) {
	type test struct {
		name   string
		days   string
		userID int64
	}

	tests := []test{
		{
			name:   "1",
			days:   "-1",
			userID: 1,
		},
		{
			name:   "2",
			days:   "0",
			userID: 1,
		},
		{
			name:   "3",
			days:   "-10",
			userID: 1,
		},
		{
			name:   "4",
			days:   "199",
			userID: 1,
		},
		{
			name:   "5",
			days:   "181",
			userID: 1,
		},
	}

	for _, tt := range tests {
		n := New(nil)

		n.SaveName(tt.userID, random.String(5))

		err := n.ProcessDaysDuration(tt.userID, tt.days)
		assert.EqualError(t, err, api_errors.ErrInvalidDays.Error())
	}
}

func TestProcessDaysDuration_NotInteger(t *testing.T) {
	n := New(nil)

	n.SaveName(1, random.String(5))

	err := n.ProcessDaysDuration(1, random.String(5))
	assert.Error(t, err)
}

func TestProcessDaysDuration_Valid(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	n.SaveName(userID, random.String(5))

	daysInt := random.Int(1, 180)
	days := strconv.Itoa(daysInt)

	err := n.ProcessDaysDuration(userID, days)
	assert.NoError(t, err)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, result.TgID)
	assert.Equal(t, days, result.Date)
}

func TestSaveCalendarDate_NotFound(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	err := n.SaveCalendarDate(userID, "")
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")
}

func TestSaveCalendarDate_Valid_OnceYearType(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	r := model.Reminder{
		Name: random.String(5),
		TgID: userID,
		Type: model.OnceYearType,
	}

	n.SaveName(r.TgID, r.Name)
	n.SaveType(r.TgID, r.Type)

	n.viewsMap[userID] = view.NewReminder()

	n.SetupCalendar(userID)

	n.viewsMap[userID].Calendar()

	err := n.SaveCalendarDate(userID, "12")
	assert.NoError(t, err)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	var curDate string

	month := time.Now().Month()

	if month < 10 {
		curDate = "12" + ".0" + strconv.Itoa(int(month))
	} else {
		curDate = "12" + strconv.Itoa(int(month))
	}

	assert.Equal(t, r.TgID, result.TgID)
	assert.Equal(t, curDate, result.Date)
}

func TestSaveCalendarDate_Valid_DateType(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	r := model.Reminder{
		Name: random.String(5),
		TgID: userID,
		Type: model.DateType,
	}

	n.SaveName(r.TgID, r.Name)
	n.SaveType(r.TgID, r.Type)

	n.viewsMap[userID] = view.NewReminder()

	n.SetupCalendar(userID)

	n.viewsMap[userID].Calendar()

	err := n.SaveCalendarDate(userID, "12")
	assert.NoError(t, err)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	var curDate string

	month := time.Now().Month()
	year := time.Now().Year()

	if month < 10 {
		curDate = fmt.Sprintf("%s.0%d.%d", "12", month, year)
	} else {
		curDate = fmt.Sprintf("%s.%d.%d", "12", month, year)
	}

	assert.Equal(t, r.TgID, result.TgID)
	assert.Equal(t, curDate, result.Date)
}

func TestCheckFields_EmptyTgID(t *testing.T) {
	n := New(nil)

	reminder := random.Reminder()
	reminder.TgID = 0

	err := n.checkFields(&reminder)
	assert.EqualError(t, err, "field TgID is not filled")
}

func TestCheckFields_EmptyName(t *testing.T) {
	n := New(nil)

	reminder := random.Reminder()
	reminder.Name = ""

	err := n.checkFields(&reminder)
	assert.EqualError(t, err, "field Name is not filled")
}

func TestCheckFields_EmptyType(t *testing.T) {
	n := New(nil)

	reminder := random.Reminder()
	reminder.Type = ""

	err := n.checkFields(&reminder)
	assert.EqualError(t, err, "field Type is not filled")
}

func TestCheckFields_EmptyDate(t *testing.T) {
	n := New(nil)

	reminder := random.Reminder()
	reminder.Date = ""

	err := n.checkFields(&reminder)
	assert.EqualError(t, err, "field Date is not filled")
}

func TestCheckFields_EmptyTime(t *testing.T) {
	n := New(nil)

	reminder := random.Reminder()
	reminder.Time = ""

	err := n.checkFields(&reminder)
	assert.EqualError(t, err, "field Time is not filled")
}

func TestCheckFields_EmptyCreated(t *testing.T) {
	n := New(nil)

	reminder := random.Reminder()
	reminder.Created = time.Time{}

	err := n.checkFields(&reminder)
	assert.EqualError(t, err, "field Created is not filled")
}

func TestFixMonth_Before10(t *testing.T) {
	months := []time.Month{
		time.January,
		time.February,
		time.March,
		time.April,
		time.May,
		time.June,
		time.July,
		time.August,
		time.September,
	}

	for _, month := range months {
		expected := "0" + strconv.Itoa(int(month))

		actual := fixMonth(month)

		assert.Equal(t, expected, actual)
	}
}

func TestFixMonth_After10(t *testing.T) {
	months := []time.Month{
		time.October,
		time.November,
		time.December,
	}

	for _, month := range months {
		expected := strconv.Itoa(int(month))

		actual := fixMonth(month)

		assert.Equal(t, expected, actual)
	}
}

func TestValidateDate_Valid(t *testing.T) {
	userID := int64(1)

	curDay := time.Now().Day() + 1
	dayOfMonth := strconv.Itoa(curDay)

	tz := time.Local

	n := New(nil)

	n.viewsMap[userID] = view.NewReminder()

	n.SetupCalendar(userID)

	err := n.ValidateDate(userID, dayOfMonth, tz)
	assert.NoError(t, err)
}

func TestValidateDate_NotInt(t *testing.T) {
	userID := int64(1)

	dayOfMonth := random.String(4)

	tz := time.Local

	n := New(nil)

	n.viewsMap[userID] = view.NewReminder()

	n.SetupCalendar(userID)

	err := n.ValidateDate(userID, dayOfMonth, tz)
	assert.Error(t, err)
}

func TestValidateDate_DatePassed(t *testing.T) {
	userID := int64(1)

	curDay := time.Now().Day() - 1
	dayOfMonth := strconv.Itoa(curDay)

	tz := time.Local

	n := New(nil)

	n.viewsMap[userID] = view.NewReminder()

	n.SetupCalendar(userID)

	err := n.ValidateDate(userID, dayOfMonth, tz)
	assert.EqualError(t, err, api_errors.ErrInvalidDate.Error())
}

func TestValidateTime_Valid(t *testing.T) {
	userID := int64(1)

	n := New(nil)

	n.SaveName(userID, random.String(10))

	year, month, day := time.Now().Date()

	var monthStr string
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	} else {
		monthStr = fmt.Sprintf("%d", month)
	}

	date := fmt.Sprintf("%d.%s.%d", day, monthStr, year)

	userDateWithTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", date, "23:59"))
	require.NoError(t, err)

	err = n.ValidateTime(time.Local, userDateWithTime)
	assert.NoError(t, err)
}

func TestValidateTime_Invalid(t *testing.T) {
	userID := int64(1)

	n := New(nil)

	n.SaveName(userID, random.String(10))

	year, month, day := time.Now().Date()

	var monthStr string
	if month < 10 {
		monthStr = fmt.Sprintf("0%d", month)
	} else {
		monthStr = fmt.Sprintf("%d", month)
	}

	date := fmt.Sprintf("%d.%s.%d", day, monthStr, year)

	userDateWithTime, err := time.Parse("02.01.2006 15:04", fmt.Sprintf("%s %s", date, "00:00"))
	require.NoError(t, err)

	err = n.ValidateTime(time.Local, userDateWithTime)
	assert.EqualError(t, err, api_errors.ErrTimeInPast.Error())
}
