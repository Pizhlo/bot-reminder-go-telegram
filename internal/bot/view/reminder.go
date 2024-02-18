package view

import (
	"fmt"
	"strconv"
	"strings"

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
}

func NewReminder() *ReminderView {
	return &ReminderView{pages: make([]string, 0), currentPage: 0, logger: logger.New()}
}

var (
	// inline кнопка для переключения на предыдущую страницу (напоминания)
	BtnPrevPgReminders = selector.Data("<", "prev")
	// inline кнопка для переключения на следующую страницу (напоминания)
	BtnNextPgReminders = selector.Data(">", "next")

	// inline кнопка для переключения на первую страницу (напоминания)
	BtnFirstPgReminders = selector.Data("<<", "start")
	// inline кнопка для переключения на последнюю страницу (напоминания)
	BtnLastPgReminders = selector.Data(">>", "end")
)

// Message формирует список сообщений из моделей заметок и возвращает первую страницу.
// Количество заметок на одной странице задает переменная noteCountPerPage (по умолчанию - 5)
func (v *ReminderView) Message(reminders []model.Reminder) string {
	if len(reminders) == 0 {
		return messages.UserDoesntHaveNotesMessage
	}

	var res = ""

	v.pages = make([]string, 0)

	for i, reminder := range reminders {
		res += fmt.Sprintf("%d. Создано: %s. Удалить: /del%d\n\n%s\n\n%s", i+1, reminder.Created.Format(dateFormat), reminder.ID, reminder.Name, string(reminder.Type)+reminder.Date+reminder.Time)
		if i%noteCountPerPage == 0 && i > 0 || len(res) == maxMessageLen {
			v.pages = append(v.pages, res)
			res = ""
		}
	}

	if len(v.pages) < 5 && res != "" {
		v.pages = append(v.pages, res)
	}

	return v.pages[0]
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
	// если страниц 1, клавиатура не нужна
	if v.total() == 1 {
		menu := &tele.ReplyMarkup{}
		menu.Inline(
			menu.Row(BtnCreateReminder),
			selector.Row(BtnDeleteAllReminders),
			menu.Row(BtnBackToMenu),
		)
		return menu
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := selector.Data(text, "s")

	selector.Inline(
		selector.Row(BtnFirstPgReminders, BtnPrevPgReminders, btn, BtnNextPgReminders, BtnLastPgReminders),
		selector.Row(BtnDeleteAllReminders),
		selector.Row(BtnBackToMenu),
	)

	return selector
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
func ReminderMessage(reminder model.Reminder) string {
	name := reminder.Name

	date := ProcessTypeAndDate(reminder.Type, reminder.Date, reminder.Time)

	return fmt.Sprintf(messages.ReminderMessage, name, date)
}

// ProcessTypeAndDate обрабатывает тип напоминания и дату. Пример:
//
// everyday 11:30 -> ежедневно в 11:30
//
// SeveralTimesDayType, minutes, 1 -> раз в 1 минуту
func ProcessTypeAndDate(reminderType model.ReminderType, date, time string) string {
	switch reminderType {
	case model.EverydayType:
		return fmt.Sprintf("ежедневно в %s", time)
	case model.SeveralTimesDayType:
		if date == "minutes" {
			minutesInt, _ := strconv.Atoi(time) // опускаем ошибку, потому что время уже было проавлидировано на предыдущих шагах
			return fmt.Sprintf("один раз в %s", processMinutes(time, minutesInt))
		}
	default:
		return ""
	}

	return ""
}

func processMinutes(minutesString string, minutesInt int) string {
	if minutesInt < 10 || minutesInt >= 20 { // [1, 9], [21, ...]
		if endsWith(minutesString, "1") { // 1, 21, 31...
			if minutesInt == 1 {
				return "минуту"
			} else {
				return fmt.Sprintf("%d минуту", minutesInt)
			}

		}

		// 2, 3, 4, [22, 23, 24, 32, 33, 34, ...]
		if endsWith(minutesString, "2", "3", "4") {
			return fmt.Sprintf("%d минуты", minutesInt)
		}

		// 5-9, 20, 25, 35, 26, 27...
		if endsWith(minutesString, "0", "5", "6", "7", "8", "9") {
			return fmt.Sprintf("%d минут", minutesInt)
		}
	}

	// [10, 19]
	if minutesInt >= 10 && minutesInt < 20 {
		return fmt.Sprintf("%d минут", minutesInt)
	}

	return ""

}

// endsWith проверяет, оканчивается ли строка на один из суффиксов
func endsWith(s string, suff ...string) bool {
	count := 0
	for _, suf := range suff {
		if strings.HasSuffix(s, suf) {
			count++
		}
	}

	return count > 0
}
