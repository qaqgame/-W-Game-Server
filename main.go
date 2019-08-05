package main

import (
	"net"
	"fmt"
	"wGame/Parser"
	"time"
	"wGame/Model"
	"wGame/Global"
	"wGame/Forward"
	"io"
	"bytes"
	"encoding/binary"
	"bufio"
)


func main() {
	Global.Count = 0                      //测试时使用的临时数据
	service := ":8080"
	tcpAddr,err:= net.ResolveTCPAddr("tcp4",service)
	if err != nil {
		fmt.Println("Resolve tcp error:",err)
	}
	listener, err := net.ListenTCP("tcp",tcpAddr)
	if err != nil {
		fmt.Println("ListenTcp Error:", err)
	}
	defer listener.Close()
	//根据连接数，判断开启转发计时器和转发器。都连接上时开始GAME
	go Forward.StartGame(Global.ConnEstablish)

	for true {
		conn,err := listener.Accept()
		//类型转换，conn转换为ReadWriter类型
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		//将rw加入到全局变量Global.Conns中存储
		//Global.Conns = append(Global.Conns, rw)
		Global.Conns[rw] = conn.RemoteAddr().String()
		fmt.Println(len(Global.Conns),"111rewr")
		//通缉已有的连接数
		Global.ConnCount++
		Global.ConnEstablish <- Global.ConnCount

		if err != nil {
			fmt.Println(err)
		}
		timerChan := make(chan int,1)
		//开启计时器和连接处理进程
		go Timer(conn,timerChan,rw)
		go handleConn(rw,timerChan)
	}
}

//连接处理进程
func handleConn(rw *bufio.ReadWriter,timerChan chan int)  {
	fmt.Println("handleConn")

	var result []byte
	msgbuf := bytes.NewBuffer(make([]byte,0,10240))
	databuf := make([]byte,4096)
	length := 0
	ulength := uint32(0)

	for true {
		//处理粘包，并读取数据
		result = ReadFromBuffer(databuf,msgbuf,length,ulength,rw)
		//重置计时器
		timerChan <- 1
		//判断与处理心跳数据包，此时为所有玩家都成功连接前的状态，
		// 使用心跳数据包来保证已建立的连接的在线状态
		if string(result) == "heart beats" {
			fmt.Println("heart beats")
			continue
		}
		if string(result) == "conn close" {
			break
		}
		//处理正常游戏内容数据包
		//解析json数据
		req := Parser.ParserReq(result)
		//重新组装新格式
		var reqex Model.ReqEx
		reqex.Request = *req
		reqex.UserId = req.UserID
		//读取的内容push进入channel
		Global.AllDataSlice <- reqex
		fmt.Println(*req)
	}
}

//计时器函数，保持长连接，并判断掉线状态
func Timer(conn net.Conn,timerChan chan int,rw *bufio.ReadWriter) {
	timer := time.Duration(3*time.Second)
	t := time.NewTimer(timer)

	defer t.Stop()
	for true {
		select {
		case <-timerChan:
			t.Reset(3*time.Second)
		case <-t.C:
			fmt.Println("OUTLINE")
			conn.Close()    //关闭连接
			//删除连接
			//TODO
			delete(Global.Conns,rw)
			return
			//
		}
	}

}

func ReadFromBuffer(databuf []byte,msgbuf *bytes.Buffer,length int, ulength uint32, rw *bufio.ReadWriter) []byte {
	var result []byte
	result = nil
	//从reader中读取数据
	for true {
		n,err := rw.Read(databuf)
		if err != nil && err != io.EOF {
			fmt.Println("Error:",err)
			if err.Error() == "read tcp 127.0.0.1:8080->"+Global.Conns[rw]+": wsarecv: An existing connection was forcibly closed by the remote host." {
				fmt.Println("Conn closed")
				delete(Global.Conns,rw)
				return []byte("conn close")
			}
		}
		if err != io.EOF{
			result = append(result,databuf[:n]...)
			break
		}
		if n == 0 {
			break
		}
		result = append(result,databuf[:n]...)
	}
	Global.Count++
	_,err := msgbuf.Write(result)
	if err != nil {
		fmt.Println("Buffer write error: ",err)
	}
	//处理粘包
	for true {
		if length == 0 && msgbuf.Len() >= 4 {
			binary.Read(msgbuf,binary.LittleEndian,&ulength)
			length = int(ulength)

			if length > 10240 {
				fmt.Printf("Message too length: %d\n", length)
			}
		}
		if length > 0 && msgbuf.Len() >= length {
			result = msgbuf.Next(length)
			fmt.Println(string(result), msgbuf.Len())
			//msgbuf.Reset()
			length = 0
		} else {
			break
		}
	}
	//返回最终结果
	return result
}