package takeOut

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

func init() {
	instance = &takeOut{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

type takeOut struct {
	mConfig
}

type mConfig struct {
	triggers []string
}

var instance *takeOut

var logger *logrus.Entry

func (t *takeOut) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "takeOut"),
		Instance: instance,
	}
}

func (t *takeOut) Init() {
	t.HotReload()
}

func (t *takeOut) HotReload() {
	t.triggers = bot.GetModConfigStringSlice(t, "triggers")
}

func (t *takeOut) PostInit() {
}

func (t *takeOut) Serve(b *bot.Bot) {
	b.OnGroupMsgAuth(t.sendRandNum, &bot.GroupAllowMsgF{Allows: t.triggers})
}

func (t *takeOut) Start(b *bot.Bot) {
}

func (t *takeOut) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

// 发送随机数
func (t *takeOut) sendRandNum(qqClient *client.QQClient, m *message.GroupMessage) {
	msgSend := message.SendingMessage{}
	msgText := "宝贝，你 roll 了个 " + strconv.Itoa(utils.GetOneRandNum(1, 101))
	atDisplay := "@"
	if m.Sender.CardName != "" {
		atDisplay += m.Sender.CardName
	} else {
		atDisplay += m.Sender.Nickname
	}
	logger.Infof("收到 %v ：%v 的外卖指令", m.Sender.Uin, atDisplay)
	msgSend.Append(message.NewAt(m.Sender.Uin, atDisplay))
	msgSend.Append(message.NewText(msgText))
	qqClient.SendGroupMessage(m.GroupCode, &msgSend)
}
