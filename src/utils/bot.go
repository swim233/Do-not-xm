package utils

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/swim233/do_not_xm/utils/logger"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	godotenv "github.com/joho/godotenv"
)

var Bot *tgbotapi.BotAPI

type Config struct {
	Token            string // Telegram Bot Token
	DebugFlag        bool   // 是否开启debug输出
	ApiLogLevel      int    // 日志等级
	EnableChecker    bool   //是否开启羡慕检查
	RandomCoolDown   int    //随机冷却时间
	StaticCoolDown   int    //固定冷却时间
	TriggerMode      string //触发模式
	WaitingForUserID int64  //监听用户id

}

var BotConfig Config

func InitBot() {

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// 如果 .env 文件不存在，创建并写入默认值
		logger.Info(".env 文件不存在，正在创建...")

		// 创建并打开 .env 文件
		file, err := os.Create(".env")
		if err != nil {
			logger.Error("创建 .env 文件失败: %v", err)
		}
		defer file.Close()

		// 写入默认的环境变量内容
		defaultEnv := `# Telegram Bot Token
Token=YOUR_TOKEN_ID

# 是否开启BotAPI debug输出(true/false)
DebugFlag=true

# API 日志等级 (可选值: DEBUG, INFO, WARN, ERROR)
ApiLogLevel=DEBUG/INFO/WARN/ERROR

# 是否开启羡慕检查
EnableChecker=true   

# 随机冷却时间
RandomCoolDown=100  

# 固定冷却时间
StaticCoolDown=100
	
# 匹配模式
TriggerMode=any

# 监听用户id
WaitingForUserID=UserID

`
		if _, err := file.WriteString(defaultEnv); err != nil {
			logger.Error("写入 .env 文件失败: %v", err)
		}
		logger.Info(".env 文件已创建，并写入默认内容.")
		os.Exit(1)
	}
	err := godotenv.Load()
	if err != nil {
		logger.Error("%s", err)
	}

	getEnv() //读取环境变量

	proxy := FetchProxy()
	if proxy != "" {
		proxyURL, err := url.Parse(proxy)
		if err != nil {
			logger.Error("Failed to parse proxy url: %s", proxy)
			return
		}
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}
		Bot.Client = client
		logger.Info("Using proxy: %s", proxy)
	}
}

func FetchProxy() string {
	proxy := os.Getenv("HTTP_PROXY")
	if proxy == "" {
		proxy = os.Getenv("HTTPS_PROXY")
	}
	if proxy == "" {
		proxy = os.Getenv("http_proxy")
	}
	if proxy == "" {
		proxy = os.Getenv("https_proxy")
	}
	return proxy
}
func UpdateEnvValue(key, newValue string) error {
	// 读取现有内容
	content, err := os.ReadFile(".env")
	if err != nil {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	found := false

	// 更新值
	for i, line := range lines {
		if strings.HasPrefix(line, key+"=") {
			lines[i] = key + "=" + newValue
			found = true
			break
		}
	}

	// 如果键不存在，添加新行
	if !found {
		lines = append(lines, key+"="+newValue)
	}

	// 写回文件
	output := strings.Join(lines, "\n")
	err = os.WriteFile(".env", []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}
	return nil
}

func getEnv() {
	var err error

	BotConfig.DebugFlag = (os.Getenv("DebugFlag") == "true") //读取debug输出

	BotConfig.Token = os.Getenv("Token") //读取token

	bot, err := tgbotapi.NewBotAPI(BotConfig.Token) //实例化BotAPI
	if err != nil {
		logger.Error("%s", err)
		err = nil
	}
	Bot = bot
	if err != nil {
		logger.Error("%s", BotConfig.Token)
		logger.Error("%s", err)
		err = nil
	}
	Bot.Debug = BotConfig.DebugFlag

	loglevel := logger.ParseLogLevel(os.Getenv("LogLevel")) //读取bot log level
	logger.SetLogLevel(loglevel)

	BotConfig.ApiLogLevel = logger.ParseLogLevel(os.Getenv("ApiLogLevel")) //读取bot api log level
	adapter := &logger.TelegramBotApiLoggerAdapter{}
	adapter.SetLogger(logger.GetInstance())
	adapter.SetLogLevel(BotConfig.ApiLogLevel)
	tgbotapi.SetLogger(adapter)

	BotConfig.TriggerMode = os.Getenv("TriggerMode") //读取匹配模式

	userIDStr := os.Getenv("WaitingForUserID") //读取监听用户id
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid WaitingForUserID: %s", userIDStr)
		userID = 0
	}
	BotConfig.WaitingForUserID = userID
}
