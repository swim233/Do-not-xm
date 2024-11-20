package pinghandler

import (
	"learn/units/bot"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func PingHandler(update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Pong！")
	bot.Bot.Send(msg)
	return nil
}
