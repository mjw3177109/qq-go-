package Server

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//创建一个用户的API

func NewUser(conn net.Conn, server *Server) *User {

	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	//启动监听当前user channel的信息
	go user.ListenMessage()

	return user
}

//用户的上线服务
func (this *User) Online() {
	//当前用户上线了,将用户加入onlinemap
	this.server.mspLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mspLock.Unlock()
	//广播当前用户上线消息
	this.server.BroadCast(this, "已上线")
}

//用户的下线服务
func (this *User) Offline() {
	//当前用户上线了,将用户加入onlinemap
	this.server.mspLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mspLock.Unlock()
	this.server.BroadCast(this, "下线")
}

//给当前user可用的客户端发送消息

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

//用户处理消息的服务
func (this *User) DoMessage(msg string) {

	if msg == "who" {
		this.server.mspLock.Lock()
		for _, user := range this.server.OnlineMap {

			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mspLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//消息格式 rename|张三
		newName := strings.Split(msg, "|")[1]

		//判断name是否存在
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名已使用\n")
		} else {
			this.server.mspLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this

			this.server.mspLock.Unlock()
			this.Name = newName
			this.SendMsg("您已经更新用户名:" + this.Name + "\n")
		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]

		if remoteName == "" {
			this.SendMsg("消息格式不正确,请使用\"to|张三|你好啊\"格式.\n")
			return
		}
		//2.根据用户名 得到对方的User对象
		rendUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("该用户名不存在\n")
			return
		}
		//3.获取消息内容 通过对方的User对象将消息内容发送过去
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("无消息内容,请重发\n")
			return
		}

		rendUser.SendMsg(this.Name + "对您说" + content)

	} else {
		this.server.BroadCast(this, msg)
	}

}

//监听信息的方法 监听当前channelgo 的方法 一旦有消息就发给客户端

func (this *User) ListenMessage() {
	for {

		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
