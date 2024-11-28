package view

import (
	"fmt"
	"time"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

const (
	createdFieldFormat = "02.01.2006 15:04:05"
	noteCountPerPage   = 5
	maxMessageLen      = 4096
)

type NoteView struct {
	pages       []string
	currentPage int
	calendar    *Calendar
}

func NewNote() *NoteView {
	return &NoteView{pages: make([]string, 0), currentPage: 0, calendar: new()}
}

var (
	// inline кнопка для переключения на предыдущую страницу (заметки)
	BtnPrevPgNotes = tele.Btn{Text: "<", Unique: "prev_pg_notes"}
	// inline кнопка для переключения на следующую страницу (заметки)
	BtnNextPgNotes = tele.Btn{Text: ">", Unique: "next_pg_notes"}

	// inline кнопка для обновления заметок
	BtnRefreshNotes = tele.Btn{Text: "🔁", Unique: "notes"}

	// inline кнопка для переключения на первую страницу (заметки)
	BtnFirstPgNotes = tele.Btn{Text: "<<", Unique: "start_pg_notes"}
	// inline кнопка для переключения на последнюю страницу (заметки)
	BtnLastPgNotes = tele.Btn{Text: ">>", Unique: "end_pg_notes"}
)

// Message формирует список сообщений из моделей заметок и возвращает первую страницу.
// Количество заметок на одной странице задает переменная noteCountPerPage (по умолчанию - 5)
func (v *NoteView) Message(notes []model.Note) (string, error) {
	var res = ""

	v.pages = make([]string, 0)

	for i, note := range notes {

		header := v.fillHeader(i+1, note)

		res += fmt.Sprintf("%s\n\n%s\n\n", header, note.Text)

		//res += fmt.Sprintf("<b>%d. Создано: %s. Удалить: /dn%d</b>\n\n%s\n\n", i+1, note.Created.Format(createdFieldFormat), note.ViewID, note.Text)
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

func (v *NoteView) fillHeader(idx int, note model.Note) string {
	messageWithEditTimetag := "<b>%d. Создано: %s. Изменено: %s. Удалить: /dn%d. Изменить: /editn%d</b>"
	defaultMessage := "<b>%d. Создано: %s. Удалить: /dn%d. Изменить: /editn%d</b>"

	// если заполнено поле последнее изменение - заполняем
	if note.LastEditSql.Valid {
		return fmt.Sprintf(messageWithEditTimetag, idx, note.Created.Format(createdFieldFormat), note.LastEditSql.Time.Format(createdFieldFormat),
			note.ViewID, note.ViewID)
	}

	return fmt.Sprintf(defaultMessage, idx, note.Created.Format(createdFieldFormat),
		note.ViewID, note.ViewID)
}

// Next возвращает следующую страницу сообщений
func (v *NoteView) Next() string {
	logrus.Debugf("noteView: getting next page. Current: %d\n", v.currentPage)

	if v.currentPage == v.total()-1 {
		logrus.Debugf("noteView: current page is the last. Setting current page to 0.\n")
		v.currentPage = 0
	} else {
		v.currentPage++
		logrus.Debugf("noteView: incrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Previous возвращает предыдущую страницу сообщений
func (v *NoteView) Previous() string {
	logrus.Debugf("noteView: getting previous page. Current: %d\n", v.currentPage)

	if v.currentPage == 0 {
		logrus.Debugf("noteView: previous page is the last. Setting current page to maximum: %d.\n", v.total())
		v.currentPage = v.total() - 1
	} else {
		v.currentPage--
		logrus.Debugf("noteView: decrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Last возвращает последнюю страницу сообщений
func (v *NoteView) Last() string {
	logrus.Debugf("noteView: getting the last page. Current: %d\n", v.currentPage)

	v.currentPage = v.total() - 1

	return v.pages[v.currentPage]
}

// First возвращает первую страницу сообщений
func (v *NoteView) First() string {
	logrus.Debugf("noteView: getting the first page. Current: %d\n", v.currentPage)

	v.currentPage = 0

	return v.pages[v.currentPage]
}

// current возвращает номер текущей страницы
func (v *NoteView) current() int {
	return v.currentPage + 1
}

// total возвращает общее количество страниц
func (v *NoteView) total() int {
	return len(v.pages)
}

// Keyboard делает клавиатуру для навигации по страницам
func (v *NoteView) Keyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	// если страниц 1, клавиатура не нужна
	if v.total() == 1 {
		menu.Inline(
			menu.Row(BtnSearchNotesByText, BtnSearchNotesByDate),
			menu.Row(BtnRefreshNotes),
			menu.Row(BtnDeleteAllNotes),
			menu.Row(BtnBackToMenu),
		)
		return menu
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := menu.Data(text, "")

	menu.Inline(
		menu.Row(BtnFirstPgNotes, BtnPrevPgNotes, btn, BtnNextPgNotes, BtnLastPgNotes),
		menu.Row(BtnSearchNotesByText, BtnSearchNotesByDate),
		menu.Row(BtnRefreshNotes),
		menu.Row(BtnDeleteAllNotes),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// Keyboard делает клавиатуру для навигации по страницам во время поиска
func (v *NoteView) KeyboardForSearch() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	// если страниц 1, клавиатура не нужна
	if v.total() == 1 {
		menu.Inline(
			menu.Row(BtnNotes),
			menu.Row(BtnBackToMenu),
		)
		return menu
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := menu.Data(text, "")

	menu.Inline(
		menu.Row(BtnFirstPgNotes, BtnPrevPgNotes, btn, BtnNextPgNotes, BtnLastPgNotes),
		menu.Row(BtnNotes),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

// SetCurrentToFirst устанавливает текущий номер страницы на 1
func (v *NoteView) SetCurrentToFirst() {
	v.currentPage = 0
}

// Clear используется когда удаляются все заметки: очищает список заметок, устанавливает текущую страницу в 0
func (v *NoteView) Clear() {
	v.currentPage = 0
	v.pages = make([]string, 0)
}

// Calendar возвращает календарь с текущим месяцем и годом
func (v *NoteView) Calendar() *tele.ReplyMarkup {
	calendar := v.calendar.currentCalendar()

	calendarWithBtns := v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToDateType)

	return calendarWithBtns
}

// PrevMonth возвращает календарь с предыдущим месяцем
func (v *NoteView) PrevMonth() *tele.ReplyMarkup {
	calendar := v.calendar.prevMonth()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToDateType)
	return calendar
}

// NextMonth возвращает календарь со следующим месяцем
func (v *NoteView) NextMonth() *tele.ReplyMarkup {
	calendar := v.calendar.nextMonth()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToDateType)
	return calendar
}

// PrevYear возвращает календарь с предыдущим годом
func (v *NoteView) PrevYear() *tele.ReplyMarkup {
	calendar := v.calendar.prevYear()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToDateType)
	return calendar
}

// NextYear возвращает календарь с следующим годом
func (v *NoteView) NextYear() *tele.ReplyMarkup {
	calendar := v.calendar.nextYear()
	calendar = v.calendar.addButns(calendar, BtnBackToMenu, BtnBackToDateType)
	return calendar
}

// GetDaysBtns возвращает слайс кнопок с числами месяца
func (v *NoteView) GetDaysBtns() []tele.Btn {
	return v.calendar.getDaysBtns()
}

// SetCurMonth устаналивает месяц в текущий
func (v *NoteView) SetCurMonth() {
	v.calendar.setCurMonth()
}

// SetCurYear устаналивает год в текущий
func (v *NoteView) SetCurYear() {
	v.calendar.setCurYear()
}

// SetCurMonth возвращает текущий месяц
func (v *NoteView) CurMonth() time.Month {
	return v.calendar.month()
}

// SetCurYear возвращает текущий год
func (v *NoteView) CurYear() int {
	return v.calendar.year()
}
