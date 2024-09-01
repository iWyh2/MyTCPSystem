package main

import (
	"flag"
	"fmt"
	myClient "github.com/iWyh2/myTcpSystem-client/client"
	"github.com/iWyh2/myTcpSystem-client/utils"
)

/*
TCP通讯系统 - 客户端
*/

// 连接服务器的地址端口全局变量
var (
	// 服务器地址
	serverIP string
	// 服务器端口
	serverPort int
)

// 初始化创建客户端时连接服务器的地址端口
func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "Set the server IP address")
	flag.IntVar(&serverPort, "port", 9200, "Set the server port")
}

func main() {
	// 命令行解析
	flag.Parse()
	// 创建客户端
	client := myClient.New(serverIP, serverPort)
	// 连接失败结束系统
	if client == nil {
		return
	}
	// 连接成功提示
	fmt.Printf("[%s] system> connect success\n", utils.TimeStr())
	// 打印进度条
	utils.PrintProgressBar()
	// 清屏
	utils.Clear()
	// 打印banner
	utils.Banner()
	// 启动客户端
	client.Run()
}
