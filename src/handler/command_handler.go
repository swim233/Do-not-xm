package handler

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/bot"
	"github.com/swim233/do_not_xm/logger"
	"github.com/swim233/do_not_xm/utils"
)

func commandProcessor(processor *MessageProcessor, message string, chatID int64, msgID int64, withPermission bool) bool {
	sendMsg := func(text string) {
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyToMessageID = int(msgID)
		bot.Bot.Send(msg)
	} //发送消息

	//判断是否为修改cd命令
	if commandContact, ok := strings.CutPrefix(message, "/changecd "); ok {
		if !withPermission {
			sendMsg("你没有使用此命令的权限！")
			return false
		}
		if err := changeCdHandler(processor, commandContact, sendMsg, TimerMap[processor.GroupID]); err != nil {
			sendMsg("发生未知错误")
			logger.Suger.Warnf("Error in changing cd : %s", err.Error())
		}
		return true
	}

	//判断是否为修改模式命令
	if commandContact, ok := strings.CutPrefix(message, "/changemode "); ok {
		if !withPermission {
			sendMsg("你没有使用此命令的权限！")
			return false
		}
		changeModeHandler(processor, commandContact, sendMsg)
		return true
	}

	//判断是否为获取cd命令
	if _, ok := strings.CutPrefix(message, "/cd"); ok {
		sendMsg("当前剩余CD : " + getCdHandler(processor))
		return true
	}

	return false

}

// 修改冷却
func changeCdHandler(p *MessageProcessor, commandContact string, sender func(string), timer *utils.Timer) error {
	var singleFunc sync.Once
	var commandContacts []string
	if commandContacts = strings.Split(commandContact, " "); len(commandContacts) != 2 {
		//判定参数数量
		return errors.New("invalid number of parameters")
	}
	tempTimeMap := make(map[int]time.Duration)
	for index, toParserTime := range commandContacts {
		if times, err := parserTime(toParserTime); err != nil {
			//格式化时间
			sender("错误的时间格式")
			return err
		} else {
			//发送修改成功信息
			defer singleFunc.Do(func() {
				sender(fmt.Sprintf("CD成功修改为: %s 固定CD %s 随机CD", p.StaticCD.String(), p.RandomCD.String()))
				timer.ChangeTime(tempTimeMap[0], tempTimeMap[1])
			})
			tempTimeMap[index] = times
		}
	}
	return nil
}

// 获取当期cd
func getCdHandler(p *MessageProcessor) string {
	return TimerMap[p.GroupID].CurrentTime.String()
}

// 修改匹配模式
func changeModeHandler(p *MessageProcessor, commandContent string, sender func(s string)) bool {
	switch commandContent {
	case "match":
		p.Mode = commandContent
		sender("当前模式为 : " + commandContent)
		return true
	case "any":
		p.Mode = commandContent
		sender("当前模式为 : " + commandContent)
		return true
	default:
		sender("参数有误 请重新输入")
		return false
	}
}

// 格式化时间
func parserTime(s string) (times time.Duration, err error) {
	times, err = time.ParseDuration(s)
	if err != nil {
		return 0, err
	}
	if times < 0 {
		return 0, errors.New("can not use negative duration")
	}
	return times, nil
}
