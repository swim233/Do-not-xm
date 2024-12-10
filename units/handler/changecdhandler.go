package handler

import (
	"fmt"
	"learn/units/bot"
	"learn/units/counter"
	"strconv"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// 修改cd
func ChangeCdHandler(update tgbotapi.Update) error {
	if update.Message.From.ID != bot.BotConfig.IntUserID && !bot.BotConfig.DebugFlag {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你没有使用该命令的权限！")
		bot.Bot.Send(msg)
		return nil
	}
	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) != 2 {
		errMessage := fmt.Errorf("格式异常: 需要 2 个参数 但是传递了 %d 个", len(args))
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMessage.Error())
		bot.Bot.Send(msg)
		return errMessage
	}

	// 获取群组和用户id
	StaticCDstr := args[0]
	RandomCDstr := args[1]
	StaticCD, err := strconv.ParseInt(StaticCDstr, 10, 64)
	if err != nil {
		errMessage := fmt.Errorf("格式有误: \"%s\"不是有效的整形数字", StaticCDstr)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMessage.Error())
		bot.Bot.Send(msg)
		return errMessage
	}
	RandomCD, err := strconv.ParseInt(RandomCDstr, 10, 64)
	if err != nil {
		errMessage := fmt.Errorf("格式有误: \"%s\"不是有效的整形数字", RandomCDstr)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMessage.Error())
		bot.Bot.Send(msg)
		return errMessage
	}

	bot.BotConfig.RandomCD = int(RandomCD)
	bot.BotConfig.StaticCD = int(StaticCD)
	tmpStaticCD := new(int)
	tmpRandomCD := new(int)
	*tmpStaticCD = bot.BotConfig.StaticCD
	*tmpRandomCD = bot.BotConfig.RandomCD

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "修改成功，当前CD为： "+counter.Calculation(tmpStaticCD)+" 固定CD + "+counter.Calculation(tmpRandomCD)+" 随机CD")
	bot.Bot.Send(msg)
	*counter.Time = 0
	return nil
}
