package messageeditor

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
)

type messageEditor struct {
	notes            []model.Note
	maxMessageLenght int
}

func New(notes []model.Note) *messageEditor {
	return &messageEditor{notes: notes, maxMessageLenght: 3000}
}

var notesMessage = "%s - создано %s\n\nДля удаления нажмите /del%d\n\n"

func (m *messageEditor) MakeNotesMessage() []string {
	sumLen := 0
	var message string
	var messages []string

	for _, note := range m.notes {
		// str := fmt.Sprintf(notesMessage, note.Text, note.Created.Format("02.01.2006 в 15:04"), note.ID)
		// nextLen := sumLen + len([]rune(str))

		// if nextLen >= m.maxMessageLenght {
		// 	messages = append(messages, message)
		// 	sumLen = 0
		// 	message = ""
		// }

		message += fmt.Sprintf(notesMessage, note.Text, note.Created.Format("02.01.2006 в 15:04"), note.ID)
		sumLen += len([]rune(message))

		if sumLen >= m.maxMessageLenght { //  telegram: message is too long (400) - length must not be more than 4096
			messages = append(messages, message)
			sumLen = 0
			message = ""
		}
	}

	if len(messages) == 0 {
		messages = append(messages, message)
	}

	return messages
}
