package changecdhandler

import (
	"fmt"
	"learn/goUnits/bot"
	"learn/goUnits/timer"
	"strconv"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func ChangeCdHandler(update tgbotapi.Update) error {
	if update.Message.From.ID != bot.BotConfig.IntUserID && !bot.BotConfig.DebugFlag {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "你没有使用该命令的权限！")
		bot.Bot.Send(msg)
		return nil
	}
	args := strings.Split(update.Message.CommandArguments(), " ")
	if len(args) != 2 {
		errMessage := fmt.Errorf("格式异常: 需要两个参数 但是传递了 %d 个", len(args))
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
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("修改成功 当前cd为%ds固定cd+%ds随机cd", StaticCD, RandomCD))
	bot.Bot.Send(msg)
	*timer.Time = 0
	return nil
}
