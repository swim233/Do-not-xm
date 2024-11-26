package handler

import (
	"fmt"
	"learn/units/bot"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

var LastXmMessageID = 0

func LastXmHandler(update tgbotapi.Update) error {
	fmtmsg := fmt.Sprintf("上一条xm为：https://t.me/ArknightsZH/%d", LastXmMessageID)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmtmsg)
	bot.Bot.Send(msg)
	return nil
}
