package handler

import (
	"learn/units/bot"
	"learn/units/counter"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func SwitchModeHandler(update tgbotapi.Update) error {
	if update.Message.From.ID == bot.BotConfig.IntUserID || bot.BotConfig.DebugFlag {
		args := strings.Split(update.Message.CommandArguments(), " ")
		if len(args) != 1 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "参数数量有误")
			bot.Bot.Send(msg)
		} else {
			if args[0] == "any" {
				bot.Mode = "any"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前模式为: 全局匹配")
				bot.Bot.Send(msg)
			} else if args[0] == "match" {
				bot.Mode = "match"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前模式为: 匹配模式")
				bot.Bot.Send(msg)
			} else if args[0] == "off" {
				bot.Mode = "off"
				*counter.Time = 0
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前bot已关闭")
				bot.Bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.From.ID, "参数有误")
				bot.Bot.Send(msg)
			}

		}

	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你没有使用该命令的权限！")
		bot.Bot.Send(msg)
	}

	return nil
}
