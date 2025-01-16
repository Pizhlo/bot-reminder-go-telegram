package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	api_errors "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/errors"
	messages "github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/messages/ru"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
)

// GetSharedAccess - хендлер для кнопки "Совметный доступ". Выгружает из базы совместные пространства пользователя, либо говорит, что их нет
func (c *Controller) GetSharedAccess(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.GetAllByUserID(ctx, telectx.Chat().ID)
	if err != nil {
		if errors.Is(err, api_errors.ErrSharedSpacesNotFound) {
			return telectx.Edit(messages.SharedSpacesNotFoundMessage, view.SharedAccessMenu())
		}

		logrus.Errorf("Error while handling shared spaces command. User ID: %d. Error: %+v\n", telectx.Chat().ID, err)

		return err
	}

	btns := c.sharedSpace.Buttons(telectx.Chat().ID)

	for _, btn := range btns {
		c.bot.Handle(&btn, func(telectx tele.Context) error {
			spaceID, err := strconv.Atoi(btn.Unique)
			if err != nil {
				return fmt.Errorf("error converting string space ID '%s' to int: %+v", btn.Unique, err)
			}

			return c.GetSharedSpace(ctx, telectx, spaceID)
		})
	}

	logrus.Debugf("Controller: successfully got all user's shared spaces. Sending message to user...\n")
	return telectx.Edit(msg, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

func (c *Controller) CreateSharedSpace(ctx context.Context, telectx tele.Context) error {
	loc, err := c.userSrv.GetLocation(ctx, telectx.Chat().ID)
	if err != nil {
		return err
	}

	space := model.SharedSpace{
		Name:    telectx.Message().Text,
		Created: time.Now().In(loc),
		Creator: model.Participant{
			User: model.User{
				TGID:     telectx.Chat().ID,
				Username: telectx.Chat().Username,
			},
			State: model.AddedState,
		},
	}

	if err := c.sharedSpace.Save(ctx, space); err != nil {
	}

	msg := fmt.Sprintf(messages.SharedSpaceCreationSuccessMessage, space.Name)
	return telectx.EditOrSend(msg, view.ShowSharedSpacesMenu())
}

func (c *Controller) GetSharedSpace(ctx context.Context, telectx tele.Context, spaceID int) error {
	msg, kb, err := c.sharedSpace.GetSharedSpace(spaceID, telectx.Chat().ID)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(msg, kb)
}

func (c *Controller) GetCurrentSharedSpace(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.CurrentSharedSpace(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(msg, kb)
}

func (c *Controller) NotesBySharedSpace(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.NotesBySpace(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(msg, kb)
}

func (c *Controller) RemindersBySharedSpace(ctx context.Context, telectx tele.Context) error {
	msg, kb, err := c.sharedSpace.RemindersBySpace(telectx.Chat().ID)
	if err != nil {
		return err
	}

	return telectx.EditOrSend(msg, kb)
}

func (c *Controller) SharedSpaceParticipants(ctx context.Context, telectx tele.Context) error {
	msg, kb := c.sharedSpace.SharedSpaceParticipants(telectx.Chat().ID)

	return telectx.EditOrSend(msg, kb)
}

// AddParticipant обрабатывает нажатие на кнопку "Добавить участника"
func (c *Controller) AddParticipant(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.AddParticipantMessage, view.BackToMenuBtn())
}

// RemoveParticipant обрабатывает нажатие на кнопку "Иключить участника"
func (c *Controller) RemoveParticipant(ctx context.Context, telectx tele.Context) error {
	kb := c.sharedSpace.ParticipantsListKeyboard(telectx.Chat().ID)
	// btns := c.sharedSpace.Buttons(telectx.Chat().ID)

	c.bot.Handle(&view.RemoveParticipantBtn, func(telectx tele.Context) error {
		return c.removeParticipant(ctx, telectx)
	})

	return telectx.EditOrSend(messages.RemoveParticipantMessage, kb)
}

func (c *Controller) removeParticipant(ctx context.Context, telectx tele.Context) error {
	// "user=%d space=%s"

	userAndSpace := strings.Split(telectx.Callback().Data, " ")

	userIDStr := strings.Split(userAndSpace[0], "=")
	spaceStr := strings.Split(userAndSpace[1], "=")

	spaceName := ""
	words := ""

	for i := 2; i < len(userAndSpace); i++ {
		words += fmt.Sprintf(" %s", userAndSpace[i])
	}

	spaceName = fmt.Sprintf("%s%s", spaceStr[1], words)

	userID, err := strconv.Atoi(userIDStr[1])
	if err != nil {
		return err
	}

	user, err := c.userSrv.GetByID(ctx, int64(userID))
	if err != nil {
		return err
	}

	space, err := c.sharedSpace.GetSharedSpaceByName(ctx, spaceName)
	if err != nil {
		return err
	}

	err = c.sharedSpace.DeleteParticipant(ctx, int64(space.ID), model.Participant{User: model.User{TGID: int64(userID)}})
	if err != nil {
		return err
	}

	msg := fmt.Sprintf(messages.UserWasRemovedMessage, space.Name, formatUsername(telectx.Chat()))

	err = c.sendMessage(int64(userID), msg, nil)
	if err != nil {
		return err
	}

	msg = fmt.Sprintf(messages.UserSuccesfullyRemoved, user.Username, space.Name)

	return telectx.EditOrSend(msg, view.BackToMenuBtn())
}

// HandleParticipant обрабатывает ссылку на пользователя, которого надо добавить в совместное пространство
func (c *Controller) HandleParticipant(ctx context.Context, telectx tele.Context, spaceName string) error {
	user := telectx.Message().Text
	from := telectx.Chat().Username

	// если пользователь прислал username с собакой
	if strings.HasPrefix(user, "@") {
		username := strings.TrimPrefix(user, "@")
		return c.handleUsername(ctx, telectx, username, from, spaceName)
	}

	// если пользователь прислал ссылку
	if strings.Contains(user, "t.me") {
		return c.handleUserLink(ctx, telectx, user, from, spaceName)
	}

	// если нет собаки, скорее всего это тоже username
	return c.handleUsername(ctx, telectx, user, from, spaceName)
}

// handleUsername обрабатывает username во время добавления пользователя как нового участника пространства
func (c *Controller) handleUsername(ctx context.Context, telectx tele.Context, username, from, spaceName string) error {
	toUser, err := c.userSrv.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, api_errors.ErrUserNotFound) {
			return telectx.EditOrSend(messages.UserNotRegisteredMessage, view.BackToMenuBtn())
		}

		return err
	}

	space := c.sharedSpace.CurrentSpace(telectx.Chat().ID)

	fromUser := model.Participant{
		User: model.User{
			TGID: telectx.Chat().ID,
		},
	}

	to := model.Participant{
		User: model.User{
			TGID:     toUser.TGID,
			Username: username,
		},
		State: model.PendingState,
	}

	err = c.sendInvitation(from, spaceName, to.TGID, fromUser.TGID)
	if err != nil {
		return err
	}

	return c.processInvitation(ctx, telectx, fromUser, to, space)
}

func (c *Controller) processInvitation(ctx context.Context, telectx tele.Context, from model.Participant, to model.Participant, space model.SharedSpace) error {
	err := c.sharedSpace.ProcessInvitation(ctx, from, to, int64(space.ID))
	if err != nil {
		if errors.Is(err, api_errors.ErrInvitationExists) {
			return telectx.EditOrSend(messages.UserAlreadyInvitedMessage, view.BackToMenuBtn())
		}

		if errors.Is(err, api_errors.ErrUserAlreadyExists) {
			msg := fmt.Sprintf(messages.UserAlreadyExistsMessage, space.Name)
			return telectx.EditOrSend(msg, view.BackToMenuBtn())
		}

		return fmt.Errorf("error processing invitation: %+v", err)
	}

	c.bot.Handle(&view.BtnAcceptInvitations, func(telectx tele.Context) error {
		return c.AcceptInvitation(ctx, telectx, from, to, space)
	})

	c.bot.Handle(&view.BtnDenyInvitations, func(telectx tele.Context) error {
		return c.DenyInvitation(ctx, telectx, from, to, space)
	})

	return telectx.EditOrSend(messages.SuccessfullySentInvitationsMessage, view.BackToMenuBtn())
}

// handleUserLink обрабатывает ссылку на пользователя во время добавления его как нового участника пространства
func (c *Controller) handleUserLink(ctx context.Context, telectx tele.Context, user, from, spaceName string) error {
	// https://t.me/iloveprogramming

	link := strings.Split(user, "/")

	if len(link) != 4 {
		return telectx.EditOrSend(messages.InvalidUserLinkMessage, view.BackToMenuBtn())
	}

	username := link[3]

	return c.handleUsername(ctx, telectx, username, from, spaceName)
}

func (c *Controller) SpaceName(ctx context.Context, telectx tele.Context) string {
	return c.sharedSpace.CurrentSpaceName(telectx.Chat().ID)
}

// SharedSpaceNewNote запрашивает у пользователя текст новой заметки для добавления в shared space
func (c *Controller) SharedSpaceNewNote(ctx context.Context, telectx tele.Context) error {
	return telectx.EditOrSend(messages.AskNoteTextMessage, view.BackToMenuBtn())
}

// AddNoteToSharedSpace принимает текст новой заметки и сохраняет ее в хранилище
func (c *Controller) AddNoteToSharedSpace(ctx context.Context, telectx tele.Context) error {
	note := model.Note{
		Text: telectx.Message().Text,
		Creator: model.User{
			TGID:     telectx.Chat().ID,
			Username: telectx.Chat().Username,
		},
	}

	err := c.sharedSpace.SaveNote(ctx, note)
	if err != nil {
		return err
	}

	// рассылка уведомлений участникам пространства
	participants := c.sharedSpace.SpaceParticipants(telectx.Chat().ID)
	creator := c.sharedSpace.SpaceCreator(telectx.Chat().ID)
	spaceName := c.SpaceName(ctx, telectx)

	for _, user := range participants {
		if user.TGID != creator.TGID {
			msg := fmt.Sprintf(messages.UserAddedNoteMessage, telectx.Chat().Username, spaceName)
			_, err = c.bot.Send(&tele.User{ID: user.TGID}, msg, view.ShowSharedSpacesMenu())
			if err != nil {
				return err
			}
		}
	}

	msg := fmt.Sprintf(messages.SuccessfullyAddedNoteMessage, spaceName)
	return telectx.EditOrSend(msg, view.BackToSharedSpaceMenu())
}

// NextPageNotes обрабатывает кнопку переключения на следующую страницу заметок в совместном пространстве
func (c *Controller) NextPageNotesSharedSpace(ctx context.Context, telectx tele.Context) error {
	next, kb := c.sharedSpace.NextPageNotes(telectx.Chat().ID)

	return telectx.Edit(next, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})
}

