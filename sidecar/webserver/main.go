package main

import (
	"flag"
	"fmt"

	"bird_qq_bot/utils"
)

var port int

// go run main.go -port 8090
func main() {
	// 从flag中读取端口
	flag.IntVar(&port, "port", 8090, "webserver port, default 8090")
	flag.Parse()
	utils.RunServer(fmt.Sprintf(":%d", port))
}
