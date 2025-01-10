package main

import (
	"flag"
	"fmt"
	"net"
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

// 菜单
func (client *Client) Menu() bool {
	var (
		flag int
	)
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名模式")
	fmt.Println("0.推出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>>>>> 请输出合法范围内的数字（1，2，3，4）>>>>>>>>>>")
		return false
	}
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
			break
		case 2:
			// 私聊模式
			fmt.Println("私聊模式选择...")
			break
		case 3:
			// 更新用户名
			fmt.Println("更新用户名选择...")
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
	fmt.Println(">>>>>>>>> 连接服务器成功 >>>>>>>>>>")

	// 启动客户端的业务
	client.Run()

}
