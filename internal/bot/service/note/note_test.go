package note

import (
	"testing"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/stretchr/testify/assert"
)

func TestSaveUser_Positive(t *testing.T) {
	type test struct {
		name    string
		userID  int64
		viewMap map[int64]*view.NoteView
	}

	tests := []test{
		{
			name:    "user is not saved in the map",
			userID:  1,
			viewMap: make(map[int64]*view.NoteView),
		},
		{
			name:   "user is saved in the map",
			userID: 1,
			viewMap: map[int64]*view.NoteView{
				1: view.NewNote(),
			},
		},
	}

	for _, tt := range tests {
		srv := New(nil)

		srv.viewsMap = tt.viewMap

		srv.SaveUser(tt.userID)

		_, ok := srv.viewsMap[tt.userID]
		assert.Equal(t, true, ok)
	}
}
