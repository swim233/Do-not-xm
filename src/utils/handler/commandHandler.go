package handler

import (
	"errors"
	"strconv"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/utils"
)

// 获取权限
func getPermissionArg(u tgbotapi.Update) (arg int64, success bool) {
	argStr := u.Message.CommandArguments()
	arg, err := strconv.ParseInt(argStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return arg, true
}

// 更改cd
func changeCoolDown(u tgbotapi.Update) (coolDownType string, time string, err error) {
	Args := strings.Split(u.Message.CommandArguments(), " ")
	if len(Args) != 2 {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "参数数量有误，请重新输入")
		_, err := utils.Bot.Request(msg)
		return "", "", err
	}
	if Args[1] == "static" {
		return Args[1], Args[2], nil
	} else if Args[2] == "random" {
		return Args[1], Args[2], nil
	} else {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "时间选项有误，请重新输入")
		utils.Bot.Send(msg)
		return "", "", errors.New("错误的时间选项")
	}
}
