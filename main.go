package main

import (
	"learn/goUnits/logger"
	"log"
	"math/rand/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	godotenv "github.com/joho/godotenv"
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
	logger.SetLogLevel(1)
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// 如果 .env 文件不存在，创建并写入默认值
		logger.Info(".env 文件不存在，正在创建...")

		// 创建并打开 .env 文件
		file, err := os.Create(".env")
		if err != nil {
			log.Fatalf("创建 .env 文件失败: %v", err)
		}
		defer file.Close()

		// 写入默认的环境变量内容
		defaultEnv := `Toekn=
UserID=
`
		if _, err := file.WriteString(defaultEnv); err != nil {
			log.Fatalf("写入 .env 文件失败: %v", err)
		}
		logger.Info(".env 文件已创建，并写入默认内容.")
	}
	err := godotenv.Load()
	if err != nil {
		logger.Error("%s", err)
	}

	BotConifg.Token = os.Getenv("Token")
	BotConifg.UserID = os.Getenv("UserID")
	BotConifg.intUserID, err = strconv.ParseInt(BotConifg.UserID, 10, 64)
	qwq, err := tgbotapi.NewBotAPI(BotConifg.Token)
	Bot = qwq
	if err != nil {
		log.Printf("%s", BotConifg.Token)
		log.Printf("%s", err)
	}
	if err != nil {
		logger.Error("%s", err)
	}
	Bot.Debug = true
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
