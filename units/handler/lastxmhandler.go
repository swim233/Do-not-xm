package handler

import (
  "fmt"
  "learn/units/bot"
  "strings"

  tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

const (
  groupXmCapacity = 10
  otherXmCapacity = 10
)

var (
  groupXm = make([]tgbotapi.Message, 0)
  otherXm = make([]tgbotapi.Message, 0)
)

func LastXmHandler(update tgbotapi.Update) error {
  // fmtmsg := fmt.Sprintf("上一条xm为：https://t.me/ArknightsZH/%d", LastXmMessageID)
  builder := strings.Builder{}
  builder.WriteString(fmt.Sprintf("上%d条群组xm:\n", len(groupXm)))
  {
    for i := len(groupXm) - 1; i >= 0; i-- {
      item := groupXm[i]
      strMsgID := string(item.Chat.ID)
      printMsgID := strMsgID[4:]
      builder.WriteString(fmt.Sprintf("%s https://t.me/c/%s/%d\n", item.Chat.Title, printMsgID, item.MessageID))
    }
  }
  builder.WriteString(fmt.Sprintf("上%d条其他xm:\n", len(otherXm)))
  {
    for i := len(otherXm) - 1; i >= 0; i-- {
      item := otherXm[i]
      identifier := item.From.FullName()
      builder.WriteString(fmt.Sprintf("%s https://t.me/%d/%d\n", identifier, item.Chat.ID, item.MessageID))
    }
  }
  // msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmtmsg)
  msg := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
  bot.Bot.Send(msg)
  return nil
}

func RecordLastXm(update *tgbotapi.Update) {
  var targetQueue *[]tgbotapi.Message
  var targetCapacity int
  if update.Message.Chat.Type == "group" || update.Message.Chat.Type == "supergroup" {
    targetQueue = &groupXm
    targetCapacity = groupXmCapacity
  } else {
    targetQueue = &otherXm
    targetCapacity = otherXmCapacity
  }
  if len(*targetQueue) == targetCapacity {
    *targetQueue = (*targetQueue)[1:]
  }
  *targetQueue = append(*targetQueue, *update.Message)
}
