package note

import (
	"context"
	"fmt"

	tele "gopkg.in/telebot.v3"
)

// GetAll возвращает все заметки пользователя
func (s *NoteService) GetAll(ctx context.Context, userID int64) (string, *tele.ReplyMarkup, error) {
	s.logger.Debugf("Getting all user's notes. User ID: %d\n", userID)

	notes, err := s.noteEditor.GetAllByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("error while getting all notes by user ID %d: %v\n", userID, err)
		return "", nil, fmt.Errorf("error while getting all notes by user ID %d: %w", userID, err)
	}

	return s.viewsMap[userID].Message(notes), s.viewsMap[userID].Keyboard(), nil
}

func (s *NoteService) NextPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Next(), s.viewsMap[userID].Keyboard()
}

func (s *NoteService) PrevPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Previous(), s.viewsMap[userID].Keyboard()
}

func (s *NoteService) LastPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Last(), s.viewsMap[userID].Keyboard()
}

func (s *NoteService) FirstPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].First(), s.viewsMap[userID].Keyboard()
}
