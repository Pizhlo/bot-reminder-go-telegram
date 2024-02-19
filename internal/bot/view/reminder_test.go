package view

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/stretchr/testify/assert"
)

func TestProcessHours_1_21(t *testing.T) {
	type test struct {
		hoursString string
		hoursInt    int
		expected    string
	}

	tests := []test{
		{
			hoursString: "1",
			hoursInt:    1,
			expected:    "час",
		},
		{
			hoursString: "21",
			hoursInt:    21,
			expected:    "21 час",
		},
	}

	for _, tt := range tests {
		actual := processHours(tt.hoursString, tt.hoursInt)
		assert.Equal(t, tt.expected, actual)
	}
}

func TestProcessHours_5_20(t *testing.T) {
	hoursMap := map[int]string{}
	expected := map[int]string{}

	for i := 5; i < 21; i++ {
		hoursMap[i] = strconv.Itoa(i)
		expected[i] = fmt.Sprintf("%d часов", i)
	}

	for k, v := range hoursMap {
		actual := processHours(v, k)
		assert.Equal(t, expected[k], actual)
	}
}

func TestProcessHours_2_4(t *testing.T) {
	hoursMap := map[int]string{}
	expected := map[int]string{}

	for i := 2; i < 5; i++ {
		hoursMap[i] = strconv.Itoa(i)
		expected[i] = fmt.Sprintf("%d часа", i)
	}

	for i := 22; i < 25; i++ {
		hoursMap[i] = strconv.Itoa(i)
		expected[i] = fmt.Sprintf("%d часа", i)
	}

	for k, v := range hoursMap {
		actual := processHours(v, k)
		assert.Equal(t, expected[k], actual)
	}
}

func TestProcessMinutes_From20(t *testing.T) {
	minutesMap := map[int]string{}
	expected := map[int]string{}

	expected[1] = "минуту"

	for i := 20; i < 59; i++ {
		str := strconv.Itoa(i)
		if strings.HasSuffix(str, "1") {
			minutesMap[i] = str
			expected[i] = fmt.Sprintf("%d минуту", i)
		}

	}

	for k, v := range minutesMap {
		actual := processMinutes(v, k)
		assert.Equal(t, expected[k], actual)
	}
}

func TestProcessMinutes_2_3_4(t *testing.T) {
	minutesMap := map[int]string{}
	expected := map[int]string{}

	for i := 2; i < 5; i++ {
		minutesMap[i] = strconv.Itoa(i)
		expected[i] = fmt.Sprintf("%d минуты", i)
	}

	for i := 22; i < 59; i++ {
		str := strconv.Itoa(i)
		if strings.HasSuffix(str, "2") || strings.HasSuffix(str, "3") || strings.HasSuffix(str, "4") {
			minutesMap[i] = strconv.Itoa(i)
			expected[i] = fmt.Sprintf("%d минуты", i)
		}

	}

	for k, v := range minutesMap {
		actual := processMinutes(v, k)
		assert.Equal(t, expected[k], actual)
	}
}

func TestProcessMinutes_EndWith234(t *testing.T) {
	minutesMap := map[int]string{}
	expected := map[int]string{}

	for i := 5; i < 10; i++ {
		minutesMap[i] = strconv.Itoa(i)
		expected[i] = fmt.Sprintf("%d минут", i)
	}

	for i := 22; i < 59; i++ {
		str := strconv.Itoa(i)
		if strings.HasSuffix(str, "0") || strings.HasSuffix(str, "5") || strings.HasSuffix(str, "6") || strings.HasSuffix(str, "7") || strings.HasSuffix(str, "8") || strings.HasSuffix(str, "9") {
			minutesMap[i] = strconv.Itoa(i)
			expected[i] = fmt.Sprintf("%d минут", i)
		}
	}

	for i := 10; i < 21; i++ {
		minutesMap[i] = strconv.Itoa(i)
		expected[i] = fmt.Sprintf("%d минут", i)
	}

	for k, v := range minutesMap {
		actual := processMinutes(v, k)
		assert.Equal(t, expected[k], actual)
	}
}

func TestProcessMinutes(t *testing.T) {
	actual := processMinutes("1", 1)
	assert.Equal(t, "минуту", actual)
}

