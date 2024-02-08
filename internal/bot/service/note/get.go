package note

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

// GetAll обрабатывает запрос всех заметки пользователя. Возвращает: первую страницу заметок (в виде целого сообщения),
// клавиатуру и ошибку
func (s *NoteService) GetAll(ctx context.Context, userID int64) (string, *tele.ReplyMarkup, error) {
	s.logger.Debugf("Note service: getting all user's notes. User ID: %d\n", userID)

	notes, err := s.noteEditor.GetAllByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Note service: error while getting all notes by user ID %d: %v\n", userID, err)
		return "", nil, err
	}

	s.logger.Debugf("Note service: got %d user's notes\n", len(notes))

	return s.viewsMap[userID].Message(notes), s.viewsMap[userID].NoteKeyboard(), nil
}

// NextPage обрабатывает кнопку переключения на следующую страницу
func (s *NoteService) NextPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Next(), s.viewsMap[userID].NoteKeyboard()
}

// PrevPage обрабатывает кнопку переключения на предыдущую страницу
func (s *NoteService) PrevPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Previous(), s.viewsMap[userID].NoteKeyboard()
}

// LastPage обрабатывает кнопку переключения на последнюю страницу
func (s *NoteService) LastPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Last(), s.viewsMap[userID].NoteKeyboard()
}

// FirstPage обрабатывает кнопку переключения на первую страницу
func (s *NoteService) FirstPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].First(), s.viewsMap[userID].NoteKeyboard()
}
