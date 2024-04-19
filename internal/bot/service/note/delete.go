package note

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

// DeleteAll удаляет все заметки пользователя по user ID
func (n *NoteService) DeleteAll(ctx context.Context, userID int64) error {
	logrus.Debug(wrap("deleting all user's notes... Setting current page to 1.\n"))

	// устанавливаем во view номер страницы на первый
	n.viewsMap[userID].Clear()

	logrus.Debug(wrap("deleting all user's notes from DB... \n"))
	// удаляем все заметки
	return n.noteEditor.DeleteAllByUserID(ctx, userID)
}

// DeleteByID удаляет заметку пользователя по user ID
func (n *NoteService) DeleteByID(ctx context.Context, userID int64, noteID int) error {
	logrus.Debug(wrap(fmt.Sprintf("deleting user's note by ID: %d. Checking if user has note with this ID...\n", noteID)))

	// проверяем, существует ли заметка с таким номером
	_, err := n.noteEditor.GetByViewID(ctx, userID, noteID)
	if err != nil {
		logrus.Debug(wrap(fmt.Sprintf("error while checking note ID %d: %v\n", noteID, err)))
		return err
	}

	logrus.Debug(wrap(fmt.Sprintf("found note by ID %d. Deleting...\n", noteID)))
	// удаляем все заметки
	return n.noteEditor.DeleteNoteByViewID(ctx, userID, noteID)
}
