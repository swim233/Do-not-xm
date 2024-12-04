package handler

import (
	"learn/units/bot"
	"learn/units/counter"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// 发送cd信息
func CdHandler(update tgbotapi.Update) error {
	fmtmsg := "当前剩余CD： " + counter.Calculation(counter.Time)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmtmsg)
	bot.Bot.Send(msg)
	return nil
}
