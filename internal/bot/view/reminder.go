package view

import (
	"fmt"

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
	// inline кнопка для переключения на предыдущую страницу (заметки)
	BtnPrevPgReminders = selector.Data("<", "prev")
	// inline кнопка для переключения на следующую страницу (заметки)
	BtnNextPgReminders = selector.Data(">", "next")

	// inline кнопка для переключения на первую страницу (заметки)
	BtnFirstPgReminders = selector.Data("<<", "start")
	// inline кнопка для переключения на последнюю страницу (заметки)
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
		res += fmt.Sprintf("%d. Создано: %s. Удалить: /del%d\n\n%s\n\n%s", i+1, reminder.Created.Format(dateFormat), reminder.ID, reminder.Text, reminder.Type+reminder.Date+reminder.Time)
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
			menu.Row(BtnSearchNotesByText, BtnSearchNotesByDate),
			menu.Row(BtnDeleteAllNotes),
			menu.Row(BtnBackToMenu),
		)
		return menu
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := selector.Data(text, "")

	selector.Inline(
		selector.Row(BtnFirstPgNotes, BtnPrevPgNotes, btn, BtnNextPgNotes, BtnLastPgNotes),
		selector.Row(BtnSearchNotesByText, BtnSearchNotesByDate),
		selector.Row(BtnDeleteAllNotes),
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
