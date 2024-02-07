package navigation

import (
	"sync"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/logger"
	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/view"
	"github.com/sirupsen/logrus"
	"gopkg.in/telebot.v3"
)

// NavigationService отвечает управление меню
type NavigationService struct {
	logger   *logrus.Logger
	viewsMap map[int64]*view.NavigationView
	mu       sync.Mutex
}

func New() *NavigationService {
	return &NavigationService{logger: logger.New(), viewsMap: make(map[int64]*view.NavigationView), mu: sync.Mutex{}}
}

// SaveUser сохраняет пользователя в мапе view
func (n *NavigationService) SaveUser(userID int64) {
	n.logger.Debugf("Navigation service: checking if user saved in the views map...\n")
	if _, ok := n.viewsMap[userID]; !ok {
		n.logger.Debugf("Navigation service: user not found in the views map. Saving...\n")
		n.viewsMap[userID] = view.NewNavigation()
	} else {
		n.logger.Debugf("Navigation service: user already saved in the views map.\n")
	}

	n.logger.Debugf("Navigation service: successfully saved user in the views map.\n")
}

// MainMenu возвращает главное меню.
// Кнопки: Профиль, Настройки, Заметки, Напоминания
func (n *NavigationService) MainMenu(userID int64) *telebot.ReplyMarkup {
	n.mu.Lock()
	defer n.mu.Unlock()

	menu := n.viewsMap[userID].MainMenu()
	return menu
}
