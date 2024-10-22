package note

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/elastic"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Save сохраняет заметку в базе данных. Для сохранения требуется: ID пользователя, содержимое заметки, дата создания
func (db *NoteRepo) Save(ctx context.Context, note model.Note) error {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("error while creating transaction: %w", err)
	}

	var id uuid.UUID
	row := tx.QueryRowContext(ctx, `insert into notes.notes (user_id, text, created) values((select id from users.users where tg_id=$1), $2, $3) returning ID`, note.TgID, note.Text, note.Created)

	err = row.Scan(&id)
	if err != nil {
		return fmt.Errorf("error scanning ID: %+v", err)
	}

	// создаем структуру для сохранения в elastic
	elasticData := elastic.Data{
		Index: elastic.NoteIndex,
		Model: elastic.Note{
			ID:   id,
			Text: note.Text,
			TgID: note.TgID,
		}}

	// сохраняем в elastic
	err = db.elasticClient.Save(ctx, elasticData)
	if err != nil {
		// отменяем транзакцию в случае ошибки (для консистентности данных)
		_ = tx.Rollback()
		return fmt.Errorf("error saving note to Elastic: %+v", err)
	}

	logrus.Debugf("NoteRepo: successfully saved user's note. User: %+v", note.TgID)

	return tx.Commit()
}
