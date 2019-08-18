package Receive

import (
	"net"
	"fmt"
	"context"
	"wGame/Buffer"
	"wGame/Global"
	"bytes"
	"wGame/Parser"
	"sync"
)

func ConnHandler(conn net.Conn, timerChan chan int) {
	//var t string = ""
	//t := 0
	oncefunc := sync.Once{}
	fmt.Println("handleConn")
	//控制子进程退出
	cxt := context.Background()
	cxt,cancle := context.WithCancel(cxt)
	cxt1 := context.Background()
	cxt1,cancle1 := context.WithCancel(cxt1)
	//defer cancle()
	//defer cancle1()

	//控制连接的锁
	//生成buffer队列top改变消息通知channel
	//向ReadFromBufferQueue进程发送队列目前head指针信息
	//bufferchange     := make(chan *Buffer.ConnBuffer,2)
	//ReadFormBufferQueue进程向本进程发送队列目前的head指针
	//bufferchangeBack := make(chan *Buffer.ConnBuffer,1)

	connbuffer  := Buffer.InitQueue()
	connbuffer2 := Buffer.InitQueue()


	//bufferchange2     := make(chan *Buffer.ConnBuffer,2)
	//ReadFormBufferQueue进程向本进程发送队列目前的head指针
	//bufferchangeBack2 := make(chan *Buffer.ConnBuffer,1)
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

	go ReadFromBufferQueue(cxt,remoteAddr,connbuffer)
	go ReadFromBufferQueue2(cxt1,remoteAddr,connbuffer2)
	for true {
		//处理粘包，并读取数据
		result = ReadFromConn(databuf,msgbuf,&length,ulength,conn)
		//fmt.Println("msgbuf len:",msgbuf.Len())
		//fmt.Println("length:",len(result))
		//for _,v := range result {
		//	fmt.Println("V",string(v))
		//}
		//重置计时器
		timerChan <- 1

		//处理从缓冲区读取的内容为nil的情况
		if result == nil {
			continue
		}
		//fmt.Println("receive",time.Now().Format(time.RFC3339Nano))
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
			if string(v) == "Ready" {
				fmt.Println("ReadyGame")
				//开始游戏，打开转发和转发计时器进程
				oncefunc.Do(StartGmae)
				continue
			}
			//
			if string(v) == "conn close" {
				conn.Close()

				Global.Connstruct.RWlock.Lock()
				Global.Connstruct.ConnCount--
				delete(Global.Connstruct.Conn,remoteAddr)
				//delete(Global.Connstruct.PlayersChannel,remoteAddr)
				Global.Connstruct.RWlock.Unlock()
				//fmt.Println("close !!!")
				cancle()
				cancle1()
				return
			}
			//处理正常游戏内容数据包
			//解析json数据

			//if t != string(v) {
			//	t = string(v)
			//	fmt.Println("value read from conn",string(v))
			//}
			//fmt.Println("value read from conn",string(v))
			//fmt.Println("value read from conn",string(v))
			//t = t + "\n value read from conn" +string(v)
			//fmt.Println("test",string(v[4:]))

			//读取的内容插入缓冲区队列中
			//head := Buffer.PushIntoQueue(reqex)
			reqmini := Parser.ParserReqMini(v)
			//fmt.Println(reqmini)
			if reqmini == nil {
				continue
			}

			//if reqmini.RoundNum == Global.Connstruct.RoundNum && t==0 {
			//	fmt.Println("value read from conn",string(v))
			//	t++
			//}

			if reqmini.DataType == 3 {
				//fmt.Println("ack",string(v))
				//fmt.Println("22222")
				Buffer.PushIntoQueue(v,connbuffer2,reqmini.RoundNum)
				//fmt.Println(connbuffer2.Size)
				//select {
				//case data := <-bufferchangeBack2:
				//	connbuffer2 = data
				//	//temp := new(Buffer.ConnBuffer)
				//	//temp.Top = connbuffer2.Top
				//	//temp.Tail = connbuffer2.Tail
				//	//temp.Size = connbuffer2.Size
				//	//temp.RWmutex = connbuffer2.RWmutex
				//	Buffer.PushIntoQueue(v,connbuffer2,reqmini.RoundNum)
				//	fmt.Println(connbuffer2.Size)
				//	//if *temp != *connbuffer2 {
				//	//	bufferchange2 <- connbuffer2
				//	//}
				//	//bufferchange2 <- connbuffer2
				//default:
				//	//temp := new(Buffer.ConnBuffer)
				//	//temp.Top = connbuffer2.Top
				//	//temp.Tail = connbuffer2.Tail
				//	//temp.Size = connbuffer2.Size
				//	//temp.RWmutex = connbuffer2.RWmutex
				//	Buffer.PushIntoQueue(v,connbuffer2,reqmini.RoundNum)
				//	fmt.Println(connbuffer2.Size)
				//	//fmt.Println(temp,connbuffer2)
				//	//if *temp != *connbuffer2 {
				//	//	bufferchange2 <- connbuffer2
				//	//}
				//	//bufferchange2 <- connbuffer2
				//}
			} else if reqmini.DataType == 2 {
				//fmt.Println("1111111")
				Buffer.PushIntoQueue(v,connbuffer,reqmini.RoundNum)
				//fmt.Println(connbuffer.Size)
				//fmt.Println("22222222222222")
				//select {
				//case data := <-bufferchangeBack:
				//	//fmt.Println("222222222222222111111111")
				//	connbuffer = data
				//	//temp := connbuffer.Top
				//	//temp := new(Buffer.ConnBuffer)
				//	//temp.Top = connbuffer.Top
				//	//temp.Tail = connbuffer.Tail
				//	//temp.Size = connbuffer.Size
				//	//temp.RWmutex = connbuffer.RWmutex
				//	Buffer.PushIntoQueue(v,connbuffer,reqmini.RoundNum)
				//	fmt.Println(connbuffer.Size)
				//	//if *temp != *connbuffer {
				//	//	bufferchange <- connbuffer
				//	//}
				//	//bufferchange <- connbuffer
				//default:
				//	//fmt.Println("2222222222222222-222222222222")
				//	//temp := new(Buffer.ConnBuffer)
				//	//temp.Top = connbuffer.Top
				//	//temp.Tail = connbuffer.Tail
				//	//temp.Size = connbuffer.Size
				//	//temp.RWmutex = connbuffer.RWmutex
				//	//temp := connbuffer
				//	Buffer.PushIntoQueue(v,connbuffer,reqmini.RoundNum)
				//	fmt.Println(connbuffer.Size)
				//	//fmt.Println("222222222222222-333333333333")
				//	//fmt.Println(temp,connbuffer)
				//	//if *temp != *connbuffer {
				//	//	bufferchange <- connbuffer
				//	//}
				//	//bufferchange <- connbuffer
				//}
				//fmt.Println("222222222222END")
			}
		}
		//Global.DebugLogger <- t
		//t=0
		//fmt.Println("11111111111111111111")
		//fmt.Println("33333")
	}
}

func StartGmae() {
	Global.ConnEstablish <- Global.Connstruct.ConnCount
}