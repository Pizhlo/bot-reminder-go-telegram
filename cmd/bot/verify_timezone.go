package bot

// import (
// 	"context"
// 	"errors"
// 	"log/slog"

// 	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/service/account"

// 	tele "gopkg.in/telebot.v3"
// )

// func NewVerifyTimezone(users account.Service, locDlg *LocationDialog) tele.MiddlewareFunc {
// 	const somethingWrongMsg = "Ой-ой, что-то пошло не так :-/"
// 	const registerMsg = "Вам надо зарегистрироваться, выполнив команду /start."

// 	return func(next tele.HandlerFunc) tele.HandlerFunc {
// 		return func(ctx tele.Context) error {
// 			tgid := ctx.Sender().ID

// 			u, err := users.FindUserByTelegramID(context.TODO(), int(tgid))
// 			if errors.Is(err, account.ErrUserNotFound) {
// 				return ctx.Send(registerMsg, tele.RemoveKeyboard)
// 			} else if err != nil {
// 				slog.Error("error verifying a location", "error", err)

// 				return ctx.Send(somethingWrongMsg, tele.RemoveKeyboard)
// 			}

// 			if !u.HasTimezone() {
// 				return locDlg.Open(ctx)
// 			}

// 			return next(ctx)
// 		}
// 	}
// }