// NextPageNotes обрабатывает кнопку переключения на предыдущую страницу заметок в совместном пространстве
func (c *Controller) PrevPageNotesSharedSpace(ctx context.Context, telectx tele.Context) error {
	next, kb := c.sharedSpace.PrevPageNotes(telectx.Chat().ID)

	return telectx.Edit(next, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

}

// NextPageNotes обрабатывает кнопку переключения на последнюю страницу заметок в совместном пространстве
func (c *Controller) LastPageNotesSharedSpace(ctx context.Context, telectx tele.Context) error {
	next, kb := c.sharedSpace.LastPageNotes(telectx.Chat().ID)

	err := telectx.Edit(next, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей
	if err != nil {
		return checkError(err)
	}

	return nil
}

// NextPageNotes обрабатывает кнопку переключения на первую страницу заметок в совместном пространстве
func (c *Controller) FirstPageNotesSharedSpace(ctx context.Context, telectx tele.Context) error {
	next, kb := c.sharedSpace.FirstPageNotes(telectx.Chat().ID)

	err := telectx.Edit(next, &tele.SendOptions{
		ReplyMarkup: kb,
		ParseMode:   htmlParseMode,
	})

	// если пришла ошибка о том, что сообщение не изменено - игнорируем.
	// такая ошибка происходит, если быть на первой странице и нажать кнопку "первая страница".
	// то же самое происходит и с последней страницей

	if err != nil {
		return checkError(err)
	}

	return nil
}

// HandleContact обрабатывает контакт при создании приглашения в совместное пространство
func (c *Controller) HandleContact(ctx context.Context, telectx tele.Context) error {
	contact := telectx.Message().Contact

	// проверяем, зарегистрирован ли пользователь в боте
	exists := c.userSrv.CheckUser(ctx, contact.UserID)
	if !exists {
		return telectx.EditOrSend(messages.UserNotRegisteredMessage, view.BackToMenuBtn())
	}

	space := c.sharedSpace.CurrentSpace(telectx.Chat().ID)

	err := c.sendInvitation(telectx.Chat().Username, space.Name, contact.UserID, telectx.Chat().ID)
	if err != nil {
		return fmt.Errorf("error sending invitation by contact: %+v", err)
	}

	fromUser := model.Participant{
		User: model.User{
			TGID:     telectx.Chat().ID,
			Username: formatUsername(telectx.Chat()),
		},
	}

	to := model.Participant{
		User: model.User{
			TGID:     contact.UserID,
			Username: formatUsername(&tele.Chat{ID: contact.UserID, FirstName: contact.FirstName, LastName: contact.LastName}),
		},
		State: model.PendingState,
	}

	return c.processInvitation(ctx, telectx, fromUser, to, space)
}

// AcceptInvitation обрабатывает кнопку согласия вступить в совместное пространство
func (c *Controller) AcceptInvitation(ctx context.Context, telectx tele.Context, from model.Participant, to model.Participant, space model.SharedSpace) error {
	// отправить пользователю, который пригласил, уведомление о том, что второй пользователь согласился

	username := formatUsername(telectx.Chat())

	msg := fmt.Sprintf(messages.UserAcceptedInvitationMessage, username, space.Name)
	err := c.sendMessage(from.TGID, msg, view.ShowSharedSpacesMenu())
	if err != nil {
		return err
	}

	// удалить приглашение из БД и обновить state у приглашенного пользователя
	err = c.sharedSpace.AcceptInvitation(ctx, from, to, int64(space.ID))
	if err != nil {
		return err
	}

	// отправить уведомление участникам беседы о том, что пользователь добавил другого пользователя
	msg = fmt.Sprintf(messages.UserWasAdded, from.Username, to.Username, space.Name)

	for _, user := range space.Participants {
		// отправляем всем, кроме того, кто пригласил
		if user.ID != int(from.TGID) {
			_, err := c.bot.Send(&tele.Chat{ID: user.TGID}, msg)
			if err != nil {
				c.HandleError(telectx, err, "notify_partisipants_new_user")
			}
		}

	}

	msg = fmt.Sprintf(messages.InvitationAcceptedMessage, space.Name)
	return telectx.EditOrSend(msg, view.ShowSharedSpacesMenu())
}

// AcceptInvitation обрабатывает кнопку отказаться вступить в совместное пространство
func (c *Controller) DenyInvitation(ctx context.Context, telectx tele.Context, from model.Participant, to model.Participant, space model.SharedSpace) error {
	// отправить пользователю, который пригласил, уведомление о том, что второй пользователь отказался
	username := formatUsername(telectx.Chat())

	msg := fmt.Sprintf(messages.UserRejecteddInvitationMessage, username, space.Name)
	err := c.sendMessage(from.TGID, msg, view.ShowSharedSpacesMenu())
	if err != nil {
		return err
	}

	// удалить приглашение и участника из БД
	err = c.sharedSpace.DenyInvitation(ctx, from, to, int64(space.ID))
	if err != nil {
		return err
	}

	msg = fmt.Sprintf(messages.InvitationRejectedMessage, space.Name)
	return telectx.EditOrSend(msg, view.ShowSharedSpacesMenu())
}

func (c *Controller) sendInvitation(from, spaceName string, toID, fromID int64) error {
	msg := fmt.Sprintf(messages.InvitationsMessage, from, spaceName)

	return c.sendMessage(toID, msg, view.InvintationKeyboard())
}

func formatUsername(chat *tele.Chat) string {
	if chat.Username != "" {
		return fmt.Sprintf("@%s", chat.Username)
	}

	return fmt.Sprintf("%s %s", chat.FirstName, chat.LastName)
}

func (c *Controller) AddReminderSharedSpace(ctx context.Context, telectx tele.Context) error {
	return nil
}
