package Server

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

	//在线用户列表
	OnlineMap map[string]*User
	mspLock   sync.RWMutex
	//消息广播的信息
	Message chan string
}

//创建一个server接口

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (this *Server) BroadCast(user *User, msg string) {

	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	this.Message <- sendMsg

}

//监听message广播 消息  一旦上线就发送给全部的在线user
func (this *Server) ListenMessager() {

	for {
		msg := <-this.Message
		this.mspLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mspLock.Unlock()
	}
}

//当前业务的操作
func (this *Server) Handler(conn net.Conn) {
	fmt.Println("连接建立成功")

	//创建用户
	user := NewUser(conn, this)
	user.Online()

	//监听用户是否活跃
	isLive := make(chan bool)

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4896)
		for {

			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				//this.BroadCast(user,"下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			//获取用户消息
			msg := string(buf[:n-1])

			//广播收到的消息
			//this.BroadCast(user,msg)
			user.DoMessage(msg)

			//判断用户是否是一个活跃的
			isLive <- true

		}

	}()

	//当前handle堵塞
	for {

		select {

		case <-isLive:
			//当前用户是活跃的，你该重置时钟
			//不做任何事情,额外i了通信select 更新下面定时器

		case <-time.After(time.Second * 1000):
			//已经超时了
			//将当前的user关闭
			user.SendMsg("你被踢了")

			//关闭通道
			close(user.C)

			//关闭连接
			conn.Close()
			//更新当前的hanler
			return //runtime.Goexit()
		}

	}

}

//启动服务器的接口

func (this *Server) Start() {
	//socket listen
	Listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	//accept
	if err != nil {

		fmt.Println("net Listen err", err)
	}
	defer Listen.Close()
	//do handler

	//启动监听message的消息
	go this.ListenMessager()

	for {
		//accept
		conn, err := Listen.Accept()
		if err != nil {

			fmt.Println("listen accept err:", err)
			continue
		}

		go this.Handler(conn)

	}

	//close accept

}
func add() {

	fmt.Println("22")
}
