package debughandler

import (
	"fmt"
	"learn/units/bot"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func DebugHandler(update tgbotapi.Update) error {
	if update.Message.From.ID == bot.BotConfig.IntUserID {
		bot.BotConfig.DebugFlag = !bot.BotConfig.DebugFlag
		fmtmsg := fmt.Sprintf("Debug模式当前为: %t", bot.BotConfig.DebugFlag)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmtmsg)
		bot.Bot.Send(msg)
	}
	return nil
}
