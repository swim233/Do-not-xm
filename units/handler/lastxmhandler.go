package handler

import (
	"fmt"
	"learn/units/bot"
	"learn/units/collections"
	"strings"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

const (
	groupXmCapacity = 10
	otherXmCapacity = 10
)

var (
	groupXm = collections.NewOverflowQueue[*tgbotapi.Update](groupXmCapacity)
	otherXm = collections.NewOverflowQueue[*tgbotapi.Update](otherXmCapacity)
)

func LastXmHandler(update tgbotapi.Update) error {
	// fmtmsg := fmt.Sprintf("上一条xm为：https://t.me/ArknightsZH/%d", LastXmMessageID)
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("上%d条群组xm:\n", groupXm.Size()))
	{
		iterator := groupXm.NewReverseIterator()
		for item, ok := iterator.Next(); ok; item, ok = iterator.Next() {
			strMsgID := string(item.Message.Chat.ID)
			printMsgID := strMsgID[4:]
			builder.WriteString(fmt.Sprintf("%s https://t.me/c/%s/%d\n", item.Message.Chat.Title, printMsgID, item.Message.MessageID))
		}
	}
	builder.WriteString(fmt.Sprintf("上%d条其他xm:\n", otherXm.Size()))
	{
		iterator := otherXm.NewReverseIterator()
		for item, ok := iterator.Next(); ok; item, ok = iterator.Next() {
			identifier := item.Message.From.FullName()
			builder.WriteString(fmt.Sprintf("%s https://t.me/%d/%d\n", identifier, item.Message.Chat.ID, item.Message.MessageID))
		}
	}
	// msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmtmsg)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
	bot.Bot.Send(msg)
	return nil
}

func RecordLastXm(update *tgbotapi.Update) {
	if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
		groupXm.Enqueue(update)
	} else {
		otherXm.Enqueue(update)
	}
}
