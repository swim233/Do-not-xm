package handler

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/utils"
)

type MessageProcessor struct {
	Handler map[int64]*XmHandler
}

// 初始化map
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		Handler: make(map[int64]*XmHandler),
	}
}

// 处理消息
func (m *MessageProcessor) Processor(u tgbotapi.Update) error {
	xmHandler, exists := m.getXmHandler(u)
	m.SendMessage(u, xmHandler)
	xmHandler.UpdateChannels <- u
	if !exists {
		go xmHandler.ListenMessage(xmHandler.UpdateChannels)
		go xmHandler.timer(xmHandler)
	}
	return nil
}

// 发送do_not_xm消息
func (m *MessageProcessor) SendMessage(u tgbotapi.Update, handler *XmHandler) {
	if handler.IsXm(u) && handler.CoolDownTime == 0 {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "不许羡慕！")
		msg.ReplyToMessageID = u.Message.MessageID
		utils.Bot.Send(msg)
		if handler.RandomCoolDown > 0 {
			handler.CoolDownTime = handler.StaticCoolDown + rand.Intn(handler.RandomCoolDown+1)
		} else {
			handler.CoolDownTime = handler.StaticCoolDown
		}
	}
}

// 更改cd
func (m *MessageProcessor) ChangeCoolDown(u tgbotapi.Update) error {
	xmHandler, _ := m.getXmHandler(u)
	if !xmHandler.CheckPermission(u.Message.From.ID, u) { //检查权限
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "你没有执行此操作的权限！")
		utils.Bot.Send(msg)
		return nil
	} else {
		coolDownType, time, err := changeCoolDown(u) //获取冷却时间和类型
		if err != nil || time == "" {
			return err
		}

		intT, err := m.parseToSeconds(time) //格式化时间
		if err != nil {
			utils.Bot.Send(func(u tgbotapi.Update) tgbotapi.Chattable {
				return tgbotapi.NewMessage(u.Message.Chat.ID, "错误的时间格式")
			}(u))
			return err
		}
		if coolDownType == "static" {
			xmHandler.StaticCoolDown = int(intT)
		} else if coolDownType == "random" {
			xmHandler.RandomCoolDown = int(intT)
		}
		msgStr := fmt.Sprintf("当前CD为 %s 固定CD %s 随机CD", formatSeconds(xmHandler.StaticCoolDown), formatSeconds(xmHandler.RandomCoolDown))
		utils.Bot.Send(func(u tgbotapi.Update) tgbotapi.Chattable {
			return tgbotapi.NewMessage(u.Message.Chat.ID, msgStr)
		}(u))

		if xmHandler.RandomCoolDown > 0 {
			xmHandler.CoolDownTime = xmHandler.StaticCoolDown + rand.Intn(xmHandler.RandomCoolDown+1)
		} else {
			xmHandler.CoolDownTime = xmHandler.StaticCoolDown
		}

	}
	return nil
}

// 输出剩余cd
func (m *MessageProcessor) CD(u tgbotapi.Update) error {
	xmHandler, _ := m.getXmHandler(u)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, fmt.Sprintf("当前剩余CD %s",
		formatSeconds(xmHandler.CoolDownTime)))
	utils.Bot.Send(msg)
	return nil
}

// 修改匹配模式
func (m *MessageProcessor) SwitchTrigger(u tgbotapi.Update) error {
	XmHandler, _ := m.getXmHandler(u)
	triggerMode, err := switchTrigger(u)
	if err != nil {
		return err
	}
	XmHandler.TriggerMode = triggerMode
	return nil
}

// 时间解析
func (m *MessageProcessor) parseToSeconds(t string) (int64, error) {
	re := regexp.MustCompile(`(\d+)(d|h|m|s)`)
	matches := re.FindAllString(t, -1)
	if len(matches) == 0 {
		return 0, fmt.Errorf("invalid time format: %q", t)
	}

	// 验证输入是否完全由合法部分组成
	if strings.Join(matches, "") != t {
		return 0, fmt.Errorf("invalid time format: %q", t)
	}

	var totalSeconds int64
	for _, part := range matches {
		// 分离数值和单位
		var valueStr, unit string
		for i, c := range part {
			if c < '0' || c > '9' {
				valueStr = part[:i]
				unit = part[i:]
				break
			}
		}

		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid value in %q: %v", part, err)
		}

		switch unit {
		case "d":
			totalSeconds += value * 86400 // 24*60*60
		case "h":
			totalSeconds += value * 3600 // 60*60
		case "m":
			totalSeconds += value * 60
		case "s":
			totalSeconds += value
		default:
			return 0, fmt.Errorf("unknown unit %q in %q", unit, part)
		}
	}
	return totalSeconds, nil
}

// 实例化xmHandler
func (m *MessageProcessor) getXmHandler(u tgbotapi.Update) (handler *XmHandler, exists bool) {

	groupID := u.Message.Chat.ID
	xmHandler, ok := m.Handler[groupID]
	if !ok {
		xmHandler = NewXmHandler(u.Message.Chat.ID, u)
		m.Handler[groupID] = xmHandler
		return xmHandler, false
	} else {
		return xmHandler, true
	}
}

// 格式化时间
func formatSeconds(seconds int) string {
	days := seconds / (24 * 3600)
	seconds %= 24 * 3600
	hours := seconds / 3600
	seconds %= 3600
	minutes := seconds / 60
	seconds %= 60

	return fmt.Sprintf("%d天%d小时%d分钟%d秒", days, hours, minutes, seconds)
}
