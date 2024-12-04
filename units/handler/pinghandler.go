package handler

import (
	"learn/units/bot"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// 发送ping
func PingHandler(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pong！")
	bot.Bot.Send(msg)
	return nil
}
