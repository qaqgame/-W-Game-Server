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
	count := 0
	var res Model.Res
	tempkey := make([]string,0)
	for {
		select {
		case <- Global.Forwardingsignal:
			if Global.ConnCount == 0 {
				continue
			}
			//到时间转发,生成response
			fmt.Println("count: ",count,"Global Count:",Global.Count)
			res.Result = "success"
			resp := Parser.CreateRes(res)
			Forwarding(Global.Conns,resp)
			count = 0

			//转发计时器重新计时，转发内容重置
			res.Content = nil
			Global.Forwardtimer <- 1
			//使用缓冲区，发送信号从缓冲区中读取数据
			for _,v := range tempkey {
				Global.PlayersChannel[v] <- 1
			}
			tempkey = nil
		case data := <- Global.AllDataSlice:
			//fmt.Println(1,count)
			tempkey = append(tempkey, data.RemoteAddr)
			//fmt.Println(len(tempkey))
			count++
			fmt.Println(count,"count")
			//组装转发内容体
			res.Content = append(res.Content, Model.Cnt{UserID:data.UserId,Opinions:data.Request.Opinions})
			//获取5个数据包时，直接转发，并重置计时器，转发结束重置转发内容体
			if count == Global.ConnCount {
				res.Result = "success"
				resp := Parser.CreateRes(res)
				Forwarding(Global.Conns,resp)
				count = 0
				Global.Forwardtimer <- 1
				res.Content = nil
				//使用缓冲区，发送信号从缓冲区中读取数据
				for _,v := range tempkey {
					Global.PlayersChannel[v] <- 1
				}
				tempkey = nil
			}
			Global.Forwardtimer <- 1
		}
	}
}

//转发response
func Forwarding(conns map[*bufio.ReadWriter]string, resp string) {
	fmt.Println("forwarding")
	if Global.ConnCount == 0 {
		return
	}
	fmt.Println(resp)
	for rw,_ := range conns {
		//获取每个conn，向每个conn转发
		_,err := rw.Write([]byte(resp))
		if err != nil {
			fmt.Println(err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println(err)
		}
	}
}
//服务端转发消息的计时器
//c是计时器使用的计时channel, send是转发时的signal channel
func ForwardingTimer() {
	timer := time.Duration(15*time.Millisecond)
	t := time.NewTimer(timer)

	defer t.Stop()

	for true {
		select {
		case <- Global.Forwardtimer:
			//fmt.Println("reset")
			t.Reset(timer)
		case <-t.C:
			fmt.Println("Timeout, Forwarding!")
			Global.Forwardingsignal <- 1
		}
	}
}