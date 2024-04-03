package reminder

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestFixInt_Before10(t *testing.T) {
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

		actual := fixInt(int(month))

		assert.Equal(t, expected, actual)
	}
}

func TestFixInt_After10(t *testing.T) {
	months := []time.Month{
		time.October,
		time.November,
		time.December,
	}

	for _, month := range months {
		expected := strconv.Itoa(int(month))

		actual := fixInt(int(month))

		assert.Equal(t, expected, actual)
	}
}
