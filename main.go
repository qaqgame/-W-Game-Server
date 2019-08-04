package main

import (
	"net"
	"fmt"
	"io/ioutil"
	"-w-Game/Parser"
	"time"
	"-w-Game/Model"
	"-w-Game/Global"
	"-w-Game/Forward"
)

func main() {
	service := ":8080"
	tcpAddr,err:= net.ResolveTCPAddr("tcp4",service)
	if err != nil {
		fmt.Println("Resolve tcp error:",err)
	}

	listener, err := net.ListenTCP("tcp",tcpAddr)
	if err != nil {
		fmt.Println("ListenTcp Error:", err)
	}

	//根据连接数，判断开启转发计时器和转发器。
	go Forward.StartGame(Global.ConnEstablish)

	for true {
		conn,err := listener.Accept()
		Global.Conns = append(Global.Conns, conn)

		Global.ConnCount++
		Global.ConnEstablish <- Global.ConnCount

		if err != nil {
			fmt.Println(err)
		}
		timerChan := make(chan int,1)
		go Timer(conn,timerChan)
		go handleConn(conn,timerChan)
	}
}

func handleConn(conn net.Conn,timerChan chan int)  {
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		fmt.Println(err)
	}
	timerChan <- 1
	//TODO
	//判断与处理心跳数据包？？或者先处理连接第一次建立时的初始状态（第一次连接发送初始状态？）？？

	//
	req := Parser.ParserReq(result)
	var reqex Model.ReqEx
	reqex.Request = *req
	reqex.UserId = req.UserID
	Global.AllDataSlice <- reqex
	fmt.Println(*req)
}

//计时器函数，保持长连接，并判断掉线状态
func Timer(conn net.Conn,timerChan chan int) {
	timer := time.Duration(5*time.Second)
	t := time.NewTimer(timer)

	defer t.Stop()
	for true {
		select {
		case <-timerChan:
			t.Reset(5*time.Second)
		case <-t.C:
			fmt.Println("OUTLINE")
			conn.Close()
			//删除连接
			//TODO

			//
		}
	}

}
