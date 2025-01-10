package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户的列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的chennel
	Message chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message广播信息channel的goroutine、一旦有信息就发送给全部1在线User
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		// 将msg发送给全部的在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播信息的方法
func (this *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg

}

func (this *Server) Handle(conn net.Conn) {
	// 用户上线、将用户加入到onlineMap中
	user := NewUser(conn, this)

	user.Online()
	// 监听用户是否活跃的channel
	isLive := make(chan bool)
	// 接收客户端发送的信息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}
			//提取用户的信息（去除'\n'）
			msg := string(buf[:n-1])

			// 用户针对MSG进行消息处理
			user.DoMessage(msg)

			// 用户任意消息、代表当前用户活跃
			isLive <- true
		}
	}()

	// 当前handler阻塞
	for {
		select {
		case <-isLive:
			// 当前用户是活跃的、需要重置定时器
		case <-time.After(time.Minute * 5):
			// 已经超时、将当前User强制关闭

			user.SendMsg("已退出群聊、请重新上线")

			// 销毁用户资源
			close(user.C)

			// 关闭连接
			conn.Close()

			//推出当前Handler
			return // runtime.Goexit()
		}
	}

}

// 启动服务器的接口
func (this *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()

	// 启动监听Message的goroutine
	go this.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// do handler
		go this.Handle(conn)

	}

	//close listen socket
}
