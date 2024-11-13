package main

import (
	"log"
	"math/rand/v2"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

var CheckFlag = 0
var Sleep = 0
var CheckXm = regexp.MustCompile(".*羡.*慕.*")
var Mode = "match"
var Bot *tgbotapi.BotAPI

type Config struct {
	Token     string
	UserID    string
	intUserID int64
}

var BotConifg Config

func main() {
	BotConifg.Token = ""
	BotConifg.intUserID = 123
	qwq, err := tgbotapi.NewBotAPI(BotConifg.Token)
	Bot = qwq
	if err != nil {
		log.Printf("%s", BotConifg.Token)
		log.Printf("%s", err)
	}
	b := Bot.AddHandle()
	b.NewCommandProcessor("switchmode", switchmodeHandle)
	b.NewProcessor(func(update tgbotapi.Update) bool {
		if update.Message != nil && Mode == "match" && Sleep <= 0 {
			if update.Message.From.ID == BotConifg.intUserID {
				CheckFlag = update.Message.MessageID
			}
			if ((update.Message.MessageID == (CheckFlag + 1)) || ((update.Message.ReplyToMessage != nil) && (update.Message.ReplyToMessage.From.ID == BotConifg.intUserID))) && IsXm(update.Message.Text) {
				Sleep = (rand.IntN(10) + 10)
				return true
			} else {
				time.Sleep(1 * time.Second)
				Sleep--
			}
		}
		if update.Message != nil && Mode == "any" && IsXm(update.Message.Text) {
			return true
		}
		return false
	}, sendXm)
	b.Run()

}
func IsXm(update string) bool {
	if CheckXm.MatchString(update) || strings.Contains(update, "xm") {
		return true
	}
	return false
}

func sendXm(update tgbotapi.Update) error {
	msgID := update.Message.MessageID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "不许羡慕！")
	msg.ReplyToMessageID = msgID
	Bot.Send(msg)
	return nil
}

func switchmodeHandle(update tgbotapi.Update) error {
	if Mode == "match" {
		Mode = "any"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前模式为：全局匹配")
		Bot.Send(msg)
	} else {
		Mode = "match"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前模式为：匹配模式")
		Bot.Send(msg)
	}
	return nil
}
