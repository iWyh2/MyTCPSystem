package myClient

import (
	"bufio"
	"fmt"
	"github.com/iWyh2/myTcpSystem-client/utils"
	"io"
	"net"
	"os"
	"strings"
)

// Client 客户端
type Client struct {
	// 服务器IP
	ServerIP string
	// 服务器端口
	ServerPort int
	// 与服务器的连接
	conn net.Conn
}

// New 获得客户端实例
func New(serverIP string, serverPort int) *Client {
	// 连接服务器
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIP, serverPort))
	if err != nil {
		// 连接失败提示
		fmt.Printf("[%s] system> connect fail\n", utils.TimeStr())
		ErrMsg(err)
		return nil
	}
	// 创建客户端实例
	client := &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		conn:       conn,
	}
	// 并返回
	return client
}

// Run 启动客户端 处理业务
func (c *Client) Run() {
	// 告知服务端上线
	c.send("cmd -login")
	// 多协程接收服务器响应
	go c.receive()
	// 输入提示
	fmt.Print("> ")
	// 创建Scanner读取标准输入
	scanner := bufio.NewScanner(os.Stdin)
	// 读取直到换行符，并将读取的文本存储在内部
	scanner.Scan()
	// 获取读取的整行文本
	msg := scanner.Text()
	// 输入exit结束客服端业务
	for msg != "exit" {
		// 向服务端发送消息
		c.send(msg)
		// 输入提示
		fmt.Print("> ")
		// 读取直到换行符
		scanner.Scan()
		// 获取读取的消息
		msg = scanner.Text()
	}
}

// 给服务器发送消息
func (c *Client) send(msg string) {
	// 向连接写数据
	n, err := c.conn.Write([]byte(msg))
	// 写数据失败提示
	if n == 0 {
		fmt.Printf("[%s] system> send fail\n", utils.TimeStr())
	}
	// 错误提示
	if err != nil {
		ErrMsg(err)
	}
}

// 接收服务器的消息
func (c *Client) receive() {
	// 4KB缓冲
	var buf = make([]byte, 1024*4)
	for {
		// 从连接读取消息
		n, err := c.conn.Read(buf)
		// 其他错误提示
		if err != nil && err != io.EOF {
			ErrMsg(err)
		}
		// 提取消息
		msg := string(buf[:n])
		// 处理消息
		msg = c.msgHandler(msg)
		// 清除当前行并打印接收到的消息
		fmt.Printf("\r%s\n", msg)
		// 输入提示
		fmt.Print("> ")
	}
}

// 处理消息
func (c *Client) msgHandler(msg string) string {
	strSlice := strings.Split(msg, "> ")
	if strSlice[0] == c.conn.LocalAddr().String() {
		return "me> " + strSlice[1]
	}
	return msg
}
