package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	flag       int // 当前client模式
}

func NewClient(serverIp string, serverPort int) *Client {
	// 创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       99,
	}
	// 连接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.Conn = conn
	// 返回对象
	return client
}

// 处理server回应的消息、直接显示到标准输出
func (client *Client) DealResponse() {
	// 一旦client.conn有数据、就直接copy到student标准输出上、永久阻塞监听
	_, _ = io.Copy(os.Stdout, client.Conn)
	//

}

// 菜单
func (client *Client) Menu() bool {
	var (
		flag int
	)
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名模式")
	fmt.Println("0.推出")

	_, _ = fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>>>> 请输出合法范围内的数字（1，2，3，4）>>>>>>>>>>")
		return false
	}
}

// 查询用户并显示
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

// 私聊模式
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string
	client.SelectUsers()
	fmt.Println(">>>>>> 请输入聊天对象p[用户名],exit退出:")
	_, _ = fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>>请输入消息内容,exit退出:")
		_, _ = fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			// 发送给服务器、消息不为空继续发送
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.Conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("client.Conn.Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>>>请输入消息内容,exit退出:")
			_, _ = fmt.Scanln(&chatMsg)
		}

		fmt.Println(">>>>>> 请输入聊天对象p[用户名],exit退出:")
		_, _ = fmt.Scanln(&remoteName)
	}

}

// 公聊模式
func (client *Client) PublicChat() {
	// 提示用户输入信息
	var chatMsg string
	fmt.Println(">>>>> 请输入聊天内容、exit推出：")
	_, _ = fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		// 发送给服务器、消息不为空继续发送
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("client.Conn.Write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>> 请输入聊天内容、exit推出：")
		fmt.Println(&chatMsg)

	}
}

// 更新用户名
func (client *Client) UpdateName() bool {
	fmt.Println(">>>>>>>> 请输入用户名:")
	_, _ = fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}
	return true
}

// 客户端run方法
func (client *Client) Run() {
	for client.flag != 0 {
		for client.Menu() != true {
		}
		// 根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			// 公聊模式
			fmt.Println("公聊模式选择...")
			client.PublicChat()
			break
		case 2:
			// 私聊模式
			fmt.Println("私聊模式选择...")
			client.PrivateChat()
			break
		case 3:
			// 更新用户名
			client.UpdateName()
			break
		}
	}

}

var (
	serverIp   string
	serverPort int
)

// ./client -ip 127.0.01 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "Port", 8888, "设置服务器端口（默认8888）")
}

func main() {
	// 命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>> 连接服务器失败 >>>>>>>")
		return
	}
	// 单独开启一个foroutine处理server的回执消息
	go client.DealResponse()

	fmt.Println(">>>>>>>>> 连接服务器成功 >>>>>>>>>>")

	// 启动客户端的业务
	client.Run()

}
