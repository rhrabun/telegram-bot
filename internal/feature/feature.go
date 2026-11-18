package feature

import "github.com/go-telegram/bot"

type Feature interface {
	Name() string
	Register(b *bot.Bot) error
}
