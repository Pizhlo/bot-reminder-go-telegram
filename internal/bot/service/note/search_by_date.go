package note

import (
	"context"
	"errors"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

// SearchByText ищет заметки, созданные в выбранную дату. Возвращает сообщение с заметками, клавиатуру, ошибку
func (s *NoteService) SearchByOneDate(ctx context.Context, note model.SearchByOneDate) (string, *tele.ReplyMarkup, error) {
	notes, err := s.noteEditor.SearchByOneDate(ctx, note)
	if err != nil {
		return "", nil, err
	}

	return s.viewsMap[note.TgID].Message(notes), s.viewsMap[note.TgID].Keyboard(), nil
}

// SaveFirstDate сохраняет первую дату для поиска заметок по двум датам
func (s *NoteService) SaveFirstDate(userID int64, date time.Time) {
	s.searchMap[userID] = model.SearchByTwoDates{
		TgID:      userID,
		FirstDate: date,
	}
}

// SaveSecondDate сохраняет вторую дату для поиска заметок по двум датам
func (s *NoteService) SaveSecondDate(userID int64, date time.Time) error {
	if val, ok := s.searchMap[userID]; ok {
		val.SecondDate = date
		s.searchMap[userID] = val
		return nil
	}

	return errors.New("no data found for this user")
}

// GetSearchNote возвращает информацию для поиска заметок по двум датам для определенного пользователя
func (s *NoteService) GetSearchNote(userID int64) (*model.SearchByTwoDates, error) {
	if val, ok := s.searchMap[userID]; ok {
		return &val, nil
	}

	return nil, errors.New("no data found for this user")
}

// ValidateSearchDate проверяет, не раньше ли вторая дата первой
func (s *NoteService) ValidateSearchDate(userID int64, date time.Time) error {
	val := s.searchMap[userID]

	if date.Before(val.FirstDate) {
		return api_errors.ErrSecondDateBeforeFirst
	}

	// вторая дата не должна быть в будущем
	today := time.Now()

	if date.After(today) {
		return api_errors.ErrSecondDateFuture
	}

	return nil
}

// SearchByTwoDates производит поиск по двум датам
func (s *NoteService) SearchByTwoDates(ctx context.Context, userID int64) (string, *tele.ReplyMarkup, error) {
	note, err := s.GetSearchNote(userID)
	if err != nil {
		return "", nil, err
	}

	notes, err := s.noteEditor.SearchByTwoDates(ctx, note)
	if err != nil {
		return "", nil, err
	}

	return s.viewsMap[note.TgID].Message(notes), s.viewsMap[note.TgID].Keyboard(), nil
}
