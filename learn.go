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

type Config struct {
	Token     string
	UserID    string
	intUserID int64
}

var BotConifg Config

func main() {
	// 创建Bot实例

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
		log.Panic(err)
	}
	BotConifg.Token = os.Getenv("Token")
	BotConifg.UserID = os.Getenv("UserID")
	BotConifg.intUserID, err = strconv.ParseInt(BotConifg.UserID, 10, 64)
	if err != nil {
		log.Panic(err)
	}
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

			if update.Message.From.ID == BotConifg.intUserID {
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
