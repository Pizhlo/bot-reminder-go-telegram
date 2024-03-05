package note

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

// SearchByText ищет заметки, созданные в выбранную дату. Возвращает сообщение с заметками, клавиатуру, ошибку
func (s *NoteService) SearchOneDate(ctx context.Context, note model.SearchByOneDate) (string, *tele.ReplyMarkup, error) {
	s.logger.Debugf("Note service: looking for note by one date. Model: %+v\n", note)

	notes, err := s.noteEditor.SearchByOneDate(ctx, note)
	if err != nil {
		return "", nil, err
	}

	return s.viewsMap[note.TgID].Message(notes), s.viewsMap[note.TgID].Keyboard(), nil
}
