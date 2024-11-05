package main

import (
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

var CheckFlag = 0

func main() {
	// 创建Bot实例
	Bot, err := tgbotapi.NewBotAPI("YOUR_BOT_TOKEN")
	if err != nil {
		log.Panic(err)
	}

	Bot.Debug = true
	updatecfg := tgbotapi.NewUpdate(0)
	updatecfg.Timeout = 60
	updates := Bot.GetUpdatesChan(updatecfg)
	checkXm := regexp.MustCompile(".*羡.*慕.*")
	for update := range updates {
		if update.Message != nil {

			if update.Message.From.ID == 5568996608 {
				CheckFlag = update.Message.MessageID
			}
			if ((update.Message.MessageID == (CheckFlag + 1)) && (checkXm.MatchString(update.Message.Text) || strings.Contains(update.Message.Text, "xm"))) || (update.Message.ReplyToMessage != nil) && (update.Message.ReplyToMessage.From.ID == 5568996608 && (checkXm.MatchString(update.Message.Text) || strings.Contains(update.Message.Text, "xm"))) {
				msgID := update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "不许羡慕！")
				msg.ReplyToMessageID = msgID
				Bot.Send(msg)
			}
		}
	}
}
