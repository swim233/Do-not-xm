package bot

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/config"
	"github.com/swim233/do_not_xm/logger"
)

var Bot *tgbotapi.BotAPI

func InitBot() {

	bot, err := tgbotapi.NewBotAPI(config.Config.GetString("bot_token"))
	if err != nil {
		logger.Suger.Errorf("Fail to Init Bot : %s", err.Error())
	} else {
		logger.Suger.Infof("Succeed login,bot token : %s", config.Config.GetString("bot_token"))
		Bot = bot
	}
	if config.Config.GetBool("debug_mode") {
		Bot.Debug = true
	}
	tgbotapi.SetLogger(&logger.ZapBotLogger{})
}
