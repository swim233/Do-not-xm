package main

import (
	"learn/units/bot"
	"learn/units/counter"
	"learn/units/handler"
	"learn/units/xmchecker"
)

func main() {
	bot.InitBot()
	go counter.Timer()
	bot.Bot.Debug = true
	b := bot.Bot.AddHandle()
	b.NewCommandProcessor("switchmode", handler.SwitchModeHandler)
	b.NewCommandProcessor("changecd", handler.ChangeCdHandler)
	b.NewCommandProcessor("debug", handler.DebugHandler)
	b.NewCommandProcessor("ping", handler.PingHandler)
	b.NewCommandProcessor("cd", handler.CdHandler)
	b.NewCommandProcessor("lastxm", handler.LastXmHandler)
	b.NewProcessor(xmchecker.XmChecker, xmchecker.SendXm)
	b.Run()
}

// TODO
// 添加独立的群组模式
// 羡慕次数统计
