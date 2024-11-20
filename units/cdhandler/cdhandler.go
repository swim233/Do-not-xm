package cdhandler

import (
	"fmt"
	"learn/units/bot"
	"learn/units/timer"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func CdHandler(update tgbotapi.Update) error {
	fmtmsg := fmt.Sprintf("当前剩余CD:%ds", *timer.Time)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmtmsg)
	bot.Bot.Send(msg)
	return nil
}
