package elastic

import (
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/operator"
	"github.com/google/uuid"
)

// Структура для хранения и поиска по заметкам
type Note struct {
	ID   uuid.UUID // id из базы
	TgID int64
	Text string
}

// ValidateNote проверяет поля структуры elastic.Data на правильность и возвращает заметку
func (d *Data) ValidateNote() (*Note, error) {
	if d.Index != NoteIndex {
		return nil, fmt.Errorf("index name not equal `notes`: `%+v`", d.Index)
	}

	val, ok := d.Model.(Note)
	if !ok {
		return nil, fmt.Errorf("type of field note.Model is not elastic.Note: %+v", d.Model)
	}

	return &val, nil
}

// GetNote возвращает модель заметки, если она валидна, либо ошибку
func (d *Data) GetNote() (*Note, error) {
	val, ok := d.Model.(Note)
	if !ok {
		return nil, fmt.Errorf("type of field note.Model is not elastic.Note: %+v", d.Model)
	}

	return &val, nil
}

// SearchNoteQuery возвращает готовый запрос для поиска по заметкам
func (d *Data) SearchNoteQuery() (*search.Request, error) {
	note, err := d.ValidateNote()
	if err != nil {
		return nil, err
	}

	req := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						Match: map[string]types.MatchQuery{
							"TgID": {
								Query: fmt.Sprintf("%d", note.TgID),
							},
						},
					},
					{
						Match: map[string]types.MatchQuery{
							"Text": {
								Query:     note.Text,
								Operator:  &operator.Or,
								Fuzziness: "auto",
							},
						},
					},
					{
						Wildcard: map[string]types.WildcardQuery{
							"Text": {
								Value:   valueToPointer(fmt.Sprintf("*%s*", note.Text)),
								Boost:   valueToPointer(float32(1.0)),
								Rewrite: valueToPointer("constant_score"),
							},
						},
					},
				},
			},
		},
	}

	return req, nil
}

func valueToPointer[T string | float32 | float64](val T) *T {
	return &val
}
