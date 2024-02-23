package reminder

import tele "gopkg.in/telebot.v3"

// Calendar возвращает календарь с текущим месяцем и годом
func (v *ReminderService) Calendar(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].Calendar()
}

// PrevMonth перелистывает календарь на месяц назад
func (v *ReminderService) PrevMonth(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].PrevMonth()
}

// NextMonth перелистывает календарь на месяц назад
func (v *ReminderService) NextMonth(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].NextMonth()
}

// PrevYear перелистывает календарь на месяц назад
func (v *ReminderService) PrevYear(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].PrevYear()
}

// NextYear перелистывает календарь на месяц назад
func (v *ReminderService) NextYear(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].NextYear()
}
