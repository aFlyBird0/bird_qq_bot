package kaoyanScore

import (
	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

// 考研分数统计

func init() {
	instance = &kaoyanScore{}
	bot.RegisterModule(instance)
	logger = utils.GetModuleLogger(instance.GetModuleInfo().String())
}

var instance *kaoyanScore

var logger *logrus.Entry

var msgFinalMap sync.Map

type kaoyanScore struct {
	mConfig
	lastUpdateTime time.Time
	cron           *cron.Cron //
}

type mConfig struct {
	triggers       []string
	AllowGroupList []int64 // 开启的群号列表
	adminList      []int64 // 管理员列表, 目前有问题，所有群的管理混在一起了
	tailMsg        string  // 尾部消息（实验室宣传语）
}

func (m *kaoyanScore) GetModuleInfo() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       bot.NewModuleID("bird", "kaoyanScore"),
		Instance: instance,
	}
}

func (m *kaoyanScore) Init() {
	m.HotReload()
	m.cron = cron.New() // 就不设置定时任务了，直接用消息触发吧
}

func (m *kaoyanScore) HotReload() {
	m.AllowGroupList = bot.GetModConfigInt64Slice(m, "allowGroupList")
	m.triggers = bot.GetModConfigStringSlice(m, "triggers")
	m.tailMsg = bot.GetModConfigString(m, "tailMsg")
}

func (m *kaoyanScore) PostInit() {
}

func (m *kaoyanScore) Serve(c *bot.Bot) {
	time.AfterFunc(time.Second*9, func() {
		m.updateAdminList(c.QQClient)
	})
	//time.AfterFunc(time.Second*10, func() {
	//	// 找到所有计算机学硕学生，并拼接邮件群发的收件人字段
	//	mailsTo := GenMailsFromUins(m.FindCSAcademicStudent(c.QQClient))
	//	fmt.Println(mailsTo)
	//})
	//time.AfterFunc(time.Second*10, func() {
	//	m.FindInvalidCardName(c.QQClient)
	//})
	filters := make([]bot.OnGroupMsgFilter, 0, 3)
	filters = append(filters, &bot.GroupAllowMsgF{Allows: m.triggers})
	filters = append(filters, &bot.GroupAllowGroupCodeF{Allows: m.AllowGroupList})
	//filters = append(filters, &bot.GroupAllowUinF{Allows: m.adminList})
	c.OnGroupMsgAuth(m.CalculateByGroupTrigger, filters...)

}

func (m *kaoyanScore) Start(c *bot.Bot) {
	runGin()
}

func (m *kaoyanScore) Stop(c *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}
