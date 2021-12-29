package takeOut

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"strconv"
	"sync"
)

func init() {
	instance = &takeOut{}
	bot.RegisterModule(instance)
}

type takeOut struct {
}

var instance *takeOut

func (t *takeOut) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "takeOut"),
		Instance: instance,
	}
}

func (t *takeOut) Init() {
}

func (t *takeOut) PostInit() {
}

func (t *takeOut) Serve(b *bot.Bot) {
	b.OnGroupMessage(t.sendRandNum)
}

func (t *takeOut) Start(b *bot.Bot) {
}

func (t *takeOut) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (t *takeOut) sendRandNum(qqClient *client.QQClient, m *message.GroupMessage) {
	if m.ToString() != "外卖" {
		return
	}
	msgSend := message.SendingMessage{}
	msgText := "宝贝，你 roll 了个 " + strconv.Itoa(utils.GetOneRandNum(1, 101))
	atDisplay := "@"
	if m.Sender.CardName != "" {
		atDisplay += m.Sender.CardName
	} else {
		atDisplay += m.Sender.Nickname
	}
	msgSend.Append(message.NewAt(m.Sender.Uin, atDisplay))
	msgSend.Append(message.NewText(msgText))
	qqClient.SendGroupMessage(m.GroupCode, &msgSend)
}
