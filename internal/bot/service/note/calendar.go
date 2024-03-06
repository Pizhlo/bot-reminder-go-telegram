package note

import (
	"time"

	tele "gopkg.in/telebot.v3"
)

// Calendar возвращает календарь с текущим месяцем и годом
func (v *NoteService) Calendar(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].Calendar()
}

// PrevMonth перелистывает календарь на месяц назад
func (v *NoteService) PrevMonth(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].PrevMonth()
}

// NextMonth перелистывает календарь на месяц назад
func (v *NoteService) NextMonth(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].NextMonth()
}

// PrevYear перелистывает календарь на месяц назад
func (v *NoteService) PrevYear(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].PrevYear()
}

// NextYear перелистывает календарь на месяц назад
func (v *NoteService) NextYear(userID int64) *tele.ReplyMarkup {
	return v.viewsMap[userID].NextYear()
}

func (v *NoteService) CurMonth(userID int64) time.Month {
	return v.viewsMap[userID].CurMonth()
}

func (v *NoteService) CurYear(userID int64) int {
	return v.viewsMap[userID].CurYear()
}

// DaysBtns возвращает слайс с кнопками, которые содержат в себе число месяца
func (v *NoteService) DaysBtns(userID int64) []tele.Btn {
	return v.viewsMap[userID].GetDaysBtns()
}