func TestEndsWith_True(t *testing.T) {
	type test struct {
		s        string
		suff     []string
		expected bool
	}

	tests := []test{
		{
			s:        "21",
			suff:     []string{"1"},
			expected: true,
		},
		{
			s:        "22",
			suff:     []string{"2", "3", "4"},
			expected: true,
		},
		{
			s:        "23",
			suff:     []string{"2", "3", "4"},
			expected: true,
		},
		{
			s:        "24",
			suff:     []string{"2", "3", "4"},
			expected: true,
		},
		{
			s:        "33",
			suff:     []string{"2", "3", "4"},
			expected: true,
		},
		{
			s:        "5",
			suff:     []string{"0", "5", "6", "7", "8", "9"},
			expected: true,
		},
		{
			s:        "10",
			suff:     []string{"0", "5", "6", "7", "8", "9"},
			expected: true,
		},
		{
			s:        "20",
			suff:     []string{"0", "5", "6", "7", "8", "9"},
			expected: true,
		},
		{
			s:        "27",
			suff:     []string{"0", "5", "6", "7", "8", "9"},
			expected: true,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, endsWith(tt.s, tt.suff...), fmt.Sprintf("case: %s", tt.s))
	}
}

func TestEndsWith_False(t *testing.T) {
	suff := []string{}
	for i := 1; i < 8; i++ {
		suff = append(suff, strconv.Itoa(i))
	}

	s := []string{"19", "29", "39"}

	for _, test := range s {
		assert.Equal(t, false, endsWith(test, suff...))
	}

	suff = []string{}
	for i := 2; i < 5; i++ {
		suff = append(suff, strconv.Itoa(i))
	}

	s = []string{"11", "28", "27"}

	for _, test := range s {
		assert.Equal(t, false, endsWith(test, suff...))
	}
}

func TestParseWeekday(t *testing.T) {
	daysOfWeek := map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}

	for k, v := range daysOfWeek {
		actual, err := ParseWeekday(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestParseWeekday_Unknown(t *testing.T) {
	str := random.String(10)
	expectedErr := fmt.Sprintf("invalid weekday '%s'", str)
	_, err := ParseWeekday(str)
	assert.EqualError(t, err, expectedErr)
}

func TestParseWeekdayRus(t *testing.T) {
	daysOfWeek := map[string]string{
		"sunday":    "воскресенье",
		"monday":    "понедельник",
		"tuesday":   "вторник",
		"wednesday": "среда",
		"thursday":  "четверг",
		"friday":    "пятница",
		"saturday":  "суббота",
	}

	for k, v := range daysOfWeek {
		actual, err := parseWeekdayRus(k)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestParseWeekdayRus_Unknown(t *testing.T) {
	str := random.String(10)
	expectedErr := fmt.Sprintf("invalid weekday '%s'", str)
	_, err := parseWeekdayRus(str)
	assert.EqualError(t, err, expectedErr)
}

func TestProcessWeekDay(t *testing.T) {
	userTime := "10:10"

	expected := map[string]string{
		"monday":    "каждый понедельник в 10:10",
		"tuesday":   "каждый вторник в 10:10",
		"wednesday": "каждую среду в 10:10",
		"thursday":  "каждый четверг в 10:10",
		"friday":    "каждую пятницу в 10:10",
		"saturday":  "каждую субботу в 10:10",
		"sunday":    "каждое воскресенье в 10:10",
	}

	for k, v := range expected {
		actual, err := processWeekDay(k, userTime)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestProcessDays_Invalid(t *testing.T) {
	_, err := processDays(random.String(10), random.String(5))
	assert.Error(t, err)
}

func TestProcessDays_1(t *testing.T) {
	actual, err := processDays("1", "10:10")
	assert.NoError(t, err)
	assert.Equal(t, "один раз в день в 10:10", actual)
}

func TestProcessDays_21_31_22_32(t *testing.T) {
	expected := map[string]string{
		"21": "один раз в 21 день в 10:10",
		"31": "один раз в 31 день в 10:10",
		"22": "один раз в 22 дня в 10:10",
		"32": "один раз в 32 дня в 10:10",
	}

	userTime := "10:10"

	for k, v := range expected {
		actual, err := processDays(k, userTime)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestProcessDays_2_3_4(t *testing.T) {
	expected := map[string]string{
		"2": "один раз в 2 дня в 10:10",
		"3": "один раз в 3 дня в 10:10",
		"4": "один раз в 4 дня в 10:10",
	}

	userTime := "10:10"

	for k, v := range expected {
		actual, err := processDays(k, userTime)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}

func TestProcessDays_From_5_To_20(t *testing.T) {
	expected := map[string]string{}
	userTime := "10:10"

	for i := 5; i < 21; i++ {
		str := strconv.Itoa(i)
		expected[str] = fmt.Sprintf("один раз в %s дней в %s", str, userTime)
	}

	for k, v := range expected {
		actual, err := processDays(k, userTime)
		assert.NoError(t, err)
		assert.Equal(t, v, actual)
	}
}
