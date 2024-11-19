package debughandler

import (
	"fmt"
	"learn/goUnits/bot"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func DebugHandler(update tgbotapi.Update) error {
	if update.Message.From.ID == bot.BotConifg.IntUserID {
		bot.BotConifg.DebugFlag = !bot.BotConifg.DebugFlag
		fmtmsg := fmt.Sprintf("Debug模式当前为: %t", bot.BotConifg.DebugFlag)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmtmsg)
		bot.Bot.Send(msg)
	}
	return nil
}
