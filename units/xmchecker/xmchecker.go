package xmchecker

import (
	"learn/units/bot"
	"learn/units/timer"
	"math/rand/v2"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func XmChecker(update tgbotapi.Update) bool {
	if update.Message != nil && bot.Mode == "match" && timer.CheckSleep() {
		if update.Message.From.ID == bot.BotConfig.IntUserID {
			bot.CheckFlag = update.Message.MessageID
		}
		if ((update.Message.MessageID == (bot.CheckFlag + 1)) || ((update.Message.ReplyToMessage != nil) && (update.Message.ReplyToMessage.From.ID == bot.BotConfig.IntUserID))) && IsXm(update.Message.Text) {
			return true
		}
	}
	if update.Message != nil && bot.Mode == "any" && IsXm(update.Message.Text) && (update.Message.From.ID != bot.BotConfig.IntUserID) && timer.CheckSleep() {

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
	*timer.Time = (rand.IntN(bot.BotConfig.RandomCD) + bot.BotConfig.StaticCD)
	msgID := update.Message.MessageID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "不许羡慕！")
	msg.ReplyToMessageID = msgID
	bot.Bot.Send(msg)
	return nil
}
