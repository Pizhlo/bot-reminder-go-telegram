package note

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// GetAll обрабатывает запрос всех заметок пользователя. Возвращает: первую страницу заметок (в виде целого сообщения),
// клавиатуру и ошибку
func (s *NoteService) GetAll(ctx context.Context, userID int64) (string, *tele.ReplyMarkup, error) {
	notes, err := s.noteEditor.GetAllByUserID(ctx, userID)
	if err != nil {
		logrus.Error(wrap(fmt.Sprintf("error while getting all notes by user ID %d: %v\n", userID, err)))
		return "", nil, err
	}

	msg, err := s.viewsMap[userID].Message(notes)

	return msg, s.viewsMap[userID].Keyboard(), err
}

// NextPage обрабатывает кнопку переключения на следующую страницу
func (s *NoteService) NextPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Next(), s.viewsMap[userID].Keyboard()
}

// PrevPage обрабатывает кнопку переключения на предыдущую страницу
func (s *NoteService) PrevPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Previous(), s.viewsMap[userID].Keyboard()
}

// LastPage обрабатывает кнопку переключения на последнюю страницу
func (s *NoteService) LastPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Last(), s.viewsMap[userID].Keyboard()
}

// FirstPage обрабатывает кнопку переключения на первую страницу
func (s *NoteService) FirstPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].First(), s.viewsMap[userID].Keyboard()
}
