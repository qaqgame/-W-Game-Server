package Forward

import (
	"wGame/Model"
	"fmt"
	"time"
	"wGame/Parser"
	"wGame/Global"
	"net"
	"bufio"
)
//获取request，构造response
//处理高优先级的超时情况
func ForwardData()  {
	count := 0
	var res Model.Res
	tempkey := make([]string,0)
	for {
		select {
		case <- Global.Forwardingsignal:
			//Global.Connstruct.RWlock.RLock()
			if Global.Connstruct.ConnCount == 0 {
				//Global.Connstruct.RWlock.RUnlock()
				continue
			} else {
				res.Result = "success"
				resp := Parser.CreateRes(res)
				if resp == "" {
					count = 0
					Global.Forwardingsignal <- 1
					continue
				}
				Forwarding(Global.Connstruct.Conn,resp)
				//使用缓冲区，发送信号从缓冲区中读取数据
				for _,v := range tempkey {
					if _,ok:= Global.Connstruct.PlayersChannel[v];ok {
						Global.Connstruct.PlayersChannel[v] <- 1
					}
				}
				tempkey = nil
				count = 0
				//转发计时器重新计时，转发内容重置
				res.Content = nil
				Global.Forwardtimer <- 1
			}

		case data := <- Global.AllDataSlice:
			tempkey = append(tempkey, data.RemoteAddr)
			count++
			//组装转发内容体
			res.Content = append(res.Content, Model.Cnt{UserID:data.UserId,Opinions:data.Request.Opinions})
			//获取5个数据包时，直接转发，并重置计时器，转发结束重置转发内容体
			//Global.Connstruct.RWlock.RLock()
			if count == Global.Connstruct.ConnCount {
				//Global.Connstruct.RWlock.RUnlock()
				res.Result = "success"
				resp := Parser.CreateRes(res)
				if resp == "" {
					Global.Forwardingsignal <- 1
					count = 0
					continue
				}
				Forwarding(Global.Connstruct.Conn,resp)
				res.Content = nil
				count = 0
				//使用缓冲区，发送信号从缓冲区中读取数据
				for _,v := range tempkey {
					if _,ok:= Global.Connstruct.PlayersChannel[v];ok {
						Global.Connstruct.PlayersChannel[v] <- 1
					}
				}
				tempkey = nil
				Global.Forwardtimer <- 1
			} else {
				//Global.Connstruct.RWlock.RUnlock()
				Global.Forwardtimer <- 1
			}
		}
	}
}

//转发response
func Forwarding(conn map[string]net.Conn, resp string) {
	//Global.Connstruct.RWlock.RLock()
	if Global.Connstruct.ConnCount == 0 {
		//Global.Connstruct.RWlock.RUnlock()
		return
	}
	//fmt.Println(resp)
	for i,v := range conn {
		//获取每个conn，向每个conn转发
		if _,ok := conn[i];ok {
			rw := bufio.NewReadWriter(bufio.NewReader(v),bufio.NewWriter(v))
			_,err := rw.Write([]byte(resp))
			if err != nil {
				fmt.Println(err)
				//loginfo := Log.GetTransferInfo()
				//Global.DebugLogger <- loginfo + err.Error()
			}
			err = rw.Flush()
			if err != nil {
				fmt.Println(err)
				//loginfo := Log.GetTransferInfo()
				//Global.DebugLogger <- loginfo + err.Error()
			}
		}
	}
}
//服务端转发消息的计时器
//c是计时器使用的计时channel, send是转发时的signal channel
func ForwardingTimer() {
	timer := time.Duration(70*time.Millisecond)
	t := time.NewTimer(timer)

	defer t.Stop()

	for true {
		select {
		case <- Global.Forwardtimer:
			t.Reset(timer)
		case <-t.C:
			fmt.Println("Timeout, Forwarding!")
			Global.Forwardingsignal <- 1
		}
	}
}