package loveMsg

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
	instance = &loveMsg{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

var instance *loveMsg

var logger *logrus.Entry

type loveMsg struct {
}

type apiResp struct {
	Success bool   `json:"success"`
	ID      int64  `json:"id"`
	Msg     string `json:"ishan"`
}

const loveMsgApiUrl = "https://api.vvhan.com/api/love?type=json"

func (l *loveMsg) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "loveMsg"),
		Instance: instance,
	}
}

func (l *loveMsg) Init() {

}

func (l *loveMsg) PostInit() {
}

func (l *loveMsg) Serve(b *bot.Bot) {
	b.OnGroupMessage(l.sendLoveMsg)
}

func (l *loveMsg) Start(b *bot.Bot) {
}

func (l *loveMsg) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (l *loveMsg) sendLoveMsg(qqClient *client.QQClient, m *message.GroupMessage) {
	if m.ToString() != "宝贝" {
		return
	}
	msgSend := message.SendingMessage{}
	atDisplay := "@"
	if m.Sender.CardName != "" {
		atDisplay += m.Sender.CardName
	} else {
		atDisplay += m.Sender.Nickname
	}
	logger.Infof("收到 %v ：%v 的情话指令", m.Sender.Uin, atDisplay)

	msgText := getLoveMsg()
	if msgText == "" {
		msgText = "情话接口异常"
	}
	msgSend.Append(message.NewAt(m.Sender.Uin, atDisplay))
	msgSend.Append(message.NewText(msgText))
	qqClient.SendGroupMessage(m.GroupCode, &msgSend)
}

func getLoveMsg() string {
	resp := apiResp{}
	_, _, errors := gorequest.New().Timeout(time.Second * 10).Get(loveMsgApiUrl).EndStruct(&resp)
	if errors != nil {
		logger.Errorf("获取情话失败：%v", errors)
		return ""
	}
	if resp.Success {
		return resp.Msg
	} else {
		return ""
	}
}
