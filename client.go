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
	conn       net.Conn
	flag       int
}

func NewClient(ServerIp string, ServerPort int) *Client {
	//创建客户端对象
	client := &Client{

		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))
	if err != nil {
		fmt.Println("net dial error", err)
	}
	client.conn = conn
	return client
}

//处理server回复的消息
func (client *Client) DealResponse() {
	//一旦conn有数据 就会copy到stdout标准输出上 永久堵塞
	io.Copy(os.Stdout, client.conn)

	//

}

//查询在线用户
func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return
	}

}

//私聊模式

func (client *Client) PrivateChat() {
	var selectName string
	var chatMsg string
	client.SelectUsers()
	fmt.Println(">>>请输入你要聊天的用户名,exit退出>>>")
	fmt.Scanln(&selectName)

	for selectName != "exit" {

		fmt.Println(">>>>请输入聊天内容,exit退出>>>>>")

		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {

			if len(chatMsg) != 0 {

				sendMsg := "to|" + selectName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println(">>>>请输入聊天内容,exit退出")
			fmt.Scanln(&chatMsg)

		}

		client.SelectUsers()
		fmt.Println(">>>请输入你要聊天的用户名,exit退出>>>")
		fmt.Scanln(&selectName)

	}

}

//更改用户名
func (client *Client) UpdateName() bool {

	fmt.Println(">>>>请输入用户名>>>>")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("CONN.Write err:", err)
		return false
	}

	return true

}

//客户公聊模式
func (client *Client) PublicChat() {
	//提示用户输入信息
	var chatMsg string

	fmt.Println(">>>>请输入关键内容,exit退出.")

	fmt.Scanln(&chatMsg)
	//发送服务器

	//消息不为空则发送
	for chatMsg != "exit" {
		if len(chatMsg) != 0 {

			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println(">>>>请输入聊天内容,exit退出")
		fmt.Scanln(&chatMsg)

	}

}

//客户端菜单
func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	fmt.Println(flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true

	} else {
		fmt.Println(">>>>>请输入合法范围内的值>>>>>>")
		return false
	}

}

//主业务
func (client *Client) Run() {
	for client.flag != 0 {

		for client.menu() != true {

		}

		//根据不同的模式处理不同的业务
		switch client.flag {
		case 1:
			//公聊模式
			fmt.Println("公聊模式选择")
			client.PublicChat()
			break

		case 2:
			//私聊模式
			fmt.Println("私聊模式选择")
			client.PrivateChat()
			break

		case 3:
			//更新用户名
			fmt.Println("更新用户名")
			client.UpdateName()
			break

		}
	}
}

var serverip string
var serverport int

func init() {
	flag.StringVar(&serverip, "ip", "127.0.0.1", "设置服务器IP链接(默认是127.0.0.1)")
	flag.IntVar(&serverport, "port", 9999, "设置服务器port端口(默认是9999)")
}

func main() {
	flag.Parse()
	client := NewClient("127.0.0.1", 9999)
	if client == nil {
		fmt.Println(">>>>>>>链接失败>>>>>>")
		return
	}

	//单独开启一个goruntime处理接受到的消息
	go client.DealResponse()
	fmt.Println(">>>>链接服务器成功>>>>>")

	//启动客户端的业务

	//select{
	//
	//}
	client.Run()

}
