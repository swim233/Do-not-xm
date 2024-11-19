package xmchecker

import (
	"learn/goUnits/bot"
	"math/rand/v2"
	"strings"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func XmChecker(update tgbotapi.Update) bool {
	if update.Message != nil && bot.Mode == "match" && bot.Sleep <= 0 {
		if update.Message.From.ID == bot.BotConifg.IntUserID {
			bot.CheckFlag = update.Message.MessageID
		}
		if ((update.Message.MessageID == (bot.CheckFlag + 1)) || ((update.Message.ReplyToMessage != nil) && (update.Message.ReplyToMessage.From.ID == bot.BotConifg.IntUserID))) && IsXm(update.Message.Text) {
			bot.Sleep = (rand.IntN(bot.BotConifg.RandomCD) + bot.BotConifg.StaticCD)
			return true
		} else {
			time.Sleep(1 * time.Second)
			bot.Sleep--
		}
	}
	if update.Message != nil && bot.Mode == "any" && IsXm(update.Message.Text) && (update.Message.From.ID != bot.BotConifg.IntUserID) && bot.Sleep <= 0 {
		bot.Sleep = (rand.IntN(bot.BotConifg.RandomCD) + bot.BotConifg.StaticCD)
		return true
	} else {
		time.Sleep(1 * time.Second)
		bot.Sleep--

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
	msgID := update.Message.MessageID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "不许羡慕！")
	msg.ReplyToMessageID = msgID
	bot.Bot.Send(msg)
	return nil
}
