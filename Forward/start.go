package Forward

import (
	"wGame/Global"
	"bufio"
	"fmt"
)

func StartGame() {
LOOP:
	for true {
		select {
		//监听连接建立数的channel
		case c := <- Global.ConnEstablish:
			//五人都建立连接，向客户端发回建立成功的数据，开启转发计时器和转发器
			//TODO
			//思考如何在某些连接断开后能正常返回数据

			//！！！
			if c == Global.PlayerNum {
				//所有人连接上，开启转发计时器和转发器
				for _,v := range Global.Connstruct.Conn {
					rw := bufio.NewReadWriter(bufio.NewReader(v),bufio.NewWriter(v))
					_,err := rw.Write([]byte("GameStart"))
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
				go ForwardData()
				fmt.Println("111212")
				go ForwardingTimer()
				break LOOP
			}
		}
	}
	return
}
