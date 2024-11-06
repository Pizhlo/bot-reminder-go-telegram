package view

import (
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

var (
	// inline кнопка просмотра заметок
	BtnNotesSharedSpace = tele.Btn{Text: "📝Заметки", Unique: "notes_by_shared_space"}
	// inline кнопка просмотра напоминаний
	BtnRemindersSharedSpace = tele.Btn{Text: "⏰Напоминания", Unique: "reminders_by_shared_space"}
	// inline кнопка для добавления пользователей в совместное пространство
	BtnAddUsersToSpace = tele.Btn{Text: "Добавить пользователей", Unique: "add_users_to_shared_space"}
)

type SharedSpaceView struct {
	pages       []string
	currentPage int
	spaces      map[int]model.SharedSpace
	btns        []tele.Btn
}

func NewSharedSpaceView() *SharedSpaceView {
	return &SharedSpaceView{
		pages:       make([]string, 0),
		currentPage: 0,
		spaces:      make(map[int]model.SharedSpace, 0),
		btns:        make([]tele.Btn, 0)}
}

func (s *SharedSpaceView) Message(spaces map[int]model.SharedSpace) string {
	var res = "Твои совместные пространства:\n\n"

	s.pages = make([]string, 0)

	// сохраняем пространства, они понадобятся для того чтобы сделать клавиатуру
	s.spaces = spaces

	for _, space := range spaces {
		participantsTxt := ""

		for _, u := range space.Participants {
			participantsTxt += fmt.Sprintf("* @%s\n", u.UsernameSQL.String)
		}

		res += fmt.Sprintf(messages.SharedSpaceMessage, space.ViewID, space.Name, participantsTxt, len(space.Notes), len(space.Reminders), space.Created.Format(createdFieldFormat))
	}

	if len(s.pages) < 5 && res != "" {
		s.pages = append(s.pages, res)
	}

	s.currentPage = 0

	return s.pages[0]
}

func (s *SharedSpaceView) MessageBySpace(spaceID int) (string, error) {
	space, ok := s.spaces[spaceID]
	if !ok {
		return "", fmt.Errorf("not found space by ID %d", spaceID)
	}

	return s.messageBySpace(space), nil
}

func (s *SharedSpaceView) messageBySpace(space model.SharedSpace) string {
	participantsTxt := ""

	for _, u := range space.Participants {
		participantsTxt += fmt.Sprintf("* @%s\n", u.UsernameSQL.String)
	}

	return fmt.Sprintf(messages.SharedSpaceMessage, space.ViewID, space.Name, participantsTxt, len(space.Notes), len(space.Reminders), space.Created.Format(createdFieldFormat))
}

func (s *SharedSpaceView) Keyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	btns := []tele.Btn{}

	for _, space := range s.spaces {
		btn := tele.Btn{Text: space.Name, Unique: fmt.Sprintf("%v", space.ID)}

		btns = append(btns, btn)
	}

	if len(btns) > 0 {
		menu.Inline(
			menu.Row(btns...),
			menu.Row(BtnCreateSharedSpace),
			menu.Row(BtnBackToMenu),
		)

		s.btns = btns
	} else {
		menu.Inline(
			menu.Row(BtnCreateSharedSpace),
			menu.Row(BtnBackToMenu),
		)
	}

	return menu
}

func (s *SharedSpaceView) KeyboardForSpace() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(
			BtnNotesSharedSpace, BtnRemindersSharedSpace,
		),
		menu.Row(BtnAddUsersToSpace),
		menu.Row(BtnBackToMenu),
	)

	return menu
}

func (s *SharedSpaceView) Buttons() []tele.Btn {
	return s.btns
}
