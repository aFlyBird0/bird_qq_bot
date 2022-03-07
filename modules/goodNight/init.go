package goodNight

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
	"fmt"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"sync"
)

func init() {
	instance = &goodNight{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

var instance *goodNight

var logger *logrus.Entry

type goodNight struct {
	mConfig
}

type mConfig struct {
	apiKey          string
	nightTriggers   []string
	tianGouTriggers []string
}

func (g *goodNight) HotReload() {
	g.apiKey = bot.GetModConfigString(g, "apiKey")
	g.nightTriggers = bot.GetModConfigStringSlice(g, "nightTriggers")
	g.tianGouTriggers = bot.GetModConfigStringSlice(g, "tianGouTriggers")
}

func (g *goodNight) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "goodNight"),
		Instance: instance,
	}
}

func (g *goodNight) Init() {
	g.HotReload()
}

func (g *goodNight) PostInit() {
}

func (g *goodNight) Serve(b *bot.Bot) {
	b.OnGroupMsgAuth(g.sendNightMsg(), &bot.GroupAllowMsgF{Allows: g.nightTriggers})
	b.OnGroupMsgAuth(g.sendTianGouMsg(), &bot.GroupAllowMsgF{Allows: g.tianGouTriggers})
}

func (g *goodNight) Start(b *bot.Bot) {
}

func (g *goodNight) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (g *goodNight) sendNightMsg() func(qqClient *client.QQClient, m *message.GroupMessage) {
	return g.msgFromTianXing(WanAnAPI)
}

func (g *goodNight) sendTianGouMsg() func(qqClient *client.QQClient, m *message.GroupMessage) {
	return g.msgFromTianXing(TianGouAPI)
}

func (g *goodNight) msgFromTianXing(api API) func(qqClient *client.QQClient, m *message.GroupMessage) {
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
		}
		msgSend.Append(message.NewAt(m.Sender.Uin, atDisplay))
		msgSend.Append(message.NewText(msgText))
		qqClient.SendGroupMessage(m.GroupCode, &msgSend)
	}
}
