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
	"wGame/Forward"
)

func ConnHandler(conn net.Conn, timerChan chan int) {
	oncefunc := sync.Once{}
	fmt.Println("handleConn")
	//控制子进程退出
	cxt := context.Background()
	cxt,cancel := context.WithCancel(cxt)
	defer cancel()
	connbuffer  := Buffer.InitQueue()
	//connbuffer2 := Buffer.InitQueue()

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
	//go ReadFromBufferQueue2(cxt1,remoteAddr,connbuffer2)
	for true {
		//处理粘包，并读取数据
		result = ReadFromConn(databuf,msgbuf,&length,ulength,conn)
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
				fmt.Println("connclese station")
				fmt.Println(Global.ConnStatus[remoteAddr])
				Forward.CloseConn(remoteAddr)
				return
			}

			//fmt.Println("read from conn",string(v))
			//解析类型和回合数
			reqmini := Parser.ParserReqMini(v)
			//fmt.Println(reqmini)
			if reqmini == nil {
				continue
			}
			//插入buffer
			Buffer.PushIntoQueue(v,connbuffer,reqmini.RoundNum)
		}
	}
}

func StartGmae() {
	Global.ConnEstablish <- Global.Connstruct.ConnCount
}