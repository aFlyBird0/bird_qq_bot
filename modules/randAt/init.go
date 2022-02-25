package randAt

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
	instance = &randAt{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

var instance *randAt

var logger *logrus.Entry

type randAt struct {
	botUin int64 //机器人自己的账号
	mConfig
}

type mConfig struct {
	triggers []string //触发词
}

func (r *randAt) HotReload() {
	r.triggers = bot.GetModConfigStringSlice(r, "triggers")
}

func (r *randAt) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "randAt"),
		Instance: instance,
	}
}

func (r *randAt) Init() {
	r.botUin = config.GlobalConfig.GetInt64("bot.account")
	r.HotReload()
}

func (r *randAt) PostInit() {
}

func (r *randAt) Serve(b *bot.Bot) {
	b.OnGroupMsgAuth(r.randAtMember, &bot.GroupAllowMsgF{Allows: r.triggers})
}

func (r *randAt) Start(b *bot.Bot) {
}

func (r *randAt) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	wg.Done()
}

// 随机@ 群成员，忽略自己
func (r *randAt) randAtMember(qqClient *client.QQClient, m *message.GroupMessage) {
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
		if members[i].Uin == r.botUin {
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
