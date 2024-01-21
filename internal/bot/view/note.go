package view

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

const (
	dateFormat       = "02.01.2006 15:04:05 MST Mon"
	noteCountPerPage = 5
)

var (
	selector = &tele.ReplyMarkup{}

	BtnPrevPg = selector.Data("<", "prev")
	BtnNextPg = selector.Data(">", "next")

	BtnFirstPg = selector.Data("<<", "start")
	BtnLastPg  = selector.Data(">>", "end")
)

type View struct {
	messages    []string
	currentPage int
	logger      *logrus.Logger
}

func New() *View {
	return &View{messages: make([]string, 0), currentPage: 0, logger: logger.New()}
}

// Message формирует список сообщений из моделей заметок и возвращает первую страницу.
// Количество заметок на одной странице задает переменная noteCountPerPage (по умолчанию - 5)
func (v *View) Message(notes []model.Note) string {
	if len(notes) == 0 {
		return "У тебя пока нет заметок. Чтобы создать, просто пришли мне текст/фото, и я сохраню это!"
	}

	var res = ""
	v.messages = make([]string, 0)

	for i, note := range notes {
		res += fmt.Sprintf("%d. Создано: %s. Удалить: /del%d\n\n%s\n\n", i+1, note.Created.Format(dateFormat), i+1, note.Text)
		if i%noteCountPerPage == 0 && i > 0 {
			v.messages = append(v.messages, res)
			res = ""
		}
	}

	return v.messages[0]
}

// Next возвращает следующую страницу сообщений
func (v *View) Next() string {
	v.logger.Debugf("View: getting next page. Current: %d\n", v.currentPage)

	if v.currentPage == len(v.messages)-1 {
		v.logger.Debugf("View: current page is the last. Setting current page to 0.\n")
		v.currentPage = 0
	} else {
		v.currentPage++
		v.logger.Debugf("View: incrementing current page. New value: %d\n", v.currentPage)
	}

	return v.messages[v.currentPage]
}

// Previous возвращает предыдущую страницу сообщений
func (v *View) Previous() string {
	v.logger.Debugf("View: getting previous page. Current: %d\n", v.currentPage)

	if v.currentPage == 0 {
		v.logger.Debugf("View: previous page is the last. Setting current page to maximum: %d.\n", len(v.messages))
		v.currentPage = len(v.messages) - 1
	} else {
		v.currentPage--
		v.logger.Debugf("View: decrementing current page. New value: %d\n", v.currentPage)
	}

	return v.messages[v.currentPage]
}

// Last возвращает последнюю страницу сообщений
func (v *View) Last() string {
	v.logger.Debugf("View: getting the last page. Current: %d\n", v.currentPage)

	v.currentPage = len(v.messages) - 1

	return v.messages[v.currentPage]
}

// First возвращает первую страницу сообщений
func (v *View) First() string {
	v.logger.Debugf("View: getting the first page. Current: %d\n", v.currentPage)

	v.currentPage = 0

	return v.messages[v.currentPage]
}

// current возвращает номер текущей страницы
func (v *View) current() int {
	return v.currentPage + 1
}

// total возвращает общее количество страниц
func (v *View) total() int {
	return len(v.messages)
}

// Keyboard делает клавиатуру для навигации по страницам
func (v *View) Keyboard() *tele.ReplyMarkup {
	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := selector.Data(text, "")

	selector.Inline(
		selector.Row(BtnFirstPg, BtnPrevPg, btn, BtnNextPg, BtnLastPg),
	)

	return selector
}
