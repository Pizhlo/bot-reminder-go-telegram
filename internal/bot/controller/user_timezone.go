package controller

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/zsefvlol/timezonemapper"
	"gopkg.in/telebot.v3"
)

func (c *Controller) parseTimezone(id int64, loc *telebot.Location) model.UserTimezone {
	tz := model.UserTimezone{
		Lat:  loc.Lat,
		Long: loc.Lng,
	}

	timezone := timezonemapper.LatLngToTimezoneString(39.9254474, 116.3870752)

	fmt.Printf("Timezone: %s\n", timezone)

	tz.Location = timezone

	return tz
}
