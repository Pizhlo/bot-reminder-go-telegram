package reminder

import (
	"context"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	tele "gopkg.in/telebot.v3"
)

// GetAll обрабатывает запрос всех напоминаний пользователя. Возвращает: первую страницу напоминаний (в виде целого сообщения),
// клавиатуру и ошибку
func (s *ReminderService) GetAll(ctx context.Context, userID int64) ([]model.Reminder, error) {
	return s.reminderEditor.GetAllByUserID(ctx, userID)
}

func (s *ReminderService) Message(userID int64, reminders []model.Reminder) (string, error) {
	return s.viewsMap[userID].Message(reminders)
}

func (s *ReminderService) Keyboard(userID int64) *tele.ReplyMarkup {
	return s.viewsMap[userID].Keyboard()
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
