package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"

	"bird_qq_bot/bot"
	"bird_qq_bot/utils"

	_ "bird_qq_bot/modules/autoCopy"
	_ "bird_qq_bot/modules/kaoyanScore"
	_ "bird_qq_bot/modules/logging"
	_ "bird_qq_bot/modules/loveMsg"
	_ "bird_qq_bot/modules/noCopy"
	_ "bird_qq_bot/modules/pong"
	_ "bird_qq_bot/modules/randAt"
	_ "bird_qq_bot/modules/restart"
	_ "bird_qq_bot/modules/takeOut"
	_ "bird_qq_bot/modules/tianXing"
)

func init() {
	utils.WriteLogToFS()
	//config.Init()
}

func main() {
	// 快速初始化
	bot.Init()

	// 初始化 Modules
	bot.StartService()

	// 使用协议
	// 不同协议可能会有部分功能无法使用
	// 在登陆前切换协议
	bot.UseProtocol(bot.IPad)

	// 登录
	err := bot.Login()
	if err != nil {
		panic(err)
	}
	//bot.SaveToken() // 存储快速登录使用的 Token, 如需使用快捷登录请解除本条注释

	// 刷新好友列表，群列表
	bot.RefreshList()

	bot.Instance.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		fmt.Printf("收到群聊信息：%v", event)
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	bot.Stop()
}
