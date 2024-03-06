package view

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

type ReminderView struct {
	pages       []string
	currentPage int
	logger      *logrus.Logger
	calendar    *Calendar
}

func NewReminder() *ReminderView {
	return &ReminderView{pages: make([]string, 0), currentPage: 0, logger: logger.New(), calendar: new()}
}

var (
	// inline кнопка для переключения на предыдущую страницу (напоминания)
	BtnPrevPgReminders = tele.Btn{Text: "<", Unique: "prev_pg_reminders"}
	// inline кнопка для переключения на следующую страницу (напоминания)
	BtnNextPgReminders = tele.Btn{Text: ">", Unique: "next_pg_reminders"}

	// inline кнопка для переключения на первую страницу (напоминания)
	BtnFirstPgReminders = tele.Btn{Text: "<<", Unique: "start_pg_reminders"}
	// inline кнопка для переключения на последнюю страницу (напоминания)
	BtnLastPgReminders = tele.Btn{Text: ">>", Unique: "end_pg_reminders"}
)

// Message формирует список сообщений из моделей заметок и возвращает первую страницу.
// Количество заметок на одной странице задает переменная noteCountPerPage (по умолчанию - 5)
func (v *ReminderView) Message(reminders []model.Reminder) (string, error) {
	// if len(reminders) == 0 {
	// 	return messages.UserDoesntHaveNotesMessage, nil
	// }

	var res = ""

	v.pages = make([]string, 0)

	for i, reminder := range reminders {
		txt, err := ProcessTypeAndDate(reminder.Type, reminder.Date, reminder.Time)
		if err != nil {
			return "", err
		}

		res += fmt.Sprintf("<b>%d. %s</b>\n\nСрабатывает: %s\nСоздано: %s\nУдалить: /del%d\n\n", i+1, reminder.Name, txt, reminder.Created.Format(createdFieldFormat), reminder.ID)
		if i%noteCountPerPage == 0 && i > 0 || len(res) == maxMessageLen {
			v.pages = append(v.pages, res)
			res = ""
		}
	}

	if len(v.pages) < 5 && res != "" {
		v.pages = append(v.pages, res)
	}

	v.currentPage = 0

	return v.pages[0], nil
}

