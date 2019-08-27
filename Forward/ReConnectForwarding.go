package Forward

import (
	"net"
	"wGame/Global"
	"fmt"
	"wGame/Parser"
	"bufio"
)

func ReconnectForwarding(conn net.Conn) {
	fmt.Println("Run ReconnectForwarding")
	select {
	case data := <-Global.ReConnData:
		//fmt.Println(string(data))
		statesync := Parser.ParserStateSync(data)
		statesync.DataType = 7
		resp := Parser.CreateReconndata(*statesync)
		//fmt.Println(string(resp))
		rw := bufio.NewReadWriter(bufio.NewReader(conn),bufio.NewWriter(conn))
		_,err := rw.Write(append(resp, Global.Connstruct.ReconnData...))
		fmt.Println(string(append(resp, Global.Connstruct.ReconnData...)))
		if err != nil {
			fmt.Println("Reconnect write error",err)
		}
		err = rw.Flush()
		if err != nil {
			fmt.Println("Reconnect flush error",err)
		}
		fmt.Println("successfully sended type7")
	}

	//在这里添加一个人？？？？
	Global.Connstruct.RWlock.Lock()
	Global.Connstruct.StartStore = false
	Global.Connstruct.FlagRoundNum = -1
	Global.Connstruct.ReconnData = nil
	Global.Connstruct.Conn[conn.RemoteAddr().String()] = conn
	Global.Connstruct.ConnStatus[conn.RemoteAddr().String()] = 1
	Global.Connstruct.ConnCount++
	Global.Connstruct.RWlock.Unlock()
}
