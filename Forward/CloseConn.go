package Forward

import (
	"wGame/Global"
)

func CloseConn(remoteAddr string) {
	if Global.Connstruct.ConnStatus[remoteAddr] == 1 {
		//fmt.Println("Close Conn with flush err")
		Global.Connstruct.ConnCount--
		Global.Connstruct.Conn[remoteAddr].Close()
		delete(Global.Connstruct.Conn,remoteAddr)
		delete(Global.Connstruct.PlayersChannel,remoteAddr)
		Global.Connstruct.ConnStatus[remoteAddr] = 0
		//fmt.Println("Close Conn end")
	}
	delete(Global.Connstruct.ConnStatus,remoteAddr)
}