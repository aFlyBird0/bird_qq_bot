package goodNight

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
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
	apiKey   string
	triggers []string
}

func (g *goodNight) HotReload() {
	g.apiKey = bot.GetModConfigString(g, "apiKey")
	g.triggers = bot.GetModConfigStringSlice(g, "triggers")
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
	b.OnGroupMsgAuth(g.sendNightMsg, &bot.GroupAllowMsgF{Allows: g.triggers})
}

func (g *goodNight) Start(b *bot.Bot) {
}

func (g *goodNight) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (g *goodNight) sendNightMsg(qqClient *client.QQClient, m *message.GroupMessage) {
	msgSend := message.SendingMessage{}
	atDisplay := "@"
	if m.Sender.CardName != "" {
		atDisplay += m.Sender.CardName
	} else {
		atDisplay += m.Sender.Nickname
	}
	logger.Infof("收到 %v ：%v 的晚安指令", m.Sender.Uin, atDisplay)

	msgText := getNightMsg(g.apiKey)
	if msgText == "" {
		msgText = "晚安消息接口异常"
	}
	msgSend.Append(message.NewAt(m.Sender.Uin, atDisplay))
	msgSend.Append(message.NewText(msgText))
	qqClient.SendGroupMessage(m.GroupCode, &msgSend)
}
