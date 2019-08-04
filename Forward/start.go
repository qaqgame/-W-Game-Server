package Forward

import "wGame/Global"

func StartGame(connestablish chan int) {
LOOP:
	for true {
		select {
		case c := <-connestablish:
			if c == 5 {
				//所有人连接上，开启转发计时器和转发器
				go ForwardData(Global.AllDataSlice,Global.Forwardingsignal,Global.Forwardtimer)
				go ForwardingTimer(Global.Forwardtimer,Global.Forwardingsignal)
				break LOOP
			}
		}
	}
	//TODO
	//是否需要发送数据回客户端，告知所有人已经连接上。

	//
}
