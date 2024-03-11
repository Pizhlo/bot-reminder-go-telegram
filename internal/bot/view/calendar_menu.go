package view

import (
	"fmt"
	"time"

	tele "gopkg.in/telebot.v3"
)

var (
	// inline кнопка для переключения на предыдущий месяц
	BtnPrevMonth = tele.Btn{Text: "<", Unique: "prev_month"}
	// inline кнопка для переключения на следующий месяц
	BtnNextMonth = tele.Btn{Text: ">", Unique: "next_month"}

	// inline кнопка для переключения на предыдущий год
	BtnPrevYear = tele.Btn{Text: "<<", Unique: "prev_year"}
	// inline кнопка для переключения на следующий год
	BtnNextYear = tele.Btn{Text: ">>", Unique: "next_year"}
)

// Calendar предоставляет календарь пользователю для выбора даты
type Calendar struct {
	curMonth time.Month
	curYear  int
	daysBtns []tele.Btn
}

type day struct {
	value   int // число месяца
	weekDay int // день недели
}

func new() *Calendar {
	return &Calendar{}
}

func (c *Calendar) getDaysBtns() []tele.Btn {
	return c.daysBtns
}

func (c *Calendar) month() time.Month {
	return c.curMonth
}

func (c *Calendar) year() int {
	return c.curYear
}

// currentCalendar предоставляет клавиатуру с календарем на текущий месяц и год
func (c *Calendar) currentCalendar() *tele.ReplyMarkup {
	return c.keyboard()
}

// prevMonth предоставляет клавиатуру с календарем на предыдущий месяц
func (c *Calendar) prevMonth() *tele.ReplyMarkup {
	if c.curMonth == 1 {
		c.curMonth = 12
		c.curYear -= 1
	} else {
		c.curMonth -= 1
	}

	return c.keyboard()
}

// nextMonth предоставляет клавиатуру с календарем на следующий месяц
func (c *Calendar) nextMonth() *tele.ReplyMarkup {
	if c.curMonth == 12 {
		c.curMonth = 1
		c.curYear += 1
	} else {
		c.curMonth += 1
	}

	return c.keyboard()
}

// prevYear предоставляет клавиатуру с календарем на предыдущий год
func (c *Calendar) prevYear() *tele.ReplyMarkup {
	c.curYear -= 1

	return c.keyboard()
}

// nextYear предоставляет клавиатуру с календарем на следующий год
func (c *Calendar) nextYear() *tele.ReplyMarkup {
	c.curYear += 1

	return c.keyboard()
}

// keyboard делает клавиатуру с календарем на указанный месяц и год
func (c *Calendar) keyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	btns := make([]tele.Btn, 0)

	// меню с названием месяца и года и кнопками для переключения между датами
	topMenu := c.topMenu()

	topMenuRows := menu.Split(6, topMenu)

	//btns = append(btns, topMenu...)

	// кнопки с днями месяца
	daysBtns := c.daysButtons()

	// кнопки с названиями дней недели - Пн, Вт, Ср, ...
	wdays := c.weekDaysButtons()

	// совмещаем вместе в одну клавиатуру
	btns = append(btns, wdays...)
	btns = append(btns, daysBtns...)

	rows := topMenuRows

	rows = append(rows, menu.Split(7, btns)...)

	//rows = append(rows, menu.Row(BtnBackToMenu, BtnBackToReminderType))

	menu.Inline(
		rows...,
	)

	return menu
}

// daysButtons генерирует кнопки с днями месяца
func (c *Calendar) daysButtons() []tele.Btn {
	res := []tele.Btn{}

	days := c.generateDays()

	// день недели первого дня в месяце
	firstWeekday := days[0].weekDay

	daysBefore := countDaysBefore(firstWeekday)

	// заполняем пробелы до первого дня. Например, если первый день - среда, будет 2 пробела
	for i := 0; i < daysBefore; i++ {
		res = append(res, tele.Btn{
			Text:   " ",
			Unique: " ",
		})
	}

	// заполняем днями с числами месяца
	for _, day := range days {
		btn := tele.Btn{
			Text:   fmt.Sprintf("%d", day.value),
			Unique: fmt.Sprintf("%d", day.value),
		}

		res = append(res, btn)

		c.daysBtns = append(c.daysBtns, btn)
	}

	// день недели последнего дня в месяце
	lastWeekDay := days[len(days)-1].weekDay

	daysAfter := countDaysAfter(lastWeekDay)

	// заполняем пробелы после последнего дня. Например, если последний день - пятница, будет 2 пробела
	for i := 0; i < daysAfter; i++ {
		res = append(res, tele.Btn{
			Text:   " ",
			Unique: " ",
		})
	}

	return res
}

