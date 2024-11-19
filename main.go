package main

import (
	"learn/units/bot"
	"learn/units/changecdhandler"
	"learn/units/debughandler"
	"learn/units/logger"
	"learn/units/switchmodehandler"
	"learn/units/timer"
	"learn/units/xmchecker"
)

func main() {
	bot.InitBot()
	go timer.Timer()
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
