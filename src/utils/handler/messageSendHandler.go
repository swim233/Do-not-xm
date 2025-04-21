package handler

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/utils"
)

type MessageProcessor struct {
	Handler map[int64]*XmHandler
}

var msgChan = make(chan tgbotapi.Update, 100)

// 初始化map
func NewMessageProcessor() *MessageProcessor {
	return &MessageProcessor{
		Handler: make(map[int64]*XmHandler),
	}
}

// 处理消息
func (m *MessageProcessor) Processor(u tgbotapi.Update) error {

	groupID := u.Message.Chat.ID

	xmHandler, ok := m.Handler[groupID]
	if !ok {

		xmHandler = NewXmHandler(u.Message.From.ID, u)
		m.Handler[groupID] = xmHandler
		m.SendMessage(u, *xmHandler)
		msgChan <- u
		go xmHandler.ListenMessage(msgChan)
		return nil
	} else {
		msgChan <- u
		m.SendMessage(u, *xmHandler)
		return nil
	}
}

// 发送do_not_xm消息
func (m MessageProcessor) SendMessage(u tgbotapi.Update, handler XmHandler) {
	if handler.IsXm(u) {
		msg := tgbotapi.NewMessage(u.Message.Chat.ID, "不许羡慕！")
		msg.ReplyToMessageID = u.Message.MessageID
		utils.Bot.Send(msg)
	}

}

// 更改cd
func (m MessageProcessor) ChangeCoolDown(u tgbotapi.Update) error {
	groupID := u.Message.Chat.ID

	xmHandler, ok := m.Handler[groupID]
	if !ok {

		xmHandler = NewXmHandler(u.Message.From.ID, u)
		m.Handler[groupID] = xmHandler
		if xmHandler.CheckPermission(u.Message.From.ID, u) {
			cdt, t, err := changeCoolDown(u)
			if err != nil {
				return err
			}
			intT, err := m.parseToSeconds(t)
			if err != nil {
				utils.Bot.Send(func(u tgbotapi.Update) tgbotapi.Chattable {
					return tgbotapi.NewMessage(u.Message.From.ID, "错误的时间格式")
				}(u))
				return err
			}
			if cdt == "static" {
				xmHandler.StaticCoolDown = int(intT)
			} else if cdt == "random" {
				xmHandler.RandomCoolDown = int(intT)
			}

		}
		return nil
	} else {
		if xmHandler.CheckPermission(u.Message.From.ID, u) {
			cdt, t, err := changeCoolDown(u)
			if err != nil {
				return err
			}
			intT, err := m.parseToSeconds(t)
			if err != nil {
				utils.Bot.Send(func(u tgbotapi.Update) tgbotapi.Chattable {
					return tgbotapi.NewMessage(u.Message.From.ID, "错误的时间格式")
				}(u))
				return err
			}
			if cdt == "static" {
				xmHandler.StaticCoolDown = int(intT)
			} else if cdt == "random" {
				xmHandler.RandomCoolDown = int(intT)
			}
		}
	}
	return nil
}
func (m MessageProcessor) Test(u tgbotapi.Update) error {

	groupID := u.Message.Chat.ID
	xmHandler, ok := m.Handler[groupID]
	if !ok {
		xmHandler = NewXmHandler(u.Message.From.ID, u)
		m.Handler[groupID] = xmHandler
		msg := tgbotapi.NewMessage(groupID, fmt.Sprintf("当前cd %d", xmHandler.StaticCoolDown))
		utils.Bot.Send(msg)
		return nil
	} else {
		msg := tgbotapi.NewMessage(groupID, fmt.Sprintf("当前cd %d", xmHandler.StaticCoolDown))
		utils.Bot.Send(msg)
	}
	return nil
}

// 时间解析
func (m MessageProcessor) parseToSeconds(t string) (int64, error) {

	re := regexp.MustCompile(`(\d+)d`)
	matches := re.FindAllStringSubmatch(t, -1)
	totalHours := 0
	for _, match := range matches {
		days, _ := strconv.Atoi(match[1])
		totalHours += days * 24
	}

	t = re.ReplaceAllString(t, "")

	if totalHours > 0 {
		t = strconv.Itoa(totalHours) + "h" + t
	}
	duration, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}
	return int64(duration.Seconds()), nil
}
