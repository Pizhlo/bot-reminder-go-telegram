package view

import (
	"fmt"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

var (
	// inline кнопка просмотра заметок в совместном пространстве
	BtnNotesSharedSpace = tele.Btn{Text: "📝Заметки", Unique: "notes_by_shared_space"}
	// inline кнопка просмотра напоминаний в совместном пространстве
	BtnRemindersSharedSpace = tele.Btn{Text: "⏰Напоминания", Unique: "reminders_by_shared_space"}

	// inline кнопка для управления участниками в совместном пространстве
	BtnSpaceParticipants = tele.Btn{Text: "Участники", Unique: "shared_space_participants"}
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

	// invintation
	BtnAcceptInvintation = tele.Btn{Text: "✅Принять", Unique: "accept_invintation"}
	BtnDenyInvintation   = tele.Btn{Text: "❌Отклонить", Unique: "deny_invintation"}

	// notes buttons

	// inline кнопка для переключения на предыдущую страницу (заметки)
	BtnPrevPgNotesSharedSpace = tele.Btn{Text: "<", Unique: "prev_pg_notes_shared_space"}
	// inline кнопка для переключения на следующую страницу (заметки)
	BtnNextPgNotesSharedSpace = tele.Btn{Text: ">", Unique: "next_pg_notes_shared_space"}

	// inline кнопка для обновления заметок
	BtnRefreshNotesSharedSpace = tele.Btn{Text: "🔁", Unique: "notes_shared_space"}

	// inline кнопка для переключения на первую страницу (заметки)
	BtnFirstPgNotesSharedSpace = tele.Btn{Text: "<<", Unique: "start_pg_notes_shared_space"}
	// inline кнопка для переключения на последнюю страницу (заметки)
	BtnLastPgNotesSharedSpace = tele.Btn{Text: ">>", Unique: "end_pg_notes_shared_space"}
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

		participants := participantsTxt(space.Participants, space.Creator.TGID)

		res += fmt.Sprintf(messages.SharedSpaceMessage, i, space.Name, participants, len(space.Notes), len(space.Reminders), space.Created.Format(createdFieldFormat))
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
	participants := participantsTxt(space.Participants, space.Creator.TGID)

	return fmt.Sprintf(messages.SharedSpaceMessage, space.ViewID, space.Name, participants, len(space.Notes), len(space.Reminders), space.Created.Format(createdFieldFormat))
}

func participantsTxt(participants []model.User, creatorID int64) string {
	participantsTxt := ""

	for _, u := range participants {
		if u.TGID == creatorID {
			participantsTxt += fmt.Sprintf("* @%s - админ\n", u.UsernameSQL.String)
		} else {
			participantsTxt += fmt.Sprintf("* @%s\n", u.UsernameSQL.String)
		}
	}

	return participantsTxt
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

// CurrentSpaceName возвращает название текущего (выбранного) совметного доступа
func (s *SharedSpaceView) CurrentSpaceName() string {
	return s.spacesMap[s.currentSpace].Name
}

// CurrentSpaceName возвращает ID текущего (выбранного) совметного доступа
func (s *SharedSpaceView) CurrentSpaceID() int {
	return s.spacesMap[s.currentSpace].ID
}

// CurrentSpace возвращает текущее выбранное совместное пространство
func (s *SharedSpaceView) CurrentSpace() model.SharedSpace {
	return s.spacesMap[s.currentSpace]
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

func (s *SharedSpaceView) Notes() (string, error) {
	space := s.spacesMap[s.currentSpace]

	if len(space.Notes) == 0 {
		return fmt.Sprintf(messages.NoNotesInSharedSpaceMessage, space.Name), nil
	}

	// pages := textForRecord(space.Notes, "")

	// return s.SharedSpaceView.Message(space.Notes)

	res := ""

	s.pages = make([]string, 0)

	for i, note := range space.Notes {
		header := s.noteView.fillHeader(i+1, note)

		msg := fmt.Sprintf("<b>%s. Автор: @%v</b>\n\n%s\n\n", header, note.Creator.Username, note.Text)

		res += msg

		if i%noteCountPerPage == 0 && i > 0 || len(res) == maxMessageLen {
			s.pages = append(s.pages, res)
			res = ""
		}
	}

	if len(s.pages) < 5 && res != "" {
		s.pages = append(s.pages, res)
	}

	s.currentPage = 0

	return s.pages[0], nil
}

// func (s *SharedSpaceView) KeyboardForNotes() *tele.ReplyMarkup {
// 	menu := &tele.ReplyMarkup{}

// 	menu.Inline(
// 		menu.Row(BtnAddNote),
// 		menu.Row(BtnBackToSharedSpace),
// 	)

// 	return menu
// }

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

func (s *SharedSpaceView) BackToSharedSpaceMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}

func (s *SharedSpaceView) InvintationKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnAcceptInvintation, BtnDenyInvintation),
	)

	return menu
}

// Next возвращает следующую страницу сообщений
func (v *SharedSpaceView) Next() string {
	logrus.Debugf("SharedSpaceView: getting next page. Current: %d\n", v.currentPage)

	if v.currentPage == v.total()-1 {
		logrus.Debugf("SharedSpaceView: current page is the last. Setting current page to 0.\n")
		v.currentPage = 0
	} else {
		v.currentPage++
		logrus.Debugf("SharedSpaceView: incrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Previous возвращает предыдущую страницу сообщений
func (v *SharedSpaceView) Previous() string {
	logrus.Debugf("SharedSpaceView: getting previous page. Current: %d\n", v.currentPage)

	if v.currentPage == 0 {
		logrus.Debugf("SharedSpaceView: previous page is the last. Setting current page to maximum: %d.\n", v.total())
		v.currentPage = v.total() - 1
	} else {
		v.currentPage--
		logrus.Debugf("SharedSpaceView: decrementing current page. New value: %d\n", v.currentPage)
	}

	return v.pages[v.currentPage]
}

// Last возвращает последнюю страницу сообщений
func (v *SharedSpaceView) Last() string {
	logrus.Debugf("SharedSpaceView: getting the last page. Current: %d\n", v.currentPage)

	v.currentPage = v.total() - 1

	return v.pages[v.currentPage]
}

// First возвращает первую страницу сообщений
func (v *SharedSpaceView) First() string {
	logrus.Debugf("SharedSpaceView: getting the first page. Current: %d\n", v.currentPage)

	v.currentPage = 0

	return v.pages[v.currentPage]
}

// current возвращает номер текущей страницы
func (v *SharedSpaceView) current() int {
	return v.currentPage + 1
}

// total возвращает общее количество страниц
func (v *SharedSpaceView) total() int {
	return len(v.pages)
}

// Keyboard делает клавиатуру для навигации по страницам заметок
func (v *SharedSpaceView) KeyboardForNotes() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	// если страниц 1, клавиатура не нужна
	if v.total() == 1 {
		menu.Inline(
			menu.Row(BtnRefreshNotesSharedSpace),
			menu.Row(BtnBackToSharedSpace),
		)
		return menu
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := menu.Data(text, "")

	menu.Inline(
		menu.Row(BtnFirstPgNotesSharedSpace, BtnPrevPgNotesSharedSpace, btn, BtnNextPgNotesSharedSpace, BtnLastPgNotesSharedSpace),
		menu.Row(BtnRefreshNotesSharedSpace),
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}
