package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

// SearchByText ищет заметки по тексту. Возвращает сообщение с заметками, клавиатуру, ошибку
func (s *NoteService) SearchByText(ctx context.Context, note model.SearchNoteByText) (string, *tele.ReplyMarkup, error) {
	s.logger.Debugf("Note service: saving user's note. Model: %+v\n", note)

	notes, err := s.noteEditor.SearchByText(ctx, note)
	if err != nil {
		return "", nil, err
	}

	return s.viewsMap[note.TgID].Message(notes), s.viewsMap[note.TgID].NoteKeyboard(), nil
}
