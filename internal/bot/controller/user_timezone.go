package controller

import (
	"fmt"

	"github.com/Pizhlo/bot-reminder-go-telegram/internal/bot/model"
	"github.com/zsefvlol/timezonemapper"
	"gopkg.in/telebot.v3"
)

func (c *Controller) saveTimezone(id int64, loc *telebot.Location) string {
	tz := model.UserTimezone{
		Lat:  loc.Lat,
		Long: loc.Lng,
	}

	timezone := timezonemapper.LatLngToTimezoneString(39.9254474, 116.3870752)
	// Should print "Timezone: Asia/Shanghai"
	fmt.Printf("Timezone: %s\n", timezone)

	tz.Location = timezone

	c.srv.TimezoneCacheEditor.SaveUserTimezone(id, tz)

	return timezone
}
