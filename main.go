package main

import (
	_ "encoding/json"
	bot "learn/goUnits/bot"
	"learn/goUnits/changecdhandler"
	"learn/goUnits/debughandler"
	"learn/goUnits/logger"
	"learn/goUnits/switchmodehandler"
	"learn/goUnits/xmchecker"
)

type Data struct {
	ChatID int
}

func main() {
	bot.InitBot()
	logger.SetLogLevel(1)
	bot.Bot.Debug = true
	b := bot.Bot.AddHandle()
	b.NewCommandProcessor("switchmode", switchmodehandler.SwitchModeHandler)
	b.NewCommandProcessor("changecd", changecdhandler.ChangeCdHandler)
	b.NewCommandProcessor("debug", debughandler.DebugHandler)
	b.NewProcessor(xmchecker.XmChecker, xmchecker.SendXm)
	b.Run()
}

// TODO
// 添加独立的群组模式
// 羡慕次数统计
