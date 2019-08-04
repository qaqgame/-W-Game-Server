package Forward

import (
	"-w-Game/Model"
	"fmt"
	"net"
	"time"
	"-w-Game/Parser"
	"-w-Game/Global"
)
//获取request，构造response
func ForwardData(reqchan chan Model.ReqEx,send chan int,c chan int)  {
	count := 0
	var res Model.Res
	for {
		select {
		case <-send:
			//到时间转发,生成response
			res.Result = "success"
			resp := Parser.CreateRes(res)
			Forwarding(Global.Conns,resp)
			//清空reqchan
		Loop1:
			for {
				select {
				case <-reqchan:
				default:
					break Loop1
				}
			}
			//转发计时器重新计时
			c <- 1
		case data := <-reqchan:
			count++
			fmt.Println(data)
			//TODO
			//(一次5个或6个？？？，根据时间或者什么标记来判断是否是同一帧的数据？？？？？)
			res.Content = append(res.Content, Model.Cnt{UserID:data.UserId,Opinions:data.Request.Opinions})
			//
		}
	}
}

//转发response
func Forwarding(conns []net.Conn, resp string) {
	for _,v := range conns {
		v.Write([]byte(resp))
	}
}
//服务端转发消息的计时器
//c是计时器使用的计时channel, send是转发时的signal channel
func ForwardingTimer(c chan int,send chan int) {
	timer := time.Duration(5*time.Second)
	t := time.NewTimer(timer)

	defer t.Stop()

	for true {
		select {
		case <-c:
			t.Reset(5*time.Second)
		case <-t.C:
			send <- 1
			fmt.Println("Timeout, Forwarding!")
		}
	}
}