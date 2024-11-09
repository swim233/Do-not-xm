package main

import (
	"learn/goUnits/logger"
	"log"
	"math/rand/v2"
	"os"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	godotenv "github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
)

var CheckFlag = 0
var Sleep = 0
var CheckXm = regexp.MustCompile(".*羡.*慕.*")

type Config struct {
	Token string
}

var BotConifg Config

func main() {
	// 创建Bot实例
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error:%s", err)
	}
	BotConifg.Token = os.Getenv("Token")
	Bot, err := tgbotapi.NewBotAPI(BotConifg.Token)
	if err != nil {
		log.Panic(err)
	}

	Bot.Debug = true
	updatecfg := tgbotapi.NewUpdate(0)
	updatecfg.Timeout = 60
	updates := Bot.GetUpdatesChan(updatecfg)

	for update := range updates {
		if update.Message != nil {

			if update.Message.From.ID == 5568996608 {
				CheckFlag = update.Message.MessageID
			}
			if Sleep <= 0 {
				if ((update.Message.MessageID == (CheckFlag + 1)) || ((update.Message.ReplyToMessage != nil) && (update.Message.ReplyToMessage.From.ID == 5568996608))) && IsXm(update.Message.Text) {
					msgID := update.Message.MessageID
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "不许羡慕！")
					msg.ReplyToMessageID = msgID
					Bot.Send(msg)
					Sleep = (rand.IntN(10) + 10)
				} else {
					time.Sleep(1 * time.Second)
					Sleep--
				}
			}
		}
	}

}
func IsXm(update string) bool {
	if CheckXm.MatchString(update) || strings.Contains(update, "xm") {
		return true
	}
	return false
}
