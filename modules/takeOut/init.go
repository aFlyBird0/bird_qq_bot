package takeOut

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/config"
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
	botUin int64 //机器人自己的账号
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
	t.botUin = config.GlobalConfig.GetInt64("bot.account")
}

func (t *takeOut) PostInit() {
}

func (t *takeOut) Serve(b *bot.Bot) {
	b.OnGroupMessage(t.sendRandNum)
	b.OnGroupMessage(t.randAtMember)
}

func (t *takeOut) Start(b *bot.Bot) {
}

func (t *takeOut) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

// 发送随机数
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
	logger.Infof("收到 %v ：%v 的外卖指令", m.Sender.Uin, atDisplay)
	msgSend.Append(message.NewAt(m.Sender.Uin, atDisplay))
	msgSend.Append(message.NewText(msgText))
	qqClient.SendGroupMessage(m.GroupCode, &msgSend)
}

// 随机@ 群成员，忽略自己
func (t *takeOut) randAtMember(qqClient *client.QQClient, m *message.GroupMessage) {
	if m.ToString() != "开枪" {
		return
	}
	strconv.Itoa(utils.GetOneRandNum(1, 101))
	msgSend := message.SendingMessage{}

	group, err1 := qqClient.GetGroupInfo(m.GroupCode)
	if err1 != nil {
		logger.Error("获取群组信息失败", err1)
	}
	members, err2 := qqClient.GetGroupMembers(group)
	if err2 != nil {
		logger.Error("获取群组成员列表失败", err2)
	}
	if err1 != nil || err2 != nil {
		msgSend.Append(message.NewText("获取群组信息失败"))
		qqClient.SendGroupMessage(m.GroupCode, &msgSend)
		return
	}

	botIndex := 0 //机器人在群成员列表中的索引
	for i := range members {
		if members[i].Uin == t.botUin {
			botIndex = i
			break
		}
	}
	// 剔除机器人自己
	members = append(members[:botIndex], members[botIndex+1:]...)

	randNum := utils.GetOneRandNum(0, len(members))
	randMember := members[randNum]

	msgText := "就决定是你了！"
	atDisplay := "@"
	if randMember.CardName != "" {
		atDisplay += randMember.CardName
	} else {
		atDisplay += randMember.Nickname
	}
	msgSend.Append(message.NewAt(randMember.Uin, atDisplay))
	msgSend.Append(message.NewText(msgText))
	qqClient.SendGroupMessage(m.GroupCode, &msgSend)
}
