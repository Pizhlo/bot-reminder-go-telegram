package controller

import "gopkg.in/telebot.v3"

type TGID struct {
	telebot.Context
}

func (p TGID) Get() int {
	return int(p.Sender().ID)
}

type tgid = TGID
