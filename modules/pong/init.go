package pong

import (
	"bird_qq_bot/bot"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"sync"
	"time"
)

func init() {
	instance = &pong{}
	bot.RegisterModule(instance)
}

var instance *pong

type pong struct {
	startTime time.Time
	mConfig
}

type mConfig struct {
	triggers []string
}

func (p *pong) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		Instance: instance,
		ID:       bot.NewModuleID("bird", "pong"),
	}
}

func (p *pong) Init() {
	p.startTime = time.Now()
	p.triggers = bot.GetModConfigStringSlice(p, "triggers")
}

func (p *pong) PostInit() {
}

func (p *pong) Serve(b *bot.Bot) {
	b.OnPrivateMsgAuth(p.sendPong, &bot.PrivateAllowMsgF{Allows: p.triggers})
}

func (p *pong) Start(b *bot.Bot) {
}

func (p *pong) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (p *pong) sendPong(c *client.QQClient, m *message.PrivateMessage) {
	msg := "pong " + p.startTime.Format("2006-01-02 15:04:05")
	c.SendPrivateMessage(m.Sender.Uin, message.NewSendingMessage().Append(message.NewText(msg)))
}
