package myServer

import (
	"fmt"
	"github.com/iWyh2/myTcpSystem-server/utils"
	"io"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	// 所有在线客户端
	onlineClients = make(map[string]net.Conn)
	// 读写锁
	mapLock sync.Mutex
)

// Server 服务端
type Server struct {
	Ip   string
	Port int
}

// New 获得服务器实例
func New(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}

// Run 运行服务器
func (s *Server) Run() {
	// 监听地址端口
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		ErrMsg(err)
		// 监听失败退出系统
		os.Exit(1)
	} else {
		// 打印提示
		fmt.Printf("[%s] system> starting...\n", utils.TimeStr())
		// 打印进度条
		utils.PrintProgressBar()
		// 清屏
		utils.Clear()
		// 打印banner
		utils.Banner()
	}
	defer listener.Close()
	// 打印提示
	fmt.Printf("[%s] system> running\n", utils.TimeStr())
	fmt.Printf("[%s] system> server ip: <%s> | server port: <%d>\n",
		utils.TimeStr(), s.Ip, s.Port)
	// 等待连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			// 连接失败提示
			ErrMsg(err)
		}
		// 连接成功提示
		fmt.Printf("[%s] system> client %v connects to the server successfully\n",
			utils.TimeStr(), conn.RemoteAddr().String())
		// 多协程处理连接
		go s.handler(conn)
	}
}

// handler 对客户端的连接进行业务处理
func (s *Server) handler(conn net.Conn) {
	// v0.9 群聊 + 私聊
	for msg, ok := s.receive(conn); ok; {
		// 处理消息
		s.msgHandler(msg, conn)
		msg, ok = s.receive(conn)
	}
}

// 处理消息
func (s *Server) msgHandler(msg string, conn net.Conn) {
	if msg == LoginCmd {
		// 登录指令
		s.loginHandler(conn)
	} else if msg == FindCmd {
		// 查找所有在线用户
		s.findAllHandler(conn)
	} else if len(msg) > 9 && msg[:7] == ChatToCmd {
		// 私聊
		s.privateChatHandler(conn, msg)
	} else {
		// 当前用户群发的消息 广播给所有在线用户
		s.broadcast(conn.RemoteAddr().String() + "> " + msg)
	}
}

// send 服务端发送消息
func (s *Server) send(conn net.Conn, msg string) {
	// 向连接写数据
	n, err := conn.Write([]byte(msg))
	// 写数据失败提示
	if n == 0 {
		fmt.Printf("[%s] system> send fail\n", utils.TimeStr())
	}
	// 错误提示
	if err != nil {
		ErrMsg(err)
	}
}

// receive 服务端接收消息
func (s *Server) receive(conn net.Conn) (string, bool) {
	// 4KB缓冲
	var buf = make([]byte, 1024*4)
	// 从连接读取消息
	n, err := conn.Read(buf)
	// 当客户端主动关闭连接时，conn.Read 会返回 n == 0 和 err == io.EOF
	if n == 0 {
		// 用户下线处理
		s.logoutHandler(conn)
		// 返回空消息 与 连接已关闭标志
		return "", false
	}
	// 其他错误提示
	if err != nil && err != io.EOF {
		ErrMsg(err)
		return "", true
	}
	// 返回消息 与 连接未关闭标志
	return string(buf[:n]), true
}

// 用户上线处理
func (s *Server) loginHandler(conn net.Conn) {
	mapLock.Lock()
	onlineClients[conn.RemoteAddr().String()] = conn
	mapLock.Unlock()
	// 告知其他在线客户端 该客户端已上线
	s.broadcast(fmt.Sprintf("system> user [%s] is online.",
		conn.RemoteAddr().String()))
	fmt.Printf("[%s] system> current number of clients online: %d\n",
		utils.TimeStr(), len(onlineClients))
}

// 用户下线处理
func (s *Server) logoutHandler(conn net.Conn) {
	// 该客户端下线 从在线列表删除
	mapLock.Lock()
	delete(onlineClients, conn.RemoteAddr().String())
	mapLock.Unlock()
	// 告知其他在线客户端 该客户端已下线
	s.broadcast(fmt.Sprintf("system> user [%s] is offline.",
		conn.RemoteAddr().String()))
	// 关闭连接
	conn.Close()
	// 打印提示
	fmt.Printf("[%s] system> client %v disconnects\n",
		utils.TimeStr(), conn.RemoteAddr().String())
	fmt.Printf("[%s] system> current number of clients online: %d\n",
		utils.TimeStr(), len(onlineClients))
}

// 查找所有在线用户处理
func (s *Server) findAllHandler(conn net.Conn) {
	resStr := "system> current online users:"
	mapLock.Lock()
	for name := range onlineClients {
		resStr += "\n" + "[" + name + "]"
	}
	mapLock.Unlock()
	s.send(conn, resStr)
}

// 广播消息
func (s *Server) broadcast(msg string) {
	mapLock.Lock()
	// 给每一个在线用户发送
	for _, conn := range onlineClients {
		s.send(conn, msg)
	}
	mapLock.Unlock()
}

// 私聊处理
func (s *Server) privateChatHandler(conn net.Conn, msg string) {
	// 处理消息
	strSlice := strings.Split(msg, " -")
	// 提取对方ip
	name := strSlice[1]
	// 提取消息内容
	content := strSlice[2]
	// 获得对方连接
	c := onlineClients[name]
	// 处理消息
	msg = conn.RemoteAddr().String() + "[DM You]> " + content
	// 给对方发送消息
	s.send(c, msg)
	// 给发送者显示消息
	msg = "me[DM " + name + "]> " + content
	s.send(conn, msg)
}
