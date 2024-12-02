package sharedaccess

import "gopkg.in/telebot.v3"

// NotesBySpace возвращает заметки, принадлежащие конкретному пространству, которое уже было выбрано, поэтому
// для запроса нужен только userID
func (s *SharedSpace) NotesBySpace(userID int64) (string, *telebot.ReplyMarkup, error) {
	msg, err := s.viewsMap[userID].Notes()
	if err != nil {
		return "", nil, err
	}

	kb := s.viewsMap[userID].KeyboardForNotes()

	return msg, kb, nil
}

// NextPage обрабатывает кнопку переключения на следующую страницу заметок совместного пространства
func (s *SharedSpace) NextPageNotes(userID int64) (string, *telebot.ReplyMarkup) {
	return s.viewsMap[userID].Next(), s.viewsMap[userID].KeyboardForNotes()
}

// PrevPage обрабатывает кнопку переключения на предыдущую страницу заметок совместного пространства
func (s *SharedSpace) PrevPageNotes(userID int64) (string, *telebot.ReplyMarkup) {
	return s.viewsMap[userID].Previous(), s.viewsMap[userID].KeyboardForNotes()
}

// LastPage обрабатывает кнопку переключения на последнюю страницу заметок совместного пространства
func (s *SharedSpace) LastPageNotes(userID int64) (string, *telebot.ReplyMarkup) {
	return s.viewsMap[userID].Last(), s.viewsMap[userID].KeyboardForNotes()
}

// FirstPage обрабатывает кнопку переключения на первую страницу заметок совместного пространства
func (s *SharedSpace) FirstPageNotes(userID int64) (string, *telebot.ReplyMarkup) {
	return s.viewsMap[userID].First(), s.viewsMap[userID].KeyboardForNotes()
}
