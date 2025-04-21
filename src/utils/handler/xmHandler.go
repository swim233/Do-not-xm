package handler

import (
	"regexp"
	"strings"
	"time"

	"slices"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/swim233/do_not_xm/utils"
)

type XmHandler struct {
	WaitingForMessageID int //监听用户消息id
	WaitingForGroupID   int //监听群组id
	UpdateChannels      chan tgbotapi.Update
	HandlerConfig       //消息处理配置
	Timer               //计时器
}
type HandlerConfig struct {
	EnableChecker    bool    //是否开启羡慕检查
	RandomCoolDown   int     //随机冷却时间
	StaticCoolDown   int     //固定冷却时间
	CoolDownTime     int     //当前冷却时间
	TriggerMode      string  //触发模式
	WaitingForUserID int64   //监听用户id
	PermissionUserID []int64 //权限用户id
}

type Timer struct {
}

var re = regexp.MustCompile(".*羡.*慕.*")

// 构造函数
func NewXmHandler(waitingForGroupID int64, u tgbotapi.Update) *XmHandler {

	newHandler := &XmHandler{
		WaitingForMessageID: 0,
		WaitingForGroupID:   int(waitingForGroupID),
		HandlerConfig: HandlerConfig{
			WaitingForUserID: utils.BotConfig.WaitingForUserID,
			PermissionUserID: []int64{utils.BotConfig.WaitingForUserID},
			EnableChecker:    utils.BotConfig.EnableChecker,
			TriggerMode:      utils.BotConfig.TriggerMode,
		},
	}
	return newHandler
}

// 检查是否符合回复条件
func (x *XmHandler) IsXm(u tgbotapi.Update) bool {

	match := re.MatchString(u.Message.Text) || strings.Contains(u.Message.Text, "xm")
	next := func(u tgbotapi.Update) bool {
		return u.Message.MessageID == x.WaitingForMessageID+1
	}(u)
	if x.TriggerMode == "match" {
		return match && next
	} else {
		return match
	}
}

// 获取监听用户消息ID
func (x *XmHandler) ListenMessage(ch <-chan tgbotapi.Update) {
	for u := range ch {
		if u.Message.From.ID == x.WaitingForUserID {
			x.WaitingForMessageID = u.Message.MessageID
		}
	}
}

// 添加权限用户 默认监听用户拥有权限
func (x *XmHandler) AddPermissionUser(u tgbotapi.Update) string {
	if x.CheckPermission(u.Message.From.ID, u) {
		id, success := getPermissionArg(u)
		if success {
			x.PermissionUserID = append(x.PermissionUserID, id)
			return "success"
		} else {
			return "false"
		}
	}
	return "without_permission"
}

// 删除权限用户
func (x *XmHandler) RemovePermissionUser(u tgbotapi.Update) string {
	if x.CheckPermission(u.Message.From.ID, u) {
		id, success := getPermissionArg(u)
		if success {
			for i, v := range x.PermissionUserID {
				if v == id {
					x.PermissionUserID = slices.Delete(x.PermissionUserID, i, i+1)
					return "success"
				}
			}
			return "not_found"
		} else {
			return "false"
		}
	}
	return "without_permission"
}

// 检查用户权限
func (x *XmHandler) CheckPermission(userID int64, u tgbotapi.Update) bool {
	if utils.Bot.IsAdminWithPermissions(u.Message.Chat.ID, userID, tgbotapi.AdminCanChangeInfo) ||
		slices.Contains(x.PermissionUserID, userID) {
		return true
	}
	return false
}

// 计时器
func (t *Timer) timer(handler *XmHandler) {
	for {
		time.Sleep(time.Second)
		if handler.CoolDownTime > 0 {
			handler.CoolDownTime--
		}
	}
}
