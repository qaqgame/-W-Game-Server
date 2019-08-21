package Forward

import (
	"wGame/Global"
)

func CloseConn(remoteAddr string) {
	if Global.ConnStatus[remoteAddr] == 1 {
		//fmt.Println("Close Conn with flush err")
		Global.Connstruct.RWlock.Lock()
		Global.Connstruct.ConnCount--
		Global.Connstruct.Conn[remoteAddr].Close()
		delete(Global.Connstruct.Conn,remoteAddr)
		delete(Global.Connstruct.PlayersChannel,remoteAddr)
		Global.Connstruct.RWlock.Unlock()
		//fmt.Println("Close Conn end")
	}
	delete(Global.ConnStatus,remoteAddr)
}

func CloseConnWhileForWarding(remoteAddr string) {
	if Global.ConnStatus[remoteAddr] == 1 {
		//fmt.Println("Close Conn with flush err")
		Global.Connstruct.ConnCount--
		Global.Connstruct.Conn[remoteAddr].Close()
		delete(Global.Connstruct.Conn,remoteAddr)
		delete(Global.Connstruct.PlayersChannel,remoteAddr)
		//fmt.Println("Close Conn end")
	}
	delete(Global.ConnStatus,remoteAddr)
}