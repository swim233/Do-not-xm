package handler

import (
	"encoding/json"
	"slices"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/logger"
)

func (m *MessageProcessor) checkKeyWord(u tgbotapi.Update) bool {
	var keywords []string
	if err := json.Unmarshal([]byte(m.KeyWord), &keywords); err != nil {
		logger.Logger.Sugar().Errorf("Error in unmarshal keyword list : %s", err.Error())
		return false
	}
	for _, kw := range keywords {
		switch m.Mode {
		case "any":
			if strings.Contains(u.Message.Text, kw) {
				return true
			}
			//匹配模式 
		case "match":
			defer m.RecordListenUserMessageLocation(u)
			if m.IsListenUserMessage {
				return strings.Contains(u.Message.Text, kw)
			} else {
				return false
			}
		}

	}
	return false
}

func (m *MessageProcessor) RecordListenUserMessageLocation(u tgbotapi.Update) {
	m.IsListenUserMessage = func() bool {
		var ids []int64
		var err error
		//获取监听用户ID
		if ids, err = m.ListenUserId.IDs(); err != nil {
			m.RenewDataBase() //如果数据为空 尝试刷新数据
			ids, _ = m.ListenUserId.IDs()
		}
		return slices.Contains(ids, u.Message.From.ID)
	}()
}
