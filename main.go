package main

import (
	_ "encoding/json"
	bot "learn/goUnits/bot"
	"learn/goUnits/changecdhandler"
	"learn/goUnits/debughandler"
	"learn/goUnits/logger"
	"learn/goUnits/switchmodehandler"
	"math/rand/v2"
	"strings"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

type Data struct {
	ChatID int
}

func main() {
	bot.InitBot()
	logger.SetLogLevel(1)
	bot.Bot.Debug = true
	b := bot.Bot.AddHandle()
	b.NewCommandProcessor("switchmode", switchmodehandler.SwitchModeHandler)
	b.NewCommandProcessor("changecd", changecdhandler.ChangeCdHandler)
	b.NewCommandProcessor("debug", debughandler.DebugHandler)
	b.NewProcessor(func(update tgbotapi.Update) bool {
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
	}, sendXm)
	b.Run()
}

func IsXm(update string) bool {
	if bot.CheckXm.MatchString(update) || strings.Contains(update, "xm") {
		return true
	}
	return false
}

func sendXm(update tgbotapi.Update) error {
	msgID := update.Message.MessageID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "不许羡慕！")
	msg.ReplyToMessageID = msgID
	bot.Bot.Send(msg)
	return nil
}

// TODO
// 拆分文件
// 添加独立的群组模式
// 羡慕次数统计
