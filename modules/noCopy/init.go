package noCopy

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/config"
	"bird_qq_bot/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func init() {
	instance = &noCopy{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

const (
	muteMinute      = 10 // 禁言时间
	maxRepeat       = 2  // 最大复读次数
	msgTraceBackNum = 10 // 复读消息追溯条数，只根据前 10 条判定是否复读
)

type noCopy struct {
	Groups      map[int64]*client.GroupInfo // 动态更新群组信息，目前未启用
	*sync.Mutex                             // 群组信息更新互斥锁
	*groupMsg
	*mConfig
}

type mConfig struct {
	whiteListWord []string
}

func (n *noCopy) InitModuleConfig() {
	n.mConfig = &mConfig{}
	n.whiteListWord = config.GlobalConfig.GetStringSlice("modules." + n.GetModuleInfo().String() + ".whiteListWord")
}

var instance *noCopy

var logger *logrus.Entry

func (n *noCopy) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "noCopy"),
		Instance: instance,
	}
}

func (n *noCopy) Init() {
	n.Groups = make(map[int64]*client.GroupInfo)
	n.groupMsg = NewGroupMsg()
	n.InitModuleConfig()
}

func (n *noCopy) PostInit() {
}

func (n *noCopy) Serve(b *bot.Bot) {
	//n.updateGroupMembers(b.QQClient)
	//b.OnGroupMessage(n.doNoCopyAndMute)
	b.OnGroupMessage(n.doNotCopyAndRecall)
}

func (n *noCopy) Start(b *bot.Bot) {
	for !b.Online.Load() {
		time.Sleep(time.Second)
	}
	time.Sleep(8 * time.Second)
}

func (n *noCopy) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

// 判断群消息是否是复读，如果是，则撤回
func (n *noCopy) doNotCopyAndRecall(qqClient *client.QQClient, m *message.GroupMessage) {
	msg := NewGroupMessageWrapper(m).ToString()
	if strings.TrimSpace(msg) == "" {
		return
	}
	// 若消息在白名单内，不触发复读判定
	if utils.InString(msg, n.whiteListWord) {
		return
	}
	if !n.isMsgRepeat(m.GroupCode, msg, strictCompare) {
		return
	}
	logger.Infof("群：%v %v: %v 复读了：%v", m.GroupCode, m.Sender.Uin, m.Sender.Nickname, msg)
	if err := qqClient.RecallGroupMessage(m.GroupCode, m.Id, m.InternalId); err != nil {
		logger.Info("群组消息撤回失败", err)
	}
}

// 判断群消息是否是复读，如果是，则禁言
func (n *noCopy) doNoCopyAndMute(client *client.QQClient, m *message.GroupMessage) {
	msg := NewGroupMessageWrapper(m).ToString()
	logger.Info(msg)
	if utils.InString(msg, n.whiteListWord) {
		return
	}
	if !n.isMsgRepeat(m.GroupCode, msg, strictCompare) {
		return
	}

	group, err1 := client.GetGroupInfo(m.GroupCode)
	if err1 != nil {
		logger.Error("获取用户组失败", err1)
	}
	members, err2 := client.GetGroupMembers(group)
	if err2 != nil {
		logger.Error("获取用户组成员列表失败", err2)
	}
	group.Members = members
	//group := n.getGroupFromCache(m.GroupCode)
	member := group.FindMember(m.Sender.Uin)
	if err := member.Mute(muteMinute * 60); err != nil {
		logger.Info("禁言失败", err)
	}
}
