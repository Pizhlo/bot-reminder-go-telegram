package reminder

import (
	"context"

	tele "gopkg.in/telebot.v3"
)

// GetAll обрабатывает запрос всех напоминаний пользователя. Возвращает: первую страницу напоминаний (в виде целого сообщения),
// клавиатуру и ошибку
func (s *ReminderService) GetAll(ctx context.Context, userID int64) (string, *tele.ReplyMarkup, error) {
	s.logger.Debugf("Reminder service: getting all user's reminders. User ID: %d\n", userID)

	reminders, err := s.reminderEditor.GetAllByUserID(ctx, userID)
	if err != nil {
		s.logger.Errorf("Reminder service: error while getting all reminders by user ID %d: %v\n", userID, err)
		return "", nil, err
	}

	s.logger.Debugf("Reminder service: got %d user's reminders\n", len(reminders))

	msg, err := s.viewsMap[userID].Message(reminders)
	if err != nil {
		s.logger.Errorf("Reminder service: error while making first message for all reminders by user ID %d: %v\n", userID, err)
		return "", nil, err
	}

	return msg, s.viewsMap[userID].Keyboard(), nil
}

// NextPage обрабатывает кнопку переключения на следующую страницу
func (s *ReminderService) NextPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Next(), s.viewsMap[userID].Keyboard()
}

// PrevPage обрабатывает кнопку переключения на предыдущую страницу
func (s *ReminderService) PrevPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Previous(), s.viewsMap[userID].Keyboard()
}

// LastPage обрабатывает кнопку переключения на последнюю страницу
func (s *ReminderService) LastPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].Last(), s.viewsMap[userID].Keyboard()
}

// FirstPage обрабатывает кнопку переключения на первую страницу
func (s *ReminderService) FirstPage(userID int64) (string, *tele.ReplyMarkup) {
	return s.viewsMap[userID].First(), s.viewsMap[userID].Keyboard()
}
