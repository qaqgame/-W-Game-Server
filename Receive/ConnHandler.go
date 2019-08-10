package Receive

import (
	"net"
	"fmt"
	"context"
	"wGame/Buffer"
	"wGame/Global"
	"bytes"
	"wGame/Parser"
	"wGame/Model"
	"sync"
)

func ConnHandler(conn net.Conn, timerChan chan int) {
	fmt.Println("handleConn")
	cxt := context.Background()
	cxt,cancle := context.WithCancel(cxt)
	defer cancle()
	//控制连接的锁
	connMutex := sync.RWMutex{}
	//生成buffer队列top改变消息通知channel
	//向ReadFromBufferQueue进程发送队列目前head指针信息
	bufferchange     := make(chan *Buffer.Node,1)
	//ReadFormBufferQueue进程向本进程发送队列目前的head指针
	bufferchangeBack := make(chan *Buffer.Node,1)

	top,tail,size,mutex := Buffer.InitQueue()
	//生成独自的channel,发送从缓冲区中读取下一条信息的signal
	MyChannel := make(chan int, 1)
	remoteAddr := conn.RemoteAddr().String()
	Global.PlayersChannel[remoteAddr] = MyChannel
	Global.PlayersChannel[remoteAddr] <- 1

	//
	var result [][]byte
	msgbuf  := bytes.NewBuffer(make([]byte,0,10240))
	databuf := make([]byte,4096)
	length  := 0
	ulength := uint32(0)

	go ReadFromBufferQueue(cxt,remoteAddr,top,size,bufferchange,bufferchangeBack,mutex)

	for true {
		//处理粘包，并读取数据
		result = ReadFromConn(databuf,msgbuf,&length,ulength,conn)
		//重置计时器
		timerChan <- 1
		//处理从缓冲区读取的内容为nil的情况
		if result == nil {
			continue
		}
		//判断与处理心跳数据包，此时为所有玩家都成功连接前的状态，
		// 使用心跳数据包来保证已建立的连接的在线状态
		if string(result[0]) == "heart beats" {
			fmt.Println("heart beats")
			continue
		}
		//重新连接的情况处理
		//TODO

		//
		if string(result[0]) == "conn close" {
			//Global.Connstruct.RWlock.Lock()
			connMutex.Lock()
			conn.Close()
			//Global.Connstruct.ConnCount--
			//delete(Global.Connstruct.Conn,remoteAddr)
			//delete(Global.Connstruct.PlayersChannel,remoteAddr)
			Global.ConnCount--
			delete(Global.Conn,remoteAddr)
			delete(Global.PlayersChannel,remoteAddr)
			//Global.Connstruct.RWlock.Unlock()
			connMutex.Unlock()
			return
		}
		//处理正常游戏内容数据包
		//解析json数据
		for _,v := range result {
			//fmt.Println("value",v)
			req := Parser.ParserReq(v)
			if req == nil {
				break
			}
			//重新组装新格式
			var reqex Model.ReqEx
			reqex.Request    = *req
			reqex.UserId     = req.UserID
			reqex.RemoteAddr = remoteAddr

			//读取的内容插入缓冲区队列中
			//head := Buffer.PushIntoQueue(reqex)
			select {
			case data := <-bufferchangeBack:
				top = data
				temp := Buffer.PushIntoQueue(reqex,top,tail,size,mutex)
				if temp[0] != top {
					bufferchange <- temp[0]
				}
				tail = temp[1]
			default:
				temp := Buffer.PushIntoQueue(reqex,top,tail,size,mutex)
				if temp[0] != top {
					bufferchange <- temp[0]
				}
				tail = temp[1]
			}
		}
	}
}
