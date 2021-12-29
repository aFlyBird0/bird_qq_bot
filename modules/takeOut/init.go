package takeOut

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
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
	msgElements := make([]message.IMessageElement, 0)
	msgElements = append(msgElements, &message.AtElement{Target: m.Sender.Uin, SubType: message.AtTypeGroupMember})
	msgText := "宝贝，你 roll 了个 " + strconv.Itoa(utils.GetOneRandNum(1, 101))
	msgElements = append(msgElements, &message.TextElement{Content: msgText})
	msgSend := message.SendingMessage{Elements: msgElements}
	qqClient.SendGroupMessage(m.GroupCode, &msgSend)
}
