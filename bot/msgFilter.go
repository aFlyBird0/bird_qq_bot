package bot

import (
	"bird_qq_bot/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"strings"
)

type OnPrivateMsgFunc func(*client.QQClient, *message.PrivateMessage)
type OnGroupMsgFunc func(*client.QQClient, *message.GroupMessage)

type OnPrivateMsgFilter interface {
	CanAccess(*message.PrivateMessage) bool
}

type OnGroupMsgFilter interface {
	CanAccess(*message.GroupMessage) bool
}

func (b *Bot) OnPrivateMsgAuth(f OnPrivateMsgFunc, filters ...OnPrivateMsgFilter) {
	b.OnPrivateMessage(func(qqClient *client.QQClient, privateMessage *message.PrivateMessage) {
		for _, filter := range filters {
			if !filter.CanAccess(privateMessage) {
				return
			}
		}
		f(qqClient, privateMessage)
	})
}

func (b *Bot) OnGroupMsgAuth(f OnGroupMsgFunc, filters ...OnGroupMsgFilter) {
	b.OnGroupMessage(func(qqClient *client.QQClient, groupMessage *message.GroupMessage) {
		for _, filter := range filters {
			if !filter.CanAccess(groupMessage) {
				return
			}
		}
		f(qqClient, groupMessage)
	})
}

type PrivateAllowUinF struct {
	Allows []int64
}

func (p *PrivateAllowUinF) CanAccess(privateMessage *message.PrivateMessage) bool {
	return utils.InInt64(privateMessage.Sender.Uin, p.Allows)
}

type PrivateBanUinF struct {
	Bans []int64
}

func (p *PrivateBanUinF) CanAccess(privateMessage *message.PrivateMessage) bool {
	return !utils.InInt64(privateMessage.Sender.Uin, p.Bans)
}

type PrivateAllowMsgF struct {
	Allows []string
}

func (p *PrivateAllowMsgF) CanAccess(privateMessage *message.PrivateMessage) bool {
	return utils.InString(privateMessage.ToString(), p.Allows)
}

type PrivateBanMsgF struct {
	Bans []string
}

func (p *PrivateBanMsgF) CanAccess(privateMessage *message.PrivateMessage) bool {
	return !utils.InString(privateMessage.ToString(), p.Bans)
}

type PrivateBanEmptyMsgF struct {
}

func (p *PrivateBanEmptyMsgF) CanAccess(privateMessage *message.PrivateMessage) bool {
	return strings.TrimSpace(privateMessage.ToString()) != ""
}

type GroupAllowGroupCodeF struct {
	Allows []int64
}

func (g *GroupAllowGroupCodeF) CanAccess(groupMessage *message.GroupMessage) bool {
	return utils.InInt64(groupMessage.Sender.Uin, g.Allows)
}

type GroupBanUinF struct {
	Bans []int64
}

func (g *GroupBanUinF) CanAccess(groupMessage *message.GroupMessage) bool {
	return !utils.InInt64(groupMessage.Sender.Uin, g.Bans)
}

type GroupAllowMsgF struct {
	Allows []string
}

func (g *GroupAllowMsgF) CanAccess(groupMessage *message.GroupMessage) bool {
	return utils.InString(groupMessage.ToString(), g.Allows)
}

type GroupBanMsgF struct {
	Bans []string
}

func (g *GroupBanMsgF) CanAccess(groupMessage *message.GroupMessage) bool {
	return !utils.InString(groupMessage.ToString(), g.Bans)
}

type GroupBanEmptyMsgF struct {
}

func (g *GroupBanEmptyMsgF) CanAccess(groupMessage *message.GroupMessage) bool {
	return strings.TrimSpace(groupMessage.ToString()) != ""
}
