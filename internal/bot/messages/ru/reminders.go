package messages

const NoRemindersMessage = "У тебя пока нет напоминаний. Чтобы создать, нажми кнопку"

// CREATE
const ReminderNameMessage = "Напиши название напоминания"
const TypeOfReminderMessage = "Ага, записываю...\n<b>Название:</b> %s (чтобы поменять название, просто пришли новое).\n\nТеперь выбери тип напоминания:\n\n<b>* Сегодня</b> - сегодня в указанное время\n<b>* Завтра</b> - завтра в указанное время\n<b>* Несколько раз в день</b> - будет предложено на выбор раз в несколько минут и часов\n<b>* Ежедневно</b> - каждый день в указанное время\n<b>* Раз в неделю</b> - раз в неделю в указанный день недели и время\n<b>* Раз в несколько дней</b> - будет предложено ввести промежуток дней, в который будет приходить сообщение (например, раз в 10 дней)\n<b>* Раз в месяц</b> - один раз в указанное число месяца (от 1 до 31) в указанное время\n<b>* Раз в год</b> - раз в год в указанную дату (например, 13 июня)\n<b>* Выбрать дату</b> - уведомление сработает один раз в указанную дату (например, 13 июня 2025 года)"
const ReminderTimeMessage = "Напиши время, когда присылать уведомления, в формате ЧЧ:ММ"
const SuccessCreationMessage = "Напоминание <b>%s</> успешно создано!🥳\n\nОно %s: %s\n\nСледующее срабатывание: %s"

const ChooseMinutesOrHoursMessage = "Ага, записываю...\n<b>Название:</b> %s (чтобы поменять название, просто пришли новое).\n<b>Тип:</b> %s\n\nТеперь выбери, с какой частотой присылать уведомления: раз в какое-то количество минут/часов или в указанное время (можно будет указать список, например: в 10:30, 15:40, 19:00)"
const ChooseMinutesOrHoursWithoutNameMessage = "Теперь выбери, с какой частотой присылать уведомления"
const MinutesDurationMessage = "Напиши, с какой частотой в минутах тебе присылать уведомления (от 1 до 59)"
const HoursDurationMessage = "Напиши, с какой частотой в часах тебе присылать уведомления (от 1 до 24)"
const TimesReminderMessage = "Напиши список, когда тебе присылать уведомления (например: в 10:30, 16:30, 17:00)"
const ChooseWeekDayMessage = "Ага, записываю...\n<b>Название:</b> %s (чтобы поменять название, просто пришли новое).\n<b>Тип:</b> %s\n\nТеперь выбери день недели, когда присылать уведомления"
const MonthDayMessage = "Ага, записываю...\n<b>Название:</b> %s\n<b>Тип:</b> %s\n\nТеперь напиши число месяца от 1 до 31, когда присылать уведомления"
const DaysDurationMessage = "Ага, записываю...\n<b>Название:</b> %s\n<b>Тип:</b> %s\n\nТеперь напиши количество дней от 1 до 180, как часто присылать уведомления"
const CalendarMessage = "Ага, записываю...\n<b>Название:</b> %s (чтобы поменять название, просто пришли новое).\n<b>Тип:</b> %s\n\nТеперь выбери необходимую дату в календаре"

// ERRORS
const InvalidTimeMessage = "🙅‍♂️Неверный формат ввода. Напиши время в формате ЧЧ:ММ (например, 13:15)"
const TimeInPastMessage = "Это время уже прошло, напиши другое"
const UserDoesntHaveRemindersMessage = `У тебя пока нет напоминаний. Чтобы создать, перейди в раздел "Напоминания" и нажми кнопку "Создать напоминание"`
const InvalidMinutesMessage = "🙅‍♂️Ты ввел неверное значение. Количество минут должно быть от целым числом 1 до 59"
const InvalidHoursMessage = "🙅‍♂️Ты ввел неверное значение. Количество часов должно быть целым числом от 1 до 24"
const InvalidTimesMessage = "🙅‍♂️Неверный ввод. Напиши списком время в формате ЧЧ:ММ, когда присылать уведомления (например: в 10:30, 16:30, 17:00)"
const InvalidDaysMessage = "🙅‍♂️Ты ввел неверное значение. Количество дней должно быть целым числом от 1 до 180"
const InvalidDaysInMonthMessage = "🙅‍♂️Ты ввел неверное значение. Количество дней должно быть целым числом от 1 до 31"
const InvalidDateMessage = "🙅‍♂️Эта дата уже прошла, но ты можешь выбрать другую"

// USER REMINDER
const ReminderMessage = "%s\n\n<i>Напоминание сработало %s</i>"

// DELETE
const ConfirmDeleteRemindersMessage = "Ой-ой... Ты точно хочешь удалить ВСЕ напоминания?😥"
const AllRemindersDeletedMessage = "Все напоминания успешно удалены!👌"
const ReminderDeletedMessage = "Напоминание <b>%s</b> успешно удалено!🥳"
const ReminderDeletedByIDMessage = "Напоминание под номером <b>%d</b> успешно удалено!🥳"
