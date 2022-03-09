package noCopy

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func init() {
	instance = &autoCopy{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

const (
	maxRepeat       = 5  // 触发自动复读次数
	msgTraceBackNum = 10 // 复读消息追溯条数，只根据前 10 条判定是否复读
)

type autoCopy struct {
	Groups      map[int64]*client.GroupInfo // 动态更新群组信息，目前未启用
	*sync.Mutex                             // 群组信息更新互斥锁
	*groupMsg
	mConfig
}

type mConfig struct {
	allowGroups []int64
}

func (m *autoCopy) HotReload() {
	m.allowGroups = bot.GetModConfigInt64Slice(m, "allowGroups")
}

var instance *autoCopy

var logger *logrus.Entry

func (m *autoCopy) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "autoCopy"),
		Instance: instance,
	}
}

func (m *autoCopy) Init() {
	m.Groups = make(map[int64]*client.GroupInfo)
	m.groupMsg = NewGroupMsg()
	m.HotReload()
}

func (m *autoCopy) PostInit() {
}

func (m *autoCopy) Serve(b *bot.Bot) {
	filters := make([]bot.OnGroupMsgFilter, 0, 3)
	filters = append(filters, &bot.GroupBanEmptyMsgF{})
	filters = append(filters, &bot.GroupAllowGroupCodeF{Allows: m.allowGroups})
	b.OnGroupMsgAuth(m.autoCopyAndJoinIn, filters...)
}

func (m *autoCopy) Start(b *bot.Bot) {
	for !b.Online.Load() {
		time.Sleep(time.Second)
	}
	time.Sleep(8 * time.Second)
}

func (m *autoCopy) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

// 判断群消息是否是复读，如果是，则加入
func (m *autoCopy) autoCopyAndJoinIn(qqClient *client.QQClient, msg *message.GroupMessage) {
	msgStr := NewGroupMessageWrapper(msg).ToString()
	// 没有复读消息，直接返回
	if !m.isMsgRepeat(msg.GroupCode, msgStr, strictCompare) {
		return
	}
	// 复读了，清空复读消息历史记录，并开始复读
	m.groupMsg.reset(msg.GroupCode)
	if msgNotSupport(msgStr) {
		return
	}
	logger.Infof("群：%v %v: %v 复读了：%v, 机器人选择加入", msg.GroupCode, msg.Sender.Uin, msg.Sender.Nickname, msg)
	msgSend := &message.SendingMessage{}
	msgSend.Append(message.NewText(msgStr))
	qqClient.SendGroupMessage(msg.GroupCode, msgSend)
}

// 有些消息QQ机器人复读麻烦，就不复读了，现在只复读纯文字
func msgNotSupport(msg string) bool {
	return strings.Contains(msg, "Image") || strings.Contains(msg, "Reply") ||
		(strings.Contains(msg, "[") && strings.Contains(msg, "]"))
}
