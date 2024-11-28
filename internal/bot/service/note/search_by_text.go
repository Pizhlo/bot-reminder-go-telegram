package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

// SearchByText ищет заметки по тексту. Возвращает сообщение с заметками, клавиатуру, ошибку
func (s *NoteService) SearchByText(ctx context.Context, note model.SearchByText) (string, *tele.ReplyMarkup, error) {
	notes, err := s.noteEditor.SearchByText(ctx, note)
	if err != nil {
		return "", nil, err
	}

	msg, err := s.viewsMap[note.TgID].Message(notes)

	return msg, s.viewsMap[note.TgID].KeyboardForSearch(), err
}
