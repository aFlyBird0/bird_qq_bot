package restart

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/parnurzeal/gorequest"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

func init() {
	instance = &restart{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

var instance *restart

var logger *logrus.Entry

type restart struct {
	mConfig
}

type mConfig struct {
	allows     []int64
	webhookUrl string
	triggers   []string
}

func (r *restart) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		Instance: instance,
		ID:       bot.NewModuleID("bird", "restart"),
	}
}

func (r *restart) HotReload() {
	r.allows = bot.GetModConfigInt64Slice(r, "allows")
	r.webhookUrl = bot.GetModConfigString(r, "webhook")
	r.triggers = bot.GetModConfigStringSlice(r, "triggers")
}

func (r *restart) Init() {
	r.HotReload()
}

func (r *restart) PostInit() {
}

func (r *restart) Serve(b *bot.Bot) {
	triggerFilter := &bot.PrivateAllowMsgF{Allows: r.triggers}
	uinFilter := &bot.PrivateAllowUinF{Allows: r.allows}
	b.OnPrivateMsgAuth(r.restartByWebHook, triggerFilter, uinFilter)
}

func (r *restart) Start(b *bot.Bot) {
}

func (r *restart) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (r *restart) restartByWebHook(c *client.QQClient, m *message.PrivateMessage) {
	c.SendPrivateMessage(m.Sender.Uin, message.NewSendingMessage().Append(message.NewText("收到重启信号，准备重启")))
	logger.Info(gorequest.New().Timeout(time.Second * 10).Get(r.webhookUrl).End())
}
