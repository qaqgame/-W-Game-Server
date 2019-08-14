package Receive

import (
	"net"
	"fmt"
	"context"
	"wGame/Buffer"
	"wGame/Global"
	"bytes"
	"wGame/Parser"
)

func ConnHandler(conn net.Conn, timerChan chan int) {
	fmt.Println("handleConn")
	//控制子进程退出
	cxt := context.Background()
	cxt,cancle := context.WithCancel(cxt)
	defer cancle()

	//控制连接的锁
	//生成buffer队列top改变消息通知channel
	//向ReadFromBufferQueue进程发送队列目前head指针信息
	bufferchange     := make(chan *Buffer.ConnBuffer,2)
	//ReadFormBufferQueue进程向本进程发送队列目前的head指针
	bufferchangeBack := make(chan *Buffer.ConnBuffer,1)

	connbuffer  := Buffer.InitQueue()
	connbuffer2 := Buffer.InitQueue()


	bufferchange2     := make(chan *Buffer.ConnBuffer,2)
	//ReadFormBufferQueue进程向本进程发送队列目前的head指针
	bufferchangeBack2 := make(chan *Buffer.ConnBuffer,1)
	//生成独自的channel,发送从缓冲区中读取下一条信息的signal
	MyChannel    := make(chan int64, 1)
	MyChannelAck := make(chan int64, 1)
	remoteAddr   := conn.RemoteAddr().String()

	Global.Connstruct.PlayersChannel[remoteAddr]    = MyChannel
	Global.Connstruct.PlayersChannelAck[remoteAddr] = MyChannelAck
	Global.Connstruct.PlayersChannel[remoteAddr]    <- 1

	//
	var result [][]byte
	msgbuf  := bytes.NewBuffer(make([]byte,0,10240))
	databuf := make([]byte,4096)
	length  := 0
	ulength := uint32(0)

	go ReadFromBufferQueue(cxt,remoteAddr,connbuffer,bufferchange,bufferchangeBack)
	go ReadFromBufferQueue2(cxt,remoteAddr,connbuffer2,bufferchange2,bufferchangeBack2)
	for true {
		//处理粘包，并读取数据
		result = ReadFromConn(databuf,msgbuf,&length,ulength,conn)
		//重置计时器
		timerChan <- 1

		//处理从缓冲区读取的内容为nil的情况
		if result == nil {
			continue
		}

		for _,v := range result {
			//判断与处理心跳数据包，此时为所有玩家都成功连接前的状态，
			// 使用心跳数据包来保证已建立的连接的在线状态
			if string(v) == "heart beats" {
				//fmt.Println("heart beats")
				continue
			}
			//重新连接的情况处理
			//TODO

			//
			if string(v) == "conn close" {
				conn.Close()

				Global.Connstruct.RWlock.Lock()
				Global.Connstruct.ConnCount--
				delete(Global.Connstruct.Conn,remoteAddr)
				delete(Global.Connstruct.PlayersChannel,remoteAddr)
				Global.Connstruct.RWlock.Unlock()

				return
			}
			//处理正常游戏内容数据包
			//解析json数据

			//fmt.Println("value read from conn",string(v))

			//读取的内容插入缓冲区队列中
			//head := Buffer.PushIntoQueue(reqex)
			reqmini := Parser.ParserReqMini(v)
			//fmt.Println(reqmini)
			if reqmini == nil {
				continue
			}

			if reqmini.DataType == 3 {
				select {
				case data := <-bufferchangeBack2:
					connbuffer2 = data
					//temp := connbuffer.Top
					connbuffer2 = Buffer.PushIntoQueue(v,connbuffer2,reqmini.RoundNum)

					bufferchange2 <- connbuffer2
				default:
					//temp := connbuffer.Top
					connbuffer2 = Buffer.PushIntoQueue(v,connbuffer2,reqmini.RoundNum)

					bufferchange2 <- connbuffer2
				}
			} else if reqmini.DataType == 2 {
				select {
				case data := <-bufferchangeBack:
					connbuffer = data
					//temp := connbuffer.Top
					connbuffer = Buffer.PushIntoQueue(v,connbuffer,reqmini.RoundNum)

					bufferchange <- connbuffer
				default:
					//temp := connbuffer.Top
					connbuffer = Buffer.PushIntoQueue(v,connbuffer,reqmini.RoundNum)

					bufferchange <- connbuffer
				}
			}
		}
	}
}
