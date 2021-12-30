package restart

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/config"
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
	webhookUrl string
}

func (r *restart) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		Instance: instance,
		ID:       bot.NewModuleID("bird", "restart"),
	}
}

func (r *restart) Init() {
	r.webhookUrl = config.GlobalConfig.GetString("ci.webhook")
}

func (r *restart) PostInit() {
}

func (r *restart) Serve(b *bot.Bot) {
	b.OnPrivateMessage(r.restartByWebHook)
}

func (r *restart) Start(b *bot.Bot) {
}

func (r *restart) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (r *restart) restartByWebHook(c *client.QQClient, m *message.PrivateMessage) {
	if m.ToString() == "#重启" {
		c.SendPrivateMessage(m.Sender.Uin, message.NewSendingMessage().Append(message.NewText("收到重启信号，准备重启")))
	}
	logger.Info(r.webhookUrl)
	logger.Info(gorequest.New().Timeout(time.Second * 10).Get(r.webhookUrl).End())
}
