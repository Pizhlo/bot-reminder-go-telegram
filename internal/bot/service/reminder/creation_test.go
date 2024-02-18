package reminder

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/stretchr/testify/assert"
)

func TestSaveReminderName(t *testing.T) {
	n := New(nil)

	userID := int64(1)
	reminderName := random.String(10)

	n.SaveName(userID, reminderName)

	result, ok := n.reminderMap[userID]
	assert.Equal(t, true, ok)

	assert.Equal(t, userID, result.TgID)
	assert.Equal(t, reminderName, result.Name)
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

		err := n.ProcessTime(tt.userID, tt.timeMsg)
		assert.NoError(t, err)

		result, ok := n.reminderMap[tt.userID]
		assert.Equal(t, true, ok)

		assert.Equal(t, tt.userID, result.TgID)
		assert.Equal(t, tt.timeMsg, result.Time)
		assert.Equal(t, tt.reminderName, result.Name)
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

		err := n.ProcessTime(tt.userID, tt.timeMsg)

		assert.Error(t, err)

		result, ok := n.reminderMap[tt.userID]
		assert.Equal(t, true, ok)

		assert.Equal(t, tt.reminderName, result.Name)

	}
}

func TestProcessTime_NotFound(t *testing.T) {
	n := New(nil)

	userID := int64(1)

	err := n.ProcessTime(userID, "10:10")
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")

	_, ok := n.reminderMap[userID]
	assert.Equal(t, false, ok)
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
	reminderID := int64(1)

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
	reminderID := int64(1)

	n := New(nil)

	err := n.SaveID(userID, reminderID)
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")

	_, ok := n.reminderMap[userID]
	assert.Equal(t, false, ok)

}

func TestGetID(t *testing.T) {
	userID := int64(1)
	reminderName := random.String(10)
	reminderID := int64(1)

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

	id, err := n.GetID(userID)
	assert.EqualError(t, err, "error while getting reminder by user ID: reminder not found")
	assert.Equal(t, int64(0), id)
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
