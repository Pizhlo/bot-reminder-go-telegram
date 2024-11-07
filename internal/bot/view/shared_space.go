package view

import (
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

var (
	// inline кнопка просмотра заметок в совместном пространстве
	BtnNotesSharedSpace = tele.Btn{Text: "📝Заметки", Unique: "notes_by_shared_space"}
	// inline кнопка просмотра напоминаний в совместном пространстве
	BtnRemindersSharedSpace = tele.Btn{Text: "⏰Напоминания", Unique: "reminders_by_shared_space"}

	// inline кнопка для управления участниками в совместном пространстве
	BtnSpaceParticipants = tele.Btn{Text: "Участники", Unique: "shared_space_participants"}
	// inline кнопка для добавления заметки в совметное пространство
	BtnAddNote = tele.Btn{Text: "📝Добавить заметку", Unique: "add_note_to_shared_space"}
	// inline кнопка для добавления напоминания в совметное пространство
	BtnAddReminder = tele.Btn{Text: "Добавить напоминание", Unique: "add_reminder_to_shared_space"}

	// inline кнопка для возврата в совместное пространство
	BtnBackToSharedSpace = tele.Btn{Text: "⬅️Назад", Unique: "back_to_shared_space"}
	// inline кнопка для возврата в совместные пространства
	BtnBackToAllSharedSpaces = tele.Btn{Text: "⬅️Назад", Unique: "back_to_all_shared_spaces"}

	// inline кнопка для добавления участников
	BtnAddParticipants = tele.Btn{Text: "Добавить", Unique: "add_users_to_shared_space"}
	// inline кнопка для исключения участников
	BtnRemoveParticipants = tele.Btn{Text: "Исключить", Unique: "add_users_to_shared_space"}
)

type SharedSpaceView struct {
	pages        []string
	currentPage  int
	spacesMap    map[int]model.SharedSpace
	btns         []tele.Btn
	currentSpace int

	noteView     *NoteView
	reminderView *ReminderView
}

func NewSharedSpaceView() *SharedSpaceView {
	return &SharedSpaceView{
		pages:        make([]string, 0),
		currentPage:  0,
		currentSpace: 0,
		spacesMap:    make(map[int]model.SharedSpace, 0),
		btns:         make([]tele.Btn, 0),
		noteView:     NewNote(),
		reminderView: NewReminder()}
}

func (s *SharedSpaceView) Message(spaces []model.SharedSpace) string {
	var res = "Твои совместные пространства:\n\n"

	s.pages = make([]string, 0)

	i := 0
	for _, space := range spaces {
		// сохраняем пространства, они понадобятся для того чтобы сделать клавиатуру
		s.spacesMap[space.ID] = space

		i++
		participantsTxt := ""

		for _, u := range space.Participants {
			if u.TGID == space.Creator.TGID {
				participantsTxt += fmt.Sprintf("* @%s - админ\n", u.UsernameSQL.String)
			} else {
				participantsTxt += fmt.Sprintf("* @%s\n", u.UsernameSQL.String)
			}
		}

		res += fmt.Sprintf(messages.SharedSpaceMessage, i, space.Name, participantsTxt, len(space.Notes), len(space.Reminders), space.Created.Format(createdFieldFormat))
	}

	if len(s.pages) < 5 && res != "" {
		s.pages = append(s.pages, res)
	}

	s.currentPage = 0

	return s.pages[0]
}

func (s *SharedSpaceView) MessageBySpace(spaceID int) (string, error) {
	space, ok := s.spacesMap[spaceID]
	if !ok {
		return "", fmt.Errorf("not found space by ID %d", spaceID)
	}

	s.currentSpace = spaceID

	return s.messageBySpace(space), nil
}

func (s *SharedSpaceView) MessageByCurrentSpace() (string, error) {
	space, ok := s.spacesMap[s.currentSpace]
	if !ok {
		return "", fmt.Errorf("not found space by ID %d", s.currentSpace)
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

	for _, space := range s.spacesMap {
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
		menu.Row(BtnSpaceParticipants),
		menu.Row(BtnBackToAllSharedSpaces),
	)

	return menu
}

func (s *SharedSpaceView) Buttons() []tele.Btn {
	return s.btns
}

func (s *SharedSpaceView) Notes() string {
	space := s.spacesMap[s.currentSpace]

	if len(space.Notes) == 0 {
		return fmt.Sprintf(messages.NoNotesInSharedSpaceMessage, space.Name)
	}

	return s.noteView.Message(space.Notes)
}

func (s *SharedSpaceView) KeyboardForNotes() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnAddNote),
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}

func (s *SharedSpaceView) Reminders() (string, error) {
	space := s.spacesMap[s.currentSpace]

	if len(space.Reminders) == 0 {
		return fmt.Sprintf(messages.NoRemindersInSharedSpaceMessage, space.Name), nil
	}

	return s.reminderView.Message(space.Reminders)
}

func (s *SharedSpaceView) KeyboardForReminders() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnAddReminder),
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}

// ParticipantsMessage возвращает сообщение для пункта меню "Участники"
func (s *SharedSpaceView) ParticipantsMessage() string {
	space := s.spacesMap[s.currentSpace]
	participants := space.Participants

	msg := fmt.Sprintf("Участники пространства <b>%s</b>:\n\n", space.Name)

	for _, u := range participants {
		msg += fmt.Sprintf("* @%s\n", u.UsernameSQL.String)
	}

	return msg
}

func (s *SharedSpaceView) ParticipantsKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnAddParticipants),
		menu.Row(BtnRemoveParticipants),
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}
