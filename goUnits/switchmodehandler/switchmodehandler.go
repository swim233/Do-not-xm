package switchmodehandler

import (
	"learn/goUnits/bot"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func SwitchModeHandler(update tgbotapi.Update) error {
	if update.Message.From.ID == bot.BotConfig.IntUserID || bot.BotConfig.DebugFlag {
		if bot.Mode == "match" {
			bot.Mode = "any"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前模式为: 全局匹配")
			bot.Bot.Send(msg)
		} else {
			bot.Mode = "match"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前模式为: 匹配模式")
			bot.Bot.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你没有使用该命令的权限！")
		bot.Bot.Send(msg)
	}

	return nil
}
