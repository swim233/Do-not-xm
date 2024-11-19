package bot

import (
	"learn/units/logger"
	"log"
	"os"
	"regexp"
	"strconv"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	godotenv "github.com/joho/godotenv"
)

var Bot *tgbotapi.BotAPI

type Config struct {
	Token     string
	UserID    string
	IntUserID int64
	RandomCD  int
	StaticCD  int
	DebugFlag bool
}

var (
	CheckFlag = 0
	CheckXm   = regexp.MustCompile(".*羡.*慕.*")
	Mode      = "match"
)
var BotConfig Config

func InitBot() {
	BotConfig.RandomCD = 30
	BotConfig.StaticCD = 30
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
		defaultEnv := `Token=
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

	BotConfig.Token = os.Getenv("Token")
	BotConfig.UserID = os.Getenv("UserID")
	BotConfig.IntUserID, err = strconv.ParseInt(BotConfig.UserID, 10, 64)
	if err != nil {
		logger.Error("%s", err)
	}
	qwq, err := tgbotapi.NewBotAPI(BotConfig.Token)
	Bot = qwq
	if err != nil {
		log.Printf("%s", BotConfig.Token)
		log.Printf("%s", err)
	}
	if err != nil {
		logger.Error("%s", err)
	}
}
