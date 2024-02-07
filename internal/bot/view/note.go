package view

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

const (
	dateFormat       = "02.01.2006 15:04:05"
	noteCountPerPage = 5
	maxMessageLen    = 4096
)

type NoteView struct {
	pages       []string
	currentPage int
	logger      *logrus.Logger
}

func NewNote() *NoteView {
	return &NoteView{pages: make([]string, 0), currentPage: 0, logger: logger.New()}
}

var (
	selector = &tele.ReplyMarkup{}

	// inline кнопка для переключения на предыдущую страницу (заметки)
	BtnPrevPgNotes = selector.Data("<", "prev")
	// inline кнопка для переключения на следующую страницу (заметки)
	BtnNextPgNotes = selector.Data(">", "next")

	// inline кнопка для переключения на первую страницу (заметки)
	BtnFirstPgNotes = selector.Data("<<", "start")
	// inline кнопка для переключения на последнюю страницу (заметки)
	BtnLastPgNotes = selector.Data(">>", "end")
)

// Message формирует список сообщений из моделей заметок и возвращает первую страницу.
// Количество заметок на одной странице задает переменная noteCountPerPage (по умолчанию - 5)
func (v *NoteView) Message(notes []model.Note) string {
	if len(notes) == 0 {
		return messages.NotesNotFoundMessage
	}

	var res = ""
	v.pages = make([]string, 0)

	for i, note := range notes {
		res += fmt.Sprintf("%d. Создано: %s. Удалить: /del%d\n\n%s\n\n", i+1, note.Created.Format(dateFormat), note.ID, note.Text)
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
func (v *NoteView) Next() string {
	v.logger.Debugf("noteView: getting next page. Current: %d\n", v.currentPage)

	if v.currentPage == v.total()-1 {
		v.logger.Debugf("noteView: current page is the last. Setting current page to 0.\n")
		v.currentPage = 0
	} else {
		v.currentPage++
		v.logger.Debugf("noteView: incrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Previous возвращает предыдущую страницу сообщений
func (v *NoteView) Previous() string {
	v.logger.Debugf("noteView: getting previous page. Current: %d\n", v.currentPage)

	if v.currentPage == 0 {
		v.logger.Debugf("noteView: previous page is the last. Setting current page to maximum: %d.\n", v.total())
		v.currentPage = v.total() - 1
	} else {
		v.currentPage--
		v.logger.Debugf("noteView: decrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Last возвращает последнюю страницу сообщений
func (v *NoteView) Last() string {
	v.logger.Debugf("noteView: getting the last page. Current: %d\n", v.currentPage)

	v.currentPage = v.total() - 1

	return v.pages[v.currentPage]
}

// First возвращает первую страницу сообщений
func (v *NoteView) First() string {
	v.logger.Debugf("noteView: getting the first page. Current: %d\n", v.currentPage)

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
	// если страниц 1, клавиатура не нужна
	if v.total() == 1 {
		return &tele.ReplyMarkup{}
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := selector.Data(text, "")

	selector.Inline(
		selector.Row(BtnFirstPgNotes, BtnPrevPgNotes, btn, BtnNextPgNotes, BtnLastPgNotes),
		selector.Row(BtnBackToMenu),
	)

	return selector
}

// SetCurrentToFirst устанавливает текущий номер страницы на 1
func (v *NoteView) SetCurrentToFirst() {
	v.currentPage = 0
}
