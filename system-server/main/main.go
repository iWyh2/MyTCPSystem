package main

import myServer "github.com/iWyh2/myTcpSystem-server/server"

/*
TCP通讯系统
作者：iWyh2
功能描述：
1.群聊
2.私聊
*/

/*
TCP通讯系统 - 服务端
*/

func main() {
	server := myServer.New("127.0.0.1", 8888)
	server.Run()
}
