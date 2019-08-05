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
func ForwardData(reqchan chan Model.ReqEx,send chan int,c chan int)  {
	count := 0
	var res Model.Res
	for {
		select {
		case <-send:
			//到时间转发,生成response
			fmt.Println("count: ",count,"Global Count:",Global.Count)
			res.Result = "success"
			resp := Parser.CreateRes(res)
			Forwarding(Global.Conns,resp)
			count = 0
			//TODO
			//清空reqchan（考虑是否需要清空）
		Loop1:
			for {
				select {
				case <-reqchan:
					fmt.Println("clear reqchan")
				default:
					break Loop1
				}
			}
			//转发计时器重新计时，转发内容重置
			res.Content = nil
			c <- 1
		case data := <-reqchan:
			count++
			fmt.Println(count,"count")
			//组装转发内容体
			res.Content = append(res.Content, Model.Cnt{UserID:data.UserId,Opinions:data.Request.Opinions})
			//获取5个数据包时，直接转发，并重置计时器，转发结束重置转发内容体
			if count == 5 {
				res.Result = "success"
				resp := Parser.CreateRes(res)
				Forwarding(Global.Conns,resp)
				count = 0
				c <- 1
				res.Content = nil
			}
		}
	}
}

//转发response
func Forwarding(conns map[*bufio.ReadWriter]string, resp string) {
	fmt.Println("forwarding")
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
func ForwardingTimer(c chan int,send chan int) {
	timer := time.Duration(20*time.Millisecond)
	t := time.NewTimer(timer)

	defer t.Stop()

	for true {
		select {
		case <-c:
			t.Reset(20*time.Millisecond)
		case <-t.C:
			fmt.Println("Timeout, Forwarding!")
			send <- 1
		}
	}
}