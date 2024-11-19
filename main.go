package main

import (
	bot "learn/goUnits/bot"
	"learn/goUnits/changecdhandler"
	"learn/goUnits/debughandler"
	"learn/goUnits/logger"
	"learn/goUnits/switchmodehandler"
	timer "learn/goUnits/timer"
	"learn/goUnits/xmchecker"
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
