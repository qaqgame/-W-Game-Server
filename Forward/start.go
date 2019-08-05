package Forward

import (
	"wGame/Global"
	"fmt"
)

func StartGame(connestablish chan int) {
LOOP:
	for true {
		select {
		//监听连接建立数的channel
		case c := <-connestablish:
			//五人都建立连接，向客户端发回建立成功的数据，开启转发计时器和转发器
			if c == 5 {
				//所有人连接上，开启转发计时器和转发器
				for rw,_ := range Global.Conns {
					_,err := rw.Write([]byte("all player connected"))
					if err != nil {
						fmt.Println("Error:",err)
					}
					err = rw.Flush()
					if err != nil {
						fmt.Println(err)
					}
				}
				go ForwardData(Global.AllDataSlice,Global.Forwardingsignal,Global.Forwardtimer)
				go ForwardingTimer(Global.Forwardtimer,Global.Forwardingsignal)
				break LOOP
			}
		}
	}
}
