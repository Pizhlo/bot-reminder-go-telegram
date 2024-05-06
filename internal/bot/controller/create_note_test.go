package controller

import (
	"context"
	"testing"
	"time"

	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/mocks"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	model_user "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model/user"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/note"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/user"
	tz_cache "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/storage/cache/timezone"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/Pizhlo/bot-reminder-go-telegram/pkg/random"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/telebot.v3"
)

func TestCreateNote(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	telectx := mocks.NewMockteleCtx(ctrl)
	noteEditor := mocks.NewMocknoteEditor(ctrl)
	tzEditor := mocks.NewMocktimezoneEditor(ctrl)
	tz := tz_cache.New()

	// при создании user service
	tzEditor.EXPECT().GetAll(gomock.Any()).Return([]*model_user.User{}, nil)

	chat := &telebot.Chat{ID: int64(1), FirstName: random.String(5)}

	expectedNote := random.Note()
	expectedNote.TgID = chat.ID

	loc := time.FixedZone("Europe/Moscow", 1)
	tz.Save(context.Background(), chat.ID, loc)

	telectx.EXPECT().Chat().Return(chat).Times(4)
	telectx.EXPECT().Text().Return(expectedNote.Text)

	// s.noteEditor.Save(ctx, note)
	noteEditor.EXPECT().Save(gomock.Any(), gomock.Any()).Do(func(ctx interface{}, note model.Note) {
		assert.Equal(t, expectedNote.Text, note.Text)
		assert.Equal(t, expectedNote.TgID, note.TgID)
		assert.Equal(t, expectedNote.Created.Local(), expectedNote.Created.Local())
	})

	noteSrv := note.New(noteEditor)
	userSrv := user.New(context.Background(), nil, tz, tzEditor)

	expectedText := messages.SuccessfullyCreatedNoteMessage
	expectedKb := view.NotesAndMenuBtns()

	telectx.EXPECT().EditOrSend(gomock.Any(), gomock.Any()).Do(func(text string, kb *telebot.ReplyMarkup) {
		assert.Equal(t, expectedText, text)
		assert.Equal(t, expectedKb, kb)
	}).Return(nil)

	controller := New(userSrv, noteSrv, nil, nil, 0)

	err := controller.CreateNote(context.Background(), telectx)
	assert.NoError(t, err)
}
