package main

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/bot"
	"github.com/swim233/do_not_xm/config"
	"github.com/swim233/do_not_xm/handler"
	"github.com/swim233/do_not_xm/logger"
)

func main() {
	logger.InitZap()
	config.InitDB()
	config.InitConfig()
	bot.InitBot()
	b := bot.Bot.AddHandle()
	b.NewProcessor(func(u tgbotapi.Update) bool {
		return u.Message != nil
	}, func(update tgbotapi.Update) error {
		return handler.MessageHandler(update)
	})
	b.Run()
}
