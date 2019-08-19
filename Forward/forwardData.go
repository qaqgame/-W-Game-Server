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
	fmt.Println("FordwardData process")
	//infotype := 2
	count := 0
	var res Model.Res
	tempkey := make([]string,0)
	for {
		select {
		case <- Global.Forwardingsignal:
			//fmt.Println("timeout fording")
			//Global.Connstruct.RWlock.RLock()
			if Global.Connstruct.ConnCount == 0 {
				fmt.Println("no connnect")
				//Global.Connstruct.RWlock.RUnlock()
				continue
			} else {
				//if infotype == 2 {
				//	res.DataType = 1
				//} else if infotype == 3 {
				//	res.DataType = 4
				//}
				if count == 0 {
					res.Result = "failed"
				} else {
					res.Result = "success"
				}
				res.DataType = 1
				res.RoundNum = Global.Connstruct.RoundNum
				resp := Parser.CreateRes(res)
				//fmt.Println("timeout resp:",resp)
				if resp == "" {
					count = 0
					Global.Connstruct.RWlock.Lock()
					Global.Connstruct.RoundNum++
					Global.Connstruct.RWlock.Unlock()
					Global.Forwardingsignal <- 1
					continue
				}
				//fmt.Println("Roundnum:",Global.Connstruct.RoundNum)
				//fmt.Println("before forward1",infotype)
				Forwarding(Global.Connstruct.Conn,resp)

				Global.Connstruct.RWlock.Lock()
				Global.Connstruct.RoundNum++
				Global.Connstruct.RWlock.Unlock()
				fmt.Println(Global.Connstruct.RoundNum)

				//if infotype == 3 {
				//	Global.Connstruct.RWlock.Lock()
				//	Global.Connstruct.RoundNum++
				//	Global.Connstruct.RWlock.Unlock()
				//	fmt.Println(Global.Connstruct.RoundNum)
				//}
				//使用缓冲区，发送信号从缓冲区中读取数据
				for _,v := range tempkey {
					if _,ok:= Global.Connstruct.PlayersChannel[v];ok {
						Global.Connstruct.PlayersChannel[v] <- 1
					}
				}
				//Global.Connstruct.RWlock.Lock()
				//if infotype == 2 {
				//	for _,v := range tempkey {
				//		if _,ok:= Global.Connstruct.PlayersChannelAck[v];ok {
				//			Global.Connstruct.PlayersChannelAck[v] <- 1
				//		}
				//	}
				//	infotype = 3
				//} else if infotype == 3 {
				//	for _,v := range tempkey {
				//		if _,ok:= Global.Connstruct.PlayersChannel[v];ok {
				//			Global.Connstruct.PlayersChannel[v] <- 1
				//		}
				//	}
				//	infotype = 2
				//}
				//Global.Connstruct.RWlock.Unlock()
				tempkey = nil
				count = 0
				//转发计时器重新计时，转发内容重置
				res.Content = nil
				Global.Forwardtimer <- 1
			}
		case data := <- Global.AllDataSlice:
			//fmt.Println("data received forward")
			//fmt.Println("data",data)
			res.DataType = 1

			//if infotype == 2 {
			//	res.DataType = 1
			//} else if infotype == 3 {
			//	res.DataType = 4
			//}

			tempkey = append(tempkey, data.RemoteAddr)
			count++
			//if data.DataType == infotype {
			//	res.Content = append(res.Content, Model.Cnt{UserID:data.UserId,Opinions:data.Request.Opinions})
			//}


			res.Content = append(res.Content, Model.Cnt{UserID:data.UserId,Opinions:data.Request.Opinions})
			//组装转发内容体
			//获取5个数据包时，直接转发，并重置计时器，转发结束重置转发内容体
			//Global.Connstruct.RWlock.RLock()
			if count == Global.Connstruct.ConnCount {
				//Global.Connstruct.RWlock.RUnlock()
				Global.Connstruct.RWlock.RLock()
				res.RoundNum = Global.Connstruct.RoundNum
				Global.Connstruct.RWlock.RUnlock()
				//res.DataType = 1
				res.Result = "success"
				resp := Parser.CreateRes(res)
				if resp == "" {
					Global.Connstruct.RWlock.Lock()
					Global.Connstruct.RoundNum++
					Global.Connstruct.RWlock.Unlock()
					Global.Forwardingsignal <- 1
					count = 0
					continue
				}
				//fmt.Println("before forward",infotype)
				Forwarding(Global.Connstruct.Conn,resp)

				Global.Connstruct.RWlock.Lock()
				Global.Connstruct.RoundNum++
				Global.Connstruct.RWlock.Unlock()

				//if infotype == 3 {
				//	Global.Connstruct.RWlock.Lock()
				//	Global.Connstruct.RoundNum++
				//	Global.Connstruct.RWlock.Unlock()
				//}
				res.Content = nil
				count = 0
				//使用缓冲区，发送信号从缓冲区中读取数据
				for _,v := range tempkey {
					if _,ok:= Global.Connstruct.PlayersChannel[v];ok {
						Global.Connstruct.PlayersChannel[v] <- 1
					}
				}
				//Global.Connstruct.RWlock.Lock()
				//if infotype == 2 {
				//	for _,v := range tempkey {
				//		if _,ok:= Global.Connstruct.PlayersChannelAck[v];ok {
				//			Global.Connstruct.PlayersChannelAck[v] <- 1
				//		}
				//	}
				//	infotype = 3
				//} else if infotype == 3 {
				//	for _,v := range tempkey {
				//		if _,ok:= Global.Connstruct.PlayersChannel[v];ok {
				//			Global.Connstruct.PlayersChannel[v] <- 1
				//		}
				//	}
				//	infotype = 2
				//}
				//Global.Connstruct.RWlock.Unlock()
				//fmt.Println("?????")
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
	//fmt.Println("forwarding")
	Global.Connstruct.RWlock.RLock()
	if Global.Connstruct.ConnCount == 0 {
		Global.Connstruct.RWlock.RUnlock()
		return
	}
	Global.Connstruct.RWlock.RUnlock()
	fmt.Println(resp)
	Global.Connstruct.RWlock.Lock()
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
	Global.Connstruct.RWlock.Unlock()
}
//服务端转发消息的计时器
//c是计时器使用的计时channel, send是转发时的signal channel
func ForwardingTimer() {
	//fmt.Println("timer",time.Now().Format(time.RFC3339Nano))
	timer := time.Duration(500*time.Millisecond)
	t := time.NewTimer(timer)

	defer t.Stop()

	for true {
		select {
		case <-Global.Forwardtimer:
			//fmt.Println("reset")
			t.Reset(timer)
		case <-t.C:
			fmt.Println("Timeout, Forwarding!")
			Global.Forwardingsignal <- 1
		}
	}
}