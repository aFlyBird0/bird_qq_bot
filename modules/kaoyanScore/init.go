package kaoyanScore

import (
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"

	"bird_qq_bot/bot"
	"bird_qq_bot/utils"
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
	tailMsg        string  // 尾部消息（可以放实验室宣传语什么的）
	webserver      webserver
	displayPicture bool   // 是否将分析结果转换为图片发到群里（和webserver可以同时开启）
	fontPath       string // 字体文件路径
}

// 注，如果 localPort 和 remoteURL 都不配置，则机器人不会向用户展示网址版数据
type webserver struct {
	localPort  string // 本地端口
	remoteURL  string // 远程webserver地址
	displayURL string // 显示在QQ机器人的消息中的webserver地址，应为配好反代后的地址
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
	m.webserver.localPort = bot.GetModConfigString(m, "webserver.localPort")
	m.webserver.remoteURL = bot.GetModConfigString(m, "webserver.remoteURL")
	m.webserver.displayURL = bot.GetModConfigString(m, "webserver.displayURL")
	m.displayPicture = bot.GetModConfigBool(m, "displayPicture")
	m.fontPath = bot.GetModConfigString(m, "fontPath")
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

	c.OnGroupMsgAuth(m.AnalyseByGroupTrigger, filters...)

	if _, err := m.cron.AddFunc("@every 10m", func() {
		m.AnalyseAndSave(c.QQClient)
	}); err != nil {
		logger.Errorf("failed to start cron, %v", err)
		return
	}
	m.cron.Start()
}

func (m *kaoyanScore) Start(_ *bot.Bot) {
	// 在本地启动一个服务器，用于展示统计结果
	if m.webserver.localPort != "" {
		RunServer(fmt.Sprintf(":%s", m.webserver.localPort), &msgFinalMap)
	}
}

func (m *kaoyanScore) Stop(_ *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
}
