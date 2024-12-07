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
	BtnSpaceParticipants = tele.Btn{Text: "🫂Участники", Unique: "shared_space_participants"}
	// inline кнопка для добавления напоминания в совметное пространство
	BtnAddReminder = tele.Btn{Text: "Добавить напоминание", Unique: "add_reminder_to_shared_space"}

	// inline кнопка для возврата в совместное пространство
	BtnBackToSharedSpace = tele.Btn{Text: "⬅️Назад", Unique: "back_to_shared_space"}
	// inline кнопка для возврата в совместные пространства
	BtnBackToAllSharedSpaces = tele.Btn{Text: "⬅️Назад", Unique: "shared_space"}
	// inline кнопка для возврата к списку участников
	BtnBackToParticipants = tele.Btn{Text: "⬅️Назад", Unique: "shared_space_participants"}

	// inline кнопка для добавления участников
	BtnAddParticipants = tele.Btn{Text: "➕Добавить", Unique: "add_users_to_shared_space"}
	// inline кнопка для исключения участников
	BtnRemoveParticipants = tele.Btn{Text: "🚫Исключить", Unique: "remove_user_from_shared_space"}

	// invitations
	BtnAcceptInvitations = tele.Btn{Text: "✅Принять", Unique: "accept_invintation"}
	BtnDenyInvitations   = tele.Btn{Text: "❌Отклонить", Unique: "deny_invintation"}

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
	pages             []string
	currentPage       int
	spacesMap         map[int]model.SharedSpace
	btns              []tele.Btn
	currentSpaceIndex int

	noteView     *NoteView
	reminderView *ReminderView
}

func NewSharedSpaceView() *SharedSpaceView {
	return &SharedSpaceView{
		pages:             make([]string, 0),
		currentPage:       0,
		currentSpaceIndex: 0,
		spacesMap:         make(map[int]model.SharedSpace, 0),
		btns:              make([]tele.Btn, 0),
		noteView:          NewNote(),
		reminderView:      NewReminder()}
}

func (s *SharedSpaceView) Message(spaces []model.SharedSpace) string {
	var res = "Твои совместные пространства:\n\n"

	s.pages = make([]string, 0)

	i := 0
	for _, space := range spaces {
		// сохраняем пространства, они понадобятся для того чтобы сделать клавиатуру
		s.spacesMap[space.ID] = space

		i++

		res += fmt.Sprintf(messages.SharedSpaceMessage,
			i, space.Name,
			fmt.Sprintf("Участников: %d",
				len(space.Participants)),
			len(space.Notes),
			len(space.Reminders),
			space.Created.Format(createdFieldFormat))
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

	s.currentSpaceIndex = spaceID

	logrus.Debugf("SharedSpaceView: MessageBySpace set currentSpaceID to %d", spaceID)

	return s.messageBySpace(space), nil
}

func (s *SharedSpaceView) MessageByCurrentSpace() (string, error) {
	space, ok := s.spacesMap[s.currentSpaceIndex]
	if !ok {
		return "", fmt.Errorf("not found space by ID %d", s.currentSpaceIndex)
	}

	logrus.Debugf("SharedSpaceView: MessageByCurrentSpace currentSpaceID %d", s.currentSpaceIndex)

	return s.messageBySpace(space), nil
}

func (s *SharedSpaceView) messageBySpace(space model.SharedSpace) string {
	participants := formatParticipants(space.Participants, space.Creator.TGID)

	return fmt.Sprintf(messages.SharedSpaceMessage, space.ViewID, space.Name, participants, len(space.Notes), len(space.Reminders), space.Created.Format(createdFieldFormat))
}

func formatParticipants(participants []model.Participant, creatorID int64) string {
	participantsTxt := "Участники:\n"

	for _, u := range participants {
		if u.TGID == creatorID {
			participantsTxt += fmt.Sprintf("* @%s - админ\n", u.UsernameSQL.String)
		} else {
			if u.State == model.PendingState {
				participantsTxt += fmt.Sprintf("* @%s - ожидание ответа\n", u.Username)
				continue
			}

			if u.State == model.RejectedState {
				continue
			}

			participantsTxt += fmt.Sprintf("* @%s\n", u.UsernameSQL.String)
		}
	}

	return participantsTxt
}

// CurrentSpaceName возвращает название текущего (выбранного) совметного доступа
func (s *SharedSpaceView) CurrentSpaceName() string {
	logrus.Debugf("SharedSpaceView: CurrentSpaceName currentSpaceID %d", s.currentSpaceIndex)
	return s.spacesMap[s.currentSpaceIndex].Name
}

// CurrentSpaceName возвращает ID текущего (выбранного) совметного доступа
func (s *SharedSpaceView) CurrentSpaceID() int {
	logrus.Debugf("SharedSpaceView: CurrentSpaceID currentSpaceID %d", s.currentSpaceIndex)
	return s.spacesMap[s.currentSpaceIndex].ID
}

// CurrentSpace возвращает текущее выбранное совместное пространство
func (s *SharedSpaceView) CurrentSpace() model.SharedSpace {
	logrus.Debugf("SharedSpaceView: CurrentSpace currentSpaceID %d", s.currentSpaceIndex)
	return s.spacesMap[s.currentSpaceIndex]
}

func (s *SharedSpaceView) Notes() (string, error) {
	space := s.spacesMap[s.currentSpaceIndex]

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
	space := s.spacesMap[s.currentSpaceIndex]

	if len(space.Reminders) == 0 {
		return fmt.Sprintf(messages.NoRemindersInSharedSpaceMessage, space.Name), nil
	}

	return s.reminderView.Message(space.Reminders)
}

// ParticipantsMessage возвращает сообщение для пункта меню "Участники"
func (s *SharedSpaceView) ParticipantsMessage() string {
	space := s.spacesMap[s.currentSpaceIndex]
	participants := space.Participants

	msg := fmt.Sprintf("Участники пространства <b>%s</b>:\n\n", space.Name)

	txt := formatParticipants(participants, space.Creator.TGID)

	logrus.Debugf("SharedSpaceView: ParticipantsMessage currentSpaceID %d", s.currentSpaceIndex)

	return fmt.Sprintf("%s%s", msg, txt)
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
