package main

import (
	"flag"

	bot "github.com/Pizhlo/bot-reminder-go-telegram/cmd/bot"
)

func main() {
	fileName := flag.String("filename", ".env", "name of config file")
	path := flag.String("path", ".", "path to config file")

	flag.Parse()

	bot.Start(*fileName, *path)
}