// Next возвращает следующую страницу сообщений
func (v *ReminderView) Next() string {
	v.logger.Debugf("ReminderView: getting next page. Current: %d\n", v.currentPage)

	if v.currentPage == v.total()-1 {
		v.logger.Debugf("ReminderView: current page is the last. Setting current page to 0.\n")
		v.currentPage = 0
	} else {
		v.currentPage++
		v.logger.Debugf("ReminderView: incrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Previous возвращает предыдущую страницу сообщений
func (v *ReminderView) Previous() string {
	v.logger.Debugf("ReminderView: getting previous page. Current: %d\n", v.currentPage)

	if v.currentPage == 0 {
		v.logger.Debugf("ReminderView: previous page is the last. Setting current page to maximum: %d.\n", v.total())
		v.currentPage = v.total() - 1
	} else {
		v.currentPage--
		v.logger.Debugf("ReminderView: decrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Last возвращает последнюю страницу сообщений
func (v *ReminderView) Last() string {
	v.logger.Debugf("ReminderView: getting the last page. Current: %d\n", v.currentPage)

	v.currentPage = v.total() - 1

	return v.pages[v.currentPage]
}

// First возвращает первую страницу сообщений
func (v *ReminderView) First() string {
	v.logger.Debugf("ReminderView: getting the first page. Current: %d\n", v.currentPage)

	v.currentPage = 0

	return v.pages[v.currentPage]
}

// current возвращает номер текущей страницы
func (v *ReminderView) current() int {
	return v.currentPage + 1
}

// total возвращает общее количество страниц
func (v *ReminderView) total() int {
	return len(v.pages)
}

// Keyboard делает клавиатуру для навигации по страницам
func (v *ReminderView) Keyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	// если страниц 1, клавиатура не нужна
	if v.total() == 1 {
		menu.Inline(
			menu.Row(BtnCreateReminder),
			menu.Row(BtnDeleteAllReminders),
			menu.Row(BtnBackToMenu),
		)
		return menu
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := menu.Data(text, "s")

	menu.Inline(
		menu.Row(BtnFirstPgReminders, BtnPrevPgReminders, btn, BtnNextPgReminders, BtnLastPgReminders),
		menu.Row(BtnCreateReminder),
		menu.Row(BtnDeleteAllReminders),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// SetCurrentToFirst устанавливает текущий номер страницы на 1
func (v *ReminderView) SetCurrentToFirst() {
	v.currentPage = 0
}

// Clear используется когда удаляются все заметки: очищает список заметок, устанавливает текущую страницу в 0
func (v *ReminderView) Clear() {
	v.currentPage = 0
	v.pages = make([]string, 0)
}

// ReminderMessage возвращает текст сообщения с напоминанием.
// Пример:
//
// купить хлеб
//
// Напоминание сработало 23.10.2023 в 18:00
func ReminderMessage(reminder model.Reminder) (string, error) {
	name := reminder.Name

	date, err := ProcessTypeAndDate(reminder.Type, reminder.Date, reminder.Time)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(messages.ReminderMessage, name, date), nil
}

// ProcessTypeAndDate обрабатывает тип напоминания и дату. Пример:
//
// everyday 11:30 -> ежедневно в 11:30
//
// SeveralTimesDayType, minutes, 1 -> раз в 1 минуту
func ProcessTypeAndDate(reminderType model.ReminderType, date, time string) (string, error) {
	switch reminderType {
	case model.EverydayType:
		return fmt.Sprintf("ежедневно в %s", time), nil
	case model.SeveralTimesDayType:
		if date == "minutes" {
			minutesInt, _ := strconv.Atoi(time) // опускаем ошибку, потому что время уже было проавлидировано на предыдущих шагах
			return fmt.Sprintf("один раз в %s", processMinutes(time, minutesInt)), nil
		} else {
			hoursInt, _ := strconv.Atoi(time)
			return fmt.Sprintf("один раз в %s", processHours(time, hoursInt)), nil
		}
	case model.EveryWeekType:
		txt, err := processWeekDay(date, time)
		if err != nil {
			return "", fmt.Errorf("error while processing week day: %w", err)
		}

		return txt, nil
	case model.SeveralDaysType:
		txt, err := processDays(date, time)
		if err != nil {
			return "", fmt.Errorf("error while processing days interval: %w", err)
		}

		return txt, nil
	case model.OnceMonthType:
		return fmt.Sprintf("каждый месяц %s числа в %s", date, time), nil
	case model.OnceYearType:
		return processDateWithoutYear(date, time), nil
	case model.DateType:
		return processDateWithYear(date, time), nil
	default:
		return "", fmt.Errorf("unknown reminder type: %s", reminderType)
	}
}

func processDateWithoutYear(date, time string) string {
	dates := strings.Split(date, ".")

	day := dates[0]
	month := dates[1]

	monthStr := processMonth(month)

	return fmt.Sprintf("раз в год %s %s в %s", day, monthStr, time)
}

func processDateWithYear(date, time string) string {
	dates := strings.Split(date, ".")

	day := dates[0]
	month := dates[1]
	year := dates[2]

	monthStr := processMonth(month)

	fullDate := fmt.Sprintf("%s %s %s года", day, monthStr, year)

	return fmt.Sprintf("%s в %s", fullDate, time)
}

func processMonth(month string) string {
	monthsMap := map[string]string{
		"01": "января",
		"02": "февраля",
		"03": "марта",
		"04": "апреля",
		"05": "мая",
		"06": "июня",
		"07": "июля",
		"08": "августа",
		"09": "сентября",
		"10": "октября",
		"11": "ноября",
		"12": "декабря",
	}

	return monthsMap[month]
}

func processDays(days, userTime string) (string, error) {
	daysInt, err := strconv.Atoi(days)
	if err != nil {
		return "", fmt.Errorf("error while converting string %s to int: %w", days, err)
	}

	// раз в день
	if daysInt == 1 {
		return fmt.Sprintf("один раз в день в %s", userTime), nil
	}

	// > 20
	if daysInt > 20 {
		// раз в 21, 31 день
		if endsWith(days, "1") {
			return fmt.Sprintf("один раз в %d день в %s", daysInt, userTime), nil
		}

		// раз в 22 дня, 32 дня
		if endsWith(days, "2") {
			return fmt.Sprintf("один раз в %d дня в %s", daysInt, userTime), nil
		}
	}

	// < 10
	if daysInt < 10 {
		// раз в 2, 3, 4 дня
		if endsWith(days, "2", "3", "4") {
			return fmt.Sprintf("один раз в %d дня в %s", daysInt, userTime), nil
		}
	}

	// 5 - 20
	return fmt.Sprintf("один раз в %d дней в %s", daysInt, userTime), nil
}

func processWeekDay(date, userTime string) (string, error) {
	wd, err := parseWeekdayRus(date)
	if err != nil {
		return "", fmt.Errorf("error while translating week day %s: %w", date, err)
	}

	switch wd {
	case "понедельник":
		return fmt.Sprintf("каждый %s в %s", wd, userTime), nil
	case "вторник":
		return fmt.Sprintf("каждый %s в %s", wd, userTime), nil
	case "среда":
		return fmt.Sprintf("каждую среду в %s", userTime), nil
	case "четверг":
		return fmt.Sprintf("каждый %s в %s", wd, userTime), nil
	case "пятница":
		return fmt.Sprintf("каждую пятницу в %s", userTime), nil
	case "суббота":
		return fmt.Sprintf("каждую субботу в %s", userTime), nil
	case "воскресенье":
		return fmt.Sprintf("каждое %s в %s", wd, userTime), nil
	default:
		return "", fmt.Errorf("unknown week day: %s", wd)
	}
}

func parseWeekdayRus(v string) (string, error) {
	daysOfWeek := map[string]string{
		"sunday":    "воскресенье",
		"monday":    "понедельник",
		"tuesday":   "вторник",
		"wednesday": "среда",
		"thursday":  "четверг",
		"friday":    "пятница",
		"saturday":  "суббота",
	}

	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return "", fmt.Errorf("invalid weekday '%s'", v)
}

func ParseWeekday(v string) (time.Weekday, error) {
	daysOfWeek := map[string]time.Weekday{
		"sunday":    time.Sunday,
		"monday":    time.Monday,
		"tuesday":   time.Tuesday,
		"wednesday": time.Wednesday,
		"thursday":  time.Thursday,
		"friday":    time.Friday,
		"saturday":  time.Saturday,
	}

	if d, ok := daysOfWeek[v]; ok {
		return d, nil
	}

	return time.Sunday, fmt.Errorf("invalid weekday '%s'", v)
}

func processHours(hoursString string, hoursInt int) string {
	// 1 час
	// 5...20 - часов
	// 2..4 - часа

	if hoursInt == 1 {
		return "час"
	}

	// 21
	if hoursInt > 20 && endsWith(hoursString, "1") {
		return fmt.Sprintf("%d час", hoursInt)
	}

	// 5...20
	if hoursInt >= 5 && hoursInt <= 20 {
		return fmt.Sprintf("%d часов", hoursInt)
	}

	// [2..4], [22, 23, 24] - часа
	return fmt.Sprintf("%d часа", hoursInt)
}

func processMinutes(minutesString string, minutesInt int) string {
	if minutesInt == 1 {
		return "минуту"
	}

	if minutesInt >= 20 && endsWith(minutesString, "1") { // [21, ...]
		return fmt.Sprintf("%d минуту", minutesInt) // 21, 31...

	}

	// 2, 3, 4, [22, 23, 24, 32, 33, 34, ...]
	if endsWith(minutesString, "2", "3", "4") && (minutesInt < 10 || minutesInt > 20) {
		return fmt.Sprintf("%d минуты", minutesInt)
	}

	// 5-9, 20, 25, 35, 26, 27...
	if endsWith(minutesString, "0", "5", "6", "7", "8", "9") || (minutesInt >= 10 && minutesInt < 20) {
		return fmt.Sprintf("%d минут", minutesInt)
	}

	return ""
}

// endsWith проверяет, оканчивается ли строка на один из суффиксов
func endsWith(s string, suff ...string) bool {
	for _, suf := range suff {
		if strings.HasSuffix(s, suf) {
			return true
		}
	}

	return false
}

// Calendar возвращает календарь с текущим месяцем и годом
func (v *ReminderView) Calendar() *tele.ReplyMarkup {
	calendar := v.calendar.currentCalendar()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToReminderType)
	return calendar
}

// PrevMonth возвращает календарь с предыдущим месяцем
func (v *ReminderView) PrevMonth() *tele.ReplyMarkup {
	calendar := v.calendar.prevMonth()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToReminderType)
	return calendar
}

// NextMonth возвращает календарь со следующим месяцем
func (v *ReminderView) NextMonth() *tele.ReplyMarkup {
	calendar := v.calendar.nextMonth()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToReminderType)
	return calendar
}

// PrevYear возвращает календарь с предыдущим годом
func (v *ReminderView) PrevYear() *tele.ReplyMarkup {
	calendar := v.calendar.prevYear()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToReminderType)
	return calendar
}

// NextYear возвращает календарь с следующим годом
func (v *ReminderView) NextYear() *tele.ReplyMarkup {
	calendar := v.calendar.nextYear()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToReminderType)
	return calendar
}

// GetDaysBtns возвращает слайс кнопок с числами месяца
func (v *ReminderView) GetDaysBtns() []tele.Btn {
	return v.calendar.getDaysBtns()
}

// Month возвращает возвращает месяц, установленный в календаре на данный момент
func (v *ReminderView) Month() time.Month {
	return v.calendar.month()
}

// Year возвращает возвращает год, установленный в календаре на данный момент
func (v *ReminderView) Year() int {
	return v.calendar.year()
}

// SetCurMonth устаналивает месяц в текущий
func (v *ReminderView) SetCurMonth() {
	v.calendar.setCurMonth()
}

// SetCurYear устаналивает год в текущий
func (v *ReminderView) SetCurYear() {
	v.calendar.setCurYear()
}
