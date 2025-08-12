package handler

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/bot"
	"github.com/swim233/do_not_xm/logger"
)

func (m *MessageProcessor) sendMessage(msgID int) {
	msg := tgbotapi.NewMessage(m.GroupID, m.ReplyMessage)
	msg.ReplyToMessageID = msgID
	_, err := bot.Bot.Send(msg)
	if err != nil {
		logger.Suger.Errorf("Error in send message :", err.Error())
	}
}
