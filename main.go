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
	"wGame/Buffer"
	"sync"
	"context"
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
	go Forward.StartGame()
	//

	for true {
		conn,err := listener.Accept()
		//类型转换，conn转换为ReadWriter类型
		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		//将rw加入到全局变量Global.Conns中存储
		//Global.Conns = append(Global.Conns, rw)
		Global.Conns[rw] = conn.RemoteAddr().String()
		fmt.Println(len(Global.Conns),"111rewr")
		//统计已有的连接数
		Global.ConnCount++
		Global.ConnEstablish <- Global.ConnCount

		if err != nil {
			fmt.Println(err)
		}
		timerChan := make(chan int,1)
		//开启计时器和连接处理进程
		go Timer(conn,timerChan,rw)
		go handleConn(conn,rw,timerChan)
	}
}

//连接处理进程
func handleConn(conn net.Conn,rw *bufio.ReadWriter,timerChan chan int)  {
	fmt.Println("handleConn")
	cxt := context.Background()
	cxt,cancle := context.WithCancel(cxt)
	//生成buffer队列top改变消息通知channel
	bufferchange     := make(chan *Buffer.Node,1)
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

	go ReadFromBufferQueue(remoteAddr,top,size,bufferchange,bufferchangeBack,mutex)

	for true {
		//处理粘包，并读取数据
		result = ReadFromBuffer(databuf,msgbuf,length,ulength,rw)
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
		if string(result[0]) == "conn close" {
			cancle()
			break
		}

		//处理正常游戏内容数据包
		//解析json数据
		for _,v := range result {
			req := Parser.ParserReq(v)
			//重新组装新格式
			var reqex Model.ReqEx
			reqex.Request    = *req
			reqex.UserId     = req.UserID
			reqex.RemoteAddr = conn.RemoteAddr().String()

			//读取的内容插入缓冲区队列中
			//head := Buffer.PushIntoQueue(reqex)
			select {
			case data := <-bufferchangeBack:
				top = data
				temp := Buffer.PushIntoQueue(reqex,top,tail,size,mutex)
				tail = temp[1]
				bufferchange <- temp[0]
			default:
				temp := Buffer.PushIntoQueue(reqex,top,tail,size,mutex)
				tail = temp[1]
				bufferchange <- temp[0]
			}
		}

	}
}

//从缓冲区读取，与conn读取并行执行
func ReadFromBufferQueue(remoteaddr string,top *Buffer.Node,size *int,bufferchange chan *Buffer.Node,bufferchangeBack chan *Buffer.Node,mutex *sync.Mutex) {
	for true {
		select {
		case <- Global.PlayersChannel[remoteaddr]:
			//fmt.Println("open")
		LOOP:
			for true {
				select {
				case data := <-bufferchange:
					top = data
					temp := Buffer.PopFromQueue(top,size,mutex)
					if temp != nil {
						//fmt.Println(temp[0].Value)
						bufferchangeBack <- temp[1]
						//if temp[0] == nil {
						//	fmt.Println("temp[0] is nil")
						//}
						Global.AllDataSlice <- temp[0].Value
						break LOOP
					}else {
						time.Sleep(5*time.Millisecond)
						continue LOOP
					}
				default:
					select {
					case data := <- bufferchange:
						top = data
						temp := Buffer.PopFromQueue(top,size,mutex)
						if temp != nil {
							//fmt.Println(temp[0].Value)
							bufferchangeBack <- temp[1]
							//if temp[0] == nil {
							//	fmt.Println("temp[0] is nil")
							//}
							Global.AllDataSlice <- temp[0].Value
							break LOOP
						}else {
							time.Sleep(5*time.Millisecond)
							continue LOOP
						}
					default:
						temp := Buffer.PopFromQueue(top,size,mutex)
						if temp != nil {
							//fmt.Println(temp[0].Value)
							Global.AllDataSlice <- temp[0].Value
							break LOOP
						}else {
							time.Sleep(5*time.Millisecond)
							continue LOOP
						}
					}
				}
			}
			//fmt.Println("over")
		}
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
			delete(Global.Conns,rw)
			return
			//
		}
	}

}

func ReadFromBuffer(databuf []byte,msgbuf *bytes.Buffer,length int, ulength uint32, rw *bufio.ReadWriter) [][]byte {
	var result []byte
	result = nil
	var ans [][]byte = nil
	//从reader中读取数据
	for true {
		n,err := rw.Read(databuf)
		if err != nil && err != io.EOF {
			fmt.Println("Error:",err)
			if err.Error() == "read tcp 127.0.0.1:8080->"+Global.Conns[rw]+": wsarecv: An existing connection was forcibly closed by the remote host." {
				fmt.Println("Conn closed")
				delete(Global.Conns,rw)
				return [][]byte{[]byte("conn close")}
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
			ans = append(ans, result)
			//fmt.Println(string(result), msgbuf.Len())
			//msgbuf.Reset()
			length = 0
		} else {
			//fmt.Println("not full data: ",string(result))
			break
		}
	}
	//Global.Count = Global.Count + int64(len(ans))
	//返回最终结果,nil or result
	return ans
}