package main

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/utils"
	"github.com/swim233/do_not_xm/utils/handler"
)

func main() {
	utils.InitBot() //初始化bot
	b := utils.Bot.AddHandle()
	processor := handler.NewMessageProcessor()
	b.NewCommandProcessor("changecd", processor.ChangeCoolDown)
	b.NewCommandProcessor("cd", processor.CD)
	b.NewProcessor(func(u tgbotapi.Update) bool {
		return u.Message != nil && !u.Message.IsCommand()
	}, processor.Processor)
	b.Run()
}
