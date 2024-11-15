package main

import (
	_ "encoding/json"
	"fmt"
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

type Data struct {
	ChatID int
}
type Config struct {
	Token     string
	UserID    string
	intUserID int64
	randomCD  int
	staticCD  int
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
	if err != nil {
		logger.Error("%s", err)
	}
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
	b.NewCommandProcessor("switchmode", switchmodeHandler)
	b.NewCommandProcessor("changecd", changecdHandler)
	b.NewProcessor(func(update tgbotapi.Update) bool {
		if update.Message != nil && Mode == "match" && Sleep <= 0 {
			if update.Message.From.ID == BotConifg.intUserID {
				CheckFlag = update.Message.MessageID
			}
			if ((update.Message.MessageID == (CheckFlag + 1)) || ((update.Message.ReplyToMessage != nil) && (update.Message.ReplyToMessage.From.ID == BotConifg.intUserID))) && IsXm(update.Message.Text) {
				Sleep = (rand.IntN(BotConifg.randomCD) + BotConifg.staticCD)
				return true
			} else {
				time.Sleep(1 * time.Second)
				Sleep--
			}
		}
		if update.Message != nil && Mode == "any" && IsXm(update.Message.Text) && (update.Message.From.ID != BotConifg.intUserID) && Sleep <= 0 {
			Sleep = (rand.IntN(BotConifg.randomCD) + BotConifg.staticCD)
			return true
		} else {
			time.Sleep(1 * time.Second)
			Sleep--

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

func switchmodeHandler(update tgbotapi.Update) error {
	if update.Message.From.ID == BotConifg.intUserID {
		if Mode == "match" {
			Mode = "any"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前模式为：全局匹配")
			Bot.Send(msg)
		} else {
			Mode = "match"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "当前模式为：匹配模式")
			Bot.Send(msg)
		}
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你没有使用该命令的权限！")
		Bot.Send(msg)
	}

	return nil
}
func changecdHandler(update tgbotapi.Update) error {
	if update.Message.From.ID == BotConifg.intUserID {
		args := strings.Split(update.Message.CommandArguments(), " ")
		if len(args) != 2 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "格式有误")
			Bot.Send(msg)
			return fmt.Errorf("参数异常")
		}
		//获取群组和用户id
		staticCDstr := args[0]
		randomCDstr := args[1]
		staticCD, err := strconv.ParseInt(staticCDstr, 10, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "格式有误")
			Bot.Send(msg)
			return fmt.Errorf("格式有误")
		}
		BotConifg.staticCD = int(staticCD)

		randomCD, err := strconv.ParseInt(randomCDstr, 10, 64)
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "格式有误")
			Bot.Send(msg)
			return fmt.Errorf("格式有误")

		}
		BotConifg.randomCD = int(randomCD)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("修改成功 当前cd为%ds固定cd+%ds随机cd", staticCD, randomCD))
		Bot.Send(msg)
		return nil
	} else {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你没有使用该命令的权限！")
		Bot.Send(msg)
		return nil
	}
}

//TODO
//添加指令修改休眠时间
//拆分文件
//添加独立的群组模式
//羡慕次数统计
