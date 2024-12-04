package view

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

func (s *SharedSpaceView) Keyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	btns := []tele.Btn{}

	for _, space := range s.spacesMap {
		btn := tele.Btn{Text: space.Name, Unique: fmt.Sprintf("%v", space.ID)}

		btns = append(btns, btn)
	}

	if len(btns) > 0 {
		menu.Inline(
			menu.Row(btns...),
			menu.Row(BtnCreateSharedSpace),
			menu.Row(BtnBackToMenu),
		)

		s.btns = btns
	} else {
		menu.Inline(
			menu.Row(BtnCreateSharedSpace),
			menu.Row(BtnBackToMenu),
		)
	}

	return menu
}

// KeyboardForSpace возвращает клавиатуру для управления совместным пространством
func KeyboardForSpace() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(
			BtnNotesSharedSpace, BtnRemindersSharedSpace,
		),
		menu.Row(BtnSpaceParticipants),
		menu.Row(BtnBackToAllSharedSpaces),
	)

	return menu
}

func (s *SharedSpaceView) Buttons() []tele.Btn {
	return s.btns
}

func KeyboardForReminders() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnAddReminder),
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}

func BackToSharedSpaceMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}

func InvintationKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnAcceptInvintation, BtnDenyInvintation),
	)

	return menu
}

// Keyboard делает клавиатуру для навигации по страницам заметок
func (v *SharedSpaceView) KeyboardForNotes() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	// если страниц 1, клавиатура не нужна
	if v.total() == 1 {
		menu.Inline(
			menu.Row(BtnRefreshNotesSharedSpace),
			menu.Row(BtnBackToSharedSpace),
		)
		return menu
	}

	text := fmt.Sprintf("%d / %d", v.current(), v.total())

	btn := menu.Data(text, "")

	menu.Inline(
		menu.Row(BtnFirstPgNotesSharedSpace, BtnPrevPgNotesSharedSpace, btn, BtnNextPgNotesSharedSpace, BtnLastPgNotesSharedSpace),
		menu.Row(BtnRefreshNotesSharedSpace),
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}

func ParticipantsKeyboard() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{}

	menu.Inline(
		menu.Row(BtnAddParticipants),
		menu.Row(BtnRemoveParticipants),
		menu.Row(BtnBackToSharedSpace),
	)

	return menu
}