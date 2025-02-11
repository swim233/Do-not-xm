package xmchecker

import (
	"learn/units/bot"
	"learn/units/counter"
	"math/rand/v2"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// 检查
func XmChecker(update tgbotapi.Update) bool {
	if update.Message != nil && bot.Mode == "match" && counter.CheckSleep() && bot.Mode != "off" {
		if update.Message.From.ID == bot.BotConfig.IntUserID {
			bot.CheckFlag = update.Message.MessageID
		}
		if ((update.Message.MessageID == (bot.CheckFlag + 1)) || ((update.Message.ReplyToMessage != nil) && (update.Message.ReplyToMessage.From.ID == bot.BotConfig.IntUserID))) && IsXm(update.Message.Text) {
			return true
		}
	}
	if update.Message != nil && bot.Mode == "any" && IsXm(update.Message.Text) && (update.Message.From.ID != bot.BotConfig.IntUserID) && counter.CheckSleep() && bot.Mode != "off" {
		return true
	}
	return false
}

func IsXm(update string) bool {
	if bot.CheckXm.MatchString(update) || strings.Contains(update, "xm") {
		return true
	}
	return false
}

func SendXm(update tgbotapi.Update) error {
	*counter.Time = (rand.IntN(bot.BotConfig.RandomCD) + bot.BotConfig.StaticCD)

	msgID := update.Message.MessageID

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "不许羡慕！")
	// handler.RecordLastXm(&update)

	msg.ReplyToMessageID = msgID

	bot.Bot.Send(msg)

	return nil
}
