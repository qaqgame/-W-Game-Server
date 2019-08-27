package main

import (
	"net"
	"fmt"
	"time"
	"wGame/Global"
	"wGame/Forward"
	"wGame/Receive"
)

func init() {
	//go Log.LogController()
	go Receive.StateSyncDataHandler()
}

func main() {
	Global.Count = 0                      //测试时使用的临时数据,统计接收了多少数据
	service := ":8080"
	tcpAddr,err:= net.ResolveTCPAddr("tcp4",service)
	if err != nil {
		fmt.Println("Resolve tcp error:",err)
		//loginfo := Log.GetTransferInfo()
		//Global.DebugLogger <- loginfo + err.Error()
	}
	listener, err := net.ListenTCP("tcp",tcpAddr)
	if err != nil {
		fmt.Println("ListenTcp Error:", err)
		//loginfo := Log.GetTransferInfo()
		//Global.DebugLogger <- loginfo + err.Error()
	}
	defer listener.Close()


	//根据连接数，判断开启转发计时器和转发器。都连接上时开始GAME
	go Forward.StartGame()
	//


	for true {
		conn,err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			//loginfo := Log.GetTransferInfo()
			//Global.DebugLogger <- loginfo + err.Error()
		}

		timerChan := make(chan int,1)
		//开启计时器和连接处理进程
		go Timer(conn,timerChan)
		go Receive.ConnHandler(conn,timerChan)
	}
}

//计时器函数，保持长连接，并判断掉线状态
func Timer(conn net.Conn,timerChan chan int) {
	remoteAddr := conn.RemoteAddr().String()
	timer := time.Duration(10*time.Second)
	t := time.NewTimer(timer)

	defer t.Stop()
	for true {
		select {
		case <-timerChan:
			t.Reset(timer)
		case <-t.C:
			fmt.Println("OUTLINE")
			//conn.Close()    //关闭连接
			//删除连接
			Forward.CloseConn(remoteAddr)
			return
			//
		}
	}
}