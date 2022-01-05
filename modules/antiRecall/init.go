package antiRecall

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/config"
	"bird_qq_bot/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"sync"
)

func init() {
	instance = &antiRecall{}
	bot.RegisterModule(instance)
}

var instance *antiRecall

type antiRecall struct {
	mConfig
}

type mConfig struct {
	allowGroups []int64
}

func (a *antiRecall) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "antiRecall"),
		Instance: instance,
	}
}

func (a *antiRecall) InitModuleConfig() {
	allows := config.GlobalConfig.GetIntSlice("modules." + a.GetModuleInfo().ID.String() + ".allowGroups")
	a.allowGroups = make([]int64, len(allows))
	for i, v := range allows {
		a.allowGroups[i] = int64(v)
	}
}

func (a *antiRecall) Init() {
	a.InitModuleConfig()
}

func (a *antiRecall) PostInit() {
}

func (a *antiRecall) Serve(b *bot.Bot) {
	b.OnGroupMessageRecalled(a.resendRecallMsg)
}

func (a *antiRecall) Start(b *bot.Bot) {
}

func (a *antiRecall) Stop(b *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}

func (a *antiRecall) resendRecallMsg(qqClient *client.QQClient, recall *client.GroupMessageRecalledEvent) {
	if !utils.InInt64(recall.GroupCode, a.allowGroups) {
		return
	}
	// todo 开发中
	//qqClient.SendGroupMessage(recall.GroupCode, recall.)

}
