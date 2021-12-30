package main

import (
	"bird_qq_bot/bot"
	"fmt"
	"os"
	"os/signal"

	"bird_qq_bot/config"
	"bird_qq_bot/utils"

	_ "bird_qq_bot/modules/logging"
	_ "bird_qq_bot/modules/loveMsg"
	_ "bird_qq_bot/modules/noCopy"
	_ "bird_qq_bot/modules/pong"
	_ "bird_qq_bot/modules/restart"
	_ "bird_qq_bot/modules/takeOut"
)

func init() {
	utils.WriteLogToFS()
	config.Init()
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
	if err := bot.Login(); err != nil {
		fmt.Println("登录失败:", err)
	}

	// 刷新好友列表，群列表
	bot.RefreshList()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
	bot.Stop()
}
