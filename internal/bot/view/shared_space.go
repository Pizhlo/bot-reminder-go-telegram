package view

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

type SharedSpaceView struct {
	pages       []string
	currentPage int
}

func NewSharedSpaceView() *SharedSpaceView {
	return &SharedSpaceView{pages: make([]string, 0), currentPage: 0}
}

func (s *SharedSpaceView) Message(spaces []model.SharedSpace) string {
	var res = "Твои совместные пространства:\n\n"

	s.pages = make([]string, 0)

	messageTxt := "%d. %s\n\nУчастники: %+v\n\nСоздано: %+v"

	for i, space := range spaces {
		res += fmt.Sprintf(messageTxt, i, space.Name, space.Participants, space.Created)
	}

	if len(s.pages) < 5 && res != "" {
		s.pages = append(s.pages, res)
	}

	s.currentPage = 0

	return s.pages[0]
}
