package view

import tele "gopkg.in/telebot.v3"

type NavigationView struct{}

func NewNavigation() *NavigationView {
	return &NavigationView{}
}

var (
	mainMenu = &tele.ReplyMarkup{}

	// inline кнопка для просмотра профиля
	BtnProfile = selector.Data("👤Профиль", "profile")
	// inline кнопка для просмотра настроек
	BtnSettings = selector.Data("⚙️Настройки", "settings")

	// inline кнопка просмотра заметок
	BtnNotes = selector.Data("📝Заметки", "notes")
	// inline кнопка просмотра напоминаний
	BtnReminders = selector.Data("⏰Напоминания", "reminders")

	// inline кнопка для возвращения в меню
	BtnMenu = selector.Data("⬅️Меню", "menu")
)

// MainMenu возвращает главное меню.
// Кнопки: Профиль, Настройки, Заметки, Напоминания
func (v *NavigationView) MainMenu() *tele.ReplyMarkup {
	mainMenu.Inline(
		mainMenu.Row(BtnProfile, BtnSettings),
		mainMenu.Row(BtnNotes, BtnReminders),
	)

	return mainMenu
}
