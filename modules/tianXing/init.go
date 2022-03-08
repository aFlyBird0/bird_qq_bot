package tianXing

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

func init() {
	instance = &tianXing{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

var instance *tianXing

var logger *logrus.Entry

type tianXing struct {
	mConfig
}

type mConfig struct {
	apiKey            string
	nightTriggers     []string
	tianGouTriggers   []string
	morningTriggers   []string
	healthTipTriggers []string
}

type msgHandle func(msg *string)

func (g *tianXing) HotReload() {
	g.apiKey = bot.GetModConfigString(g, "apiKey")
	g.nightTriggers = bot.GetModConfigStringSlice(g, "triggers.night")
	g.tianGouTriggers = bot.GetModConfigStringSlice(g, "triggers.dog")
	g.morningTriggers = bot.GetModConfigStringSlice(g, "triggers.morning")
	g.healthTipTriggers = bot.GetModConfigStringSlice(g, "triggers.healthTip")
}

func (g *tianXing) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "tianXing"),
		Instance: instance,
	}
}

func (g *tianXing) Init() {
	g.HotReload()
}

func (g *tianXing) PostInit() {
}

func (g *tianXing) Serve(b *bot.Bot) {
	b.OnGroupMsgAuth(g.sendNightMsg(), &bot.GroupAllowMsgF{Allows: g.nightTriggers})
	b.OnGroupMsgAuth(g.sendTianGouMsg(), &bot.GroupAllowMsgF{Allows: g.tianGouTriggers})
	b.OnGroupMsgAuth(g.sendMorningMsg(), &bot.GroupAllowMsgF{Allows: g.morningTriggers})
	b.OnGroupMsgAuth(g.sendHealthTipMsg(), &bot.GroupAllowMsgF{Allows: g.healthTipTriggers})
}

func (g *tianXing) Start(b *bot.Bot) {
}

func (g *tianXing) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (g *tianXing) sendNightMsg() func(qqClient *client.QQClient, m *message.GroupMessage) {
	return g.msgFromTianXing(wanAnAPI, addWanAnHint)
}

func (g *tianXing) sendTianGouMsg() func(qqClient *client.QQClient, m *message.GroupMessage) {
	return g.msgFromTianXing(tianGouAPI)
}

func (g *tianXing) sendMorningMsg() func(qqClient *client.QQClient, m *message.GroupMessage) {
	return g.msgFromTianXing(zaoAnAPI, addZaoAnHint)
}

func (g *tianXing) sendHealthTipMsg() func(qqClient *client.QQClient, m *message.GroupMessage) {
	return g.msgFromTianXing(healthTipAPI)
}

func (g *tianXing) msgFromTianXing(api API, handlers ...msgHandle) func(qqClient *client.QQClient, m *message.GroupMessage) {
	return func(qqClient *client.QQClient, m *message.GroupMessage) {
		msgSend := message.SendingMessage{}
		atDisplay := "@"
		if m.Sender.CardName != "" {
			atDisplay += m.Sender.CardName
		} else {
			atDisplay += m.Sender.Nickname
		}
		logger.Infof("收到 %v ：%v 的 %v 指令", m.Sender.Uin, atDisplay, api.Name)

		c := NewClient(g.apiKey)
		msgText, err := c.getFirstMsg(api)
		if err != nil {
			logger.Error(err)
			msgText = fmt.Sprintf("%v 接口请求失败", api.Name)
		} else {
			for _, h := range handlers {
				h(&msgText)
			}
		}
		msgSend.Append(message.NewAt(m.Sender.Uin, atDisplay))
		msgSend.Append(message.NewText(msgText))
		qqClient.SendGroupMessage(m.GroupCode, &msgSend)
	}
}

func addWanAnHint(msg *string) {
	if !strings.Contains(*msg, "晚安") {
		*msg += "晚安！"
	}
}

func addZaoAnHint(msg *string) {
	if !strings.Contains(*msg, "早安") {
		*msg += "早安！"
	}
}
