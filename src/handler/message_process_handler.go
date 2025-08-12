package handler

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"github.com/swim233/do_not_xm/bot"
	"github.com/swim233/do_not_xm/config"
	"github.com/swim233/do_not_xm/logger"
	"github.com/swim233/do_not_xm/module"
	"github.com/swim233/do_not_xm/utils"
)

type Group struct {
	GroupID      int64
	RandomCD     time.Duration
	StaticCD     time.Duration
	Mode         string
	ListenUserId string
	KeyWord      string
	ReplyMessage string
}

type MessageProcessor struct {
	GroupID             int64
	RandomCD            time.Duration
	StaticCD            time.Duration
	Mode                string
	ListenUserId        ListenUserId
	KeyWord             string
	ReplyMessage        string
	IsListenUserMessage bool
}
type ListenUserId string

var MessageMap = make(map[int64]chan tgbotapi.Update)
var TimerMap = make(map[int64]*utils.Timer)
var GroupDataMu sync.Mutex
var Groups []Group
var ProcessorMap = make(map[int64]MessageProcessor)
var MessageLock sync.Locker
var Onces sync.Once

func MessageHandler(u tgbotapi.Update) error {
	Onces.Do(func() { //读取数据库内容
		if config.ConfigDB == nil {
			logger.Suger.Warn("ConfigDB is not initialized")
		}
		result := config.ConfigDB.Find(&Groups)
		if result.Error != nil {
			logger.Suger.Errorf("Fail to load groups config from database : %s", result.Error.Error())
		} else {
			for _, group := range Groups {
				ProcessorMap[group.GroupID] = MessageProcessor{
					GroupID:      group.GroupID,
					RandomCD:     group.RandomCD,
					StaticCD:     group.StaticCD,
					Mode:         group.Mode,
					ListenUserId: ListenUserId(group.ListenUserId),
					KeyWord:      group.KeyWord,
					ReplyMessage: group.ReplyMessage,
				}
				timer := &utils.Timer{

					CurrentTime:  0,
					CoolDownTime: time.Duration(ProcessorMap[group.GroupID].RandomCD) + time.Duration(ProcessorMap[group.GroupID].StaticCD),
				}
				go timer.Timer()
				MessageMap[group.GroupID] = make(chan tgbotapi.Update)
				processor := ProcessorMap[group.GroupID]
				go processor.msgProcessor(timer)
				TimerMap[group.GroupID] = timer
			}
		}
	})
	defer logger.Suger.Sync()
	var gid int64
	if u.Message.Chat.IsSuperGroup() {
		gid = u.Message.Chat.ID
	} else {
		return nil
	}
	if _, ok := ProcessorMap[gid]; ok {
		MessageMap[gid] <- u
	} else {
		ProcessorMap[gid] = MessageProcessor{
			GroupID:      gid,
			RandomCD:     300,
			StaticCD:     300,
			Mode:         "match",
			ListenUserId: getListenIDs(gid),
			KeyWord: func() string {
				str, err := json.Marshal([]string{"xm", "羡慕"})
				if err != nil {
					logger.Suger.Errorf("Error when marshaling keyword : %s", err.Error())
				}
				return string(str)
			}(),
			ReplyMessage: "不许羡慕！",
		}

		timer := &utils.Timer{

			CurrentTime:  0,
			CoolDownTime: time.Duration(ProcessorMap[gid].RandomCD) + time.Duration(ProcessorMap[gid].StaticCD),
		}
		go timer.Timer()
		TimerMap[gid] = timer
		MessageMap[gid] = make(chan tgbotapi.Update)
		processor := ProcessorMap[gid]
		go processor.msgProcessor(timer)

		MessageMap[gid] <- u

	}
	return nil

}

// 处理消息发送
func (m *MessageProcessor) msgProcessor(timer *utils.Timer) {
	for update := range MessageMap[m.GroupID] {
		if update.Message.IsCommand() {

			commandProcessor(m, update.Message.Text, update.Message.Chat.ID, int64(update.Message.MessageID), m.checkPermission(update))
			config.WriteBackDatabase(func() config.GroupDB {
				return config.GroupDB{
					GroupID:      m.GroupID,
					StaticCD:     m.StaticCD,
					RandomCD:     m.RandomCD,
					Mode:         m.Mode,
					ListenID:     string(m.ListenUserId),
					KeyWord:      m.KeyWord,
					ReplyMessage: m.ReplyMessage,
				}
			}())

		}

		if m.checkKeyWord(update) && timer.CurrentTime <= 0 {
			msgID := update.Message.MessageID
			m.sendMessage(msgID)
			timer.CurrentTime = timer.CoolDownTime
		}
	}
}

func (l ListenUserId) IDs() ([]int64, error) {
	if l == "" {
		return nil, errors.New("IDs is nil")
	} else {
		return module.ParseJson2Data[[]int64](l.string())
	}
}

func (l ListenUserId) string() string {
	return string(l)
}
func (m *MessageProcessor) RenewDataBase() {
	m.ListenUserId = getListenIDs(m.GroupID)
	config.WriteBackDatabase(func() config.GroupDB {
		return config.GroupDB{
			GroupID:      m.GroupID,
			StaticCD:     m.StaticCD,
			RandomCD:     m.RandomCD,
			Mode:         m.Mode,
			ListenID:     string(m.ListenUserId),
			KeyWord:      m.KeyWord,
			ReplyMessage: m.ReplyMessage,
		}
	}())
}

func getListenIDs(gid int64) ListenUserId {
	admins, err := bot.Bot.GetChatAdministrators(tgbotapi.ChatAdministratorsConfig{
		ChatConfig: tgbotapi.ChatConfig{
			ChatID: gid,
		},
	})
	if err != nil {
		logger.Suger.Warnf("Fail to get chat admins : %s", err.Error())
	}
	var adminsId []int64
	for _, admin := range admins {
		adminsId = append(adminsId, admin.User.ID)
	}

	strListenIDs, err := json.Marshal(adminsId)
	if err != nil {
		logger.Suger.Errorf("Error when marshaling adminsID : %s", err.Error())
	}
	return ListenUserId(strListenIDs)
}
func (m *MessageProcessor) checkPermission(u tgbotapi.Update) bool {
	return bot.Bot.IsAdmin(m.GroupID, u.Message.From.ID) || u.Message.From.ID == viper.GetInt64("owner_id")
}
