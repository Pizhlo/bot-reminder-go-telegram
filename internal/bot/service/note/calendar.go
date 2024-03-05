package note

import tele "gopkg.in/telebot.v3"

// Calendar возвращает календарь с текущим месяцем и годом
func (v *NoteService) Calendar(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].Calendar()
}

// DaysBtns возвращает слайс с кнопками, которые содержат в себе число месяца
func (v *NoteService) DaysBtns(userID int64) []tele.Btn {
	return v.viewsMap[userID].GetDaysBtns()
}
