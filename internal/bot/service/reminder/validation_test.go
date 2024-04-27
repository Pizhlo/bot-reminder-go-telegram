package reminder

import (
	"strconv"
	"testing"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/stretchr/testify/assert"
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

	n := New(nil)

	l := time.FixedZone("Europe/Moscow", 0)

	userTime := time.Now().In(l).Add(24 * time.Hour)

	err := n.ValidateTime(time.Local, userTime)
	assert.NoError(t, err)
}

func TestValidateTime_Invalid(t *testing.T) {
	n := New(nil)

	l := time.FixedZone("Europe/Moscow", 0)

	userTime := time.Now().In(l).Add(-24 * time.Hour)

	err := n.ValidateTime(l, userTime)
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
