package Forward

import (
	"wGame/Model"
	"fmt"
	"time"
	"wGame/Parser"
	"wGame/Global"
	"bufio"
)
//获取request，构造response
//处理高优先级的超时情况
func ForwardData()  {
	fmt.Println("FordwardData process")
	count := 0
	var res Model.Res
	tempkey := make([]string,0)
	for {
		select {
		case <- Global.Forwardingsignal:
			Global.Connstruct.RWlock.RLock()
			if Global.Connstruct.ConnCount == 0 {
				fmt.Println("no connnect")
				Global.Connstruct.RWlock.RUnlock()
				continue
			} else {
				Global.Connstruct.RWlock.RUnlock()
				res.Result = "timeout"
				res.RoundNum = Global.Connstruct.RoundNum

				Forwarding(res)

				fmt.Println(Global.Connstruct.RoundNum)

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
			res.Content = append(res.Content, Model.Cnt{UserID:data.UserId,Opinions:data.Request.Opinions})
			//组装转发内容体
			//获取5个数据包时，直接转发，并重置计时器，转发结束重置转发内容体
			Global.Connstruct.RWlock.RLock()
			if count == Global.Connstruct.ConnCount {
				Global.Connstruct.RWlock.RUnlock()
				Global.Connstruct.RWlock.RLock()
				res.RoundNum = Global.Connstruct.RoundNum
				Global.Connstruct.RWlock.RUnlock()
				res.Result = "success"
				Forwarding(res)

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
				Global.Connstruct.RWlock.RUnlock()
				Global.Forwardtimer <- 1
			}
		}
	}
}

//转发response
func Forwarding(res Model.Res) {
	//fmt.Println("forwarding")
	Global.Connstruct.RWlock.RLock()
	if Global.Connstruct.ConnCount == 0 {
		Global.Connstruct.RWlock.RUnlock()
		return
	}
	res.DataType = 1
	respt := Parser.CreateRes(res)

	if Global.Connstruct.HaveReConn {
		res.DataType = 8
		Global.Connstruct.HaveReConn = false
	} else {
		res.DataType = 1
	}
	Global.Connstruct.RWlock.RUnlock()

	resp := Parser.CreateRes(res)
	//fmt.Println(string(resp))

	Global.Connstruct.RWlock.Lock()
	if Global.Connstruct.StartStore==true && Global.Connstruct.FlagRoundNum==-1 {
		Global.Connstruct.ReconnData = append(Global.Connstruct.ReconnData, respt...)
	} else if Global.Connstruct.StartStore==true&&res.RoundNum>Global.Connstruct.FlagRoundNum {
		Global.Connstruct.ReconnData = append(Global.Connstruct.ReconnData, respt...)
	}
	for _,v := range Global.Connstruct.Conn {
		remoteAddr := v.RemoteAddr().String()
		//获取每个conn，向每个conn转发
		rw := bufio.NewReadWriter(bufio.NewReader(v),bufio.NewWriter(v))
		_,err := rw.Write(resp)
		if err != nil {
			fmt.Println("Write error",err)
			//loginfo := Log.GetTransferInfo()
			//Global.DebugLogger <- loginfo + err.Error()
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Flush err",err)
			fmt.Println(Global.Connstruct.ConnStatus)
			CloseConn(remoteAddr)
			fmt.Println("Closed")
			//loginfo := Log.GetTransferInfo()
			//Global.DebugLogger <- loginfo + err.Error()
		}
	}
	Global.Connstruct.RoundNum++
	Global.Connstruct.RWlock.Unlock()
	Global.Forwardtimer <- 1
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