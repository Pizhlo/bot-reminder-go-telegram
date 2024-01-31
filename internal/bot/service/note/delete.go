package note

import "context"

// DeleteAll удаляет все заметки пользователя по user ID
func (n *NoteService) DeleteAll(ctx context.Context, userID int64) error {
	n.logger.Debugf("Note service: deleting all user's notes... Setting current page to 1.\n")
	// устанавливаем во view номер страницы на первый
	n.viewsMap[userID].SetCurrentToFirst()

	n.logger.Debugf("Note service: deleting all user's notes from DB... \n")
	// удаляем все заметки
	return n.noteEditor.DeleteAllByUserID(ctx, userID)
}
