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
	ID        uuid.UUID // id из базы
	ElasticID string    // id в elastic Search
	TgID      int64
	Text      string
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

// SearchByTextQuery возвращает готовый запрос для поиска по тексту
func (d *Data) SearchByTextQuery() (*search.Request, error) {
	switch d.Index {
	case NoteIndex:
		return d.searchNoteByTextQuery()
	case ReminderIndex:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown elastic index: %s", d.Index)
	}
}

// SearchByIDQuery возвращает готовый запрос для поиска по ID.
// Ищет в эластике по ID из базы
func (d *Data) SearchByIDQuery() (*search.Request, error) {
	switch d.Index {
	case NoteIndex:
		return d.searchNoteByID()
	case ReminderIndex:
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown elastic index: %s", d.Index)
	}
}

func (d *Data) searchNoteByTextQuery() (*search.Request, error) {
	note, err := d.ValidateNote()
	if err != nil {
		return nil, err
	}

	must1 := []types.Query{
		{
			Bool: &types.BoolQuery{
				Should: []types.Query{
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
		{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						Match: map[string]types.MatchQuery{
							"TgID": {
								Query: fmt.Sprintf("%d", note.TgID),
							},
						},
					},
				},
			},
		},
	}

	req := &search.Request{
		Query: &types.Query{
			Bool: &types.BoolQuery{
				Must: must1,
			},
		},
	}

	return req, nil
}

func (d *Data) searchNoteByID() (*search.Request, error) {
	note, err := d.ValidateNote()
	if err != nil {
		return nil, err
	}

	req := &search.Request{
		Query: &types.Query{
			Match: map[string]types.MatchQuery{
				"ID": {
					Query: note.ID.String(),
				},
			},
		},
	}

	return req, nil
}

func valueToPointer[T string | float32 | float64](val T) *T {
	return &val
}
