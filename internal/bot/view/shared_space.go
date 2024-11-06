package view

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

type SharedSpaceView struct {
	pages       []string
	currentPage int
	spaces      []model.SharedSpace
}

func NewSharedSpaceView() *SharedSpaceView {
	return &SharedSpaceView{pages: make([]string, 0), currentPage: 0, spaces: make([]model.SharedSpace, 0)}
}

func (s *SharedSpaceView) Message(spaces []model.SharedSpace) string {
	var res = "Твои совместные пространства:\n\n"

	s.pages = make([]string, 0)

	// сохраняем пространства, они понадобятся для того чтобы сделать клавиатуру
	s.spaces = spaces

	messageTxt := "<b>%d. %s</b>\n\nУчастники: %+v\n\nСоздано: %+v\n\n"

	for _, space := range spaces {
		res += fmt.Sprintf(messageTxt, space.ViewID, space.Name, space.Participants, space.Created.Format(createdFieldFormat))
	}

	if len(s.pages) < 5 && res != "" {
		s.pages = append(s.pages, res)
	}

	s.currentPage = 0

	return s.pages[0]
}

func (s *SharedSpaceView) Keyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	btns := []tele.Btn{}

	for _, space := range s.spaces {
		btn := tele.Btn{Text: space.Name, Unique: space.Name}

		btns = append(btns, btn)
	}

	if len(btns) > 0 {
		menu.Inline(
			menu.Row(btns...),
			menu.Row(BtnCreateSharedSpace),
			menu.Row(BtnBackToMenu),
		)
	} else {
		menu.Inline(
			menu.Row(BtnCreateSharedSpace),
			menu.Row(BtnBackToMenu),
		)
	}

	return menu
}