// countDaysBefore возвращает количество дней, предшествующих первому дню месяца
func countDaysBefore(weekDay int) int {
	// день недели воскресенье - нужно добавить 6 дней перед ним
	if weekDay == 0 {
		return 6
	}

	return weekDay - 1
}

// countDaysBefore возвращает количество дней, идущих после последнего дня месяца
func countDaysAfter(weekDay int) int {
	// день недели воскресенье - добавлять дни не нужно
	if weekDay == 0 {
		return 0
	}

	return 7 - weekDay
}

// setCurYear устанавливает год в текущий
func (c *Calendar) setCurYear() {
	c.curYear = time.Now().Year()
}

// setCurMonth устанавливает месяц в текущий
func (c *Calendar) setCurMonth() {
	c.curMonth = time.Now().Month()
}

// daysInMonthCount возвращает количество дней в месяце, установленном в поле curMonth
func (c *Calendar) daysInMonthCount() int {
	return time.Date(c.curYear, c.curMonth+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// generateDays составляет дни на весь месяц
func (c *Calendar) generateDays() []day {
	days := []day{}
	daysInMonth := c.daysInMonthCount()

	for i := 1; i <= daysInMonth; i++ {
		weeekDay := int(time.Date(c.curYear, c.curMonth, i, 0, 0, 0, 0, time.Local).Weekday())
		day := day{value: i, weekDay: weeekDay}
		days = append(days, day)
	}

	return days
}

// weekDaysButtons возвращает кнопки с названиями дней недели
func (c *Calendar) weekDaysButtons() []tele.Btn {
	weekDays := map[int]string{
		1: "Пн",
		2: "Вт",
		3: "Ср",
		4: "Чт",
		5: "Пт",
		6: "Сб",
		7: "Вс",
	}

	res := []tele.Btn{}

	for i := 1; i < 8; i++ {
		btn := tele.Btn{
			Text:   weekDays[i],
			Unique: " ",
		}

		res = append(res, btn)
	}

	return res
}

// topMenu возвращает меню, расположенное над календарем - там находятся название месяца, года,
// и кнопки для переключения между датами
func (c *Calendar) topMenu() []tele.Btn {
	menu := []tele.Btn{}

	monthBtn, yearBtn := c.monthTitle(), c.yearTitle()

	menu = append(menu, BtnPrevYear, BtnPrevMonth, monthBtn, yearBtn, BtnNextMonth, BtnNextYear)

	return menu
}

// monthTitle возвращает кнопку с названием текущего месяца
func (c *Calendar) monthTitle() tele.Btn {
	monthsTranslation := map[string]string{
		"January":   "Янв",
		"February":  "Фев",
		"March":     "Март",
		"April":     "Апр",
		"May":       "Май",
		"June":      "Июнь",
		"July":      "Июль",
		"August":    "Авг",
		"September": "Сент",
		"October":   "Окт",
		"November":  "Ноя",
		"December":  "Дек",
	}

	monthTitle := monthsTranslation[c.curMonth.String()]

	return tele.Btn{
		Text:   monthTitle,
		Unique: monthTitle,
	}
}

// yearTitle возвращает кнопку с надписью текущего года
func (c *Calendar) yearTitle() tele.Btn {
	return tele.Btn{
		Text:   fmt.Sprintf("%d", c.curYear),
		Unique: fmt.Sprintf("%d", c.curYear),
	}
}

func (c *Calendar) addButns(dst *tele.ReplyMarkup, btns ...tele.Btn) *tele.ReplyMarkup {
	inlineBtns := []tele.InlineButton{}

	for _, btn := range btns {
		inlineBtns = append(inlineBtns, *btn.Inline())
	}

	dst.InlineKeyboard = append(dst.InlineKeyboard, inlineBtns)

	return dst
}
