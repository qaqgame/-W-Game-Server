package Receive

import (
	"bytes"
	"wGame/Global"
	"encoding/binary"
	"fmt"
	"net"
	"io"
)

func ReadFromConn(databuf []byte, msgbuf *bytes.Buffer, length *int, ulength uint32, conn net.Conn) [][]byte {
	var result []byte = nil
	var ans [][]byte = nil
	remoteaddr := conn.RemoteAddr().String()
	localaddrlen := len(conn.LocalAddr().String())
	//从reader中读取数据
	//n,err := conn.Read(databuf)
	//if err != nil && err != io.EOF {
	//	fmt.Println("Error:",err)
	//	if err.Error()[13+localaddrlen+len(remoteaddr):] == "wsarecv: An existing connection was forcibly closed by the remote host." ||
	//		err.Error()[13+localaddrlen+len(remoteaddr):] == "read: connection reset by peer"{
	//		fmt.Println("Conn closed")
	//		//delete(Global.Conns,rw)
	//		// delete(Global.Conn,remoteaddr)
	//		return [][]byte{[]byte("conn close")}
	//	}
	//}
	//result = append(result,databuf[:n]...)

	for true {
		n,err := conn.Read(databuf)
		if err != nil && err != io.EOF {
			fmt.Println("ReadError:",err)
			if err.Error()[13+localaddrlen+len(remoteaddr):] == "wsarecv: An existing connection was forcibly closed by the remote host." ||
				err.Error()[13+localaddrlen+len(remoteaddr):] == "read: connection reset by peer"{
					fmt.Println("Conn closed")
					//delete(Global.Conns,rw)
					// delete(Global.Conn,remoteaddr)
					return [][]byte{[]byte("conn close")}
			}
		}
		if err != io.EOF{
			result = append(result,databuf[:n]...)
			break
		}
		if n == 0 {
			break
		}
		result = append(result,databuf[:n]...)
	}
	//Global.Count++
	//fmt.Println(string(result))

	_,err := msgbuf.Write(result)
	if err != nil {
		fmt.Println("Buffer write error: ",err)
		//loginfo := Log.GetTransferInfo()
		//Global.DebugLogger <- loginfo + err.Error()
	}
	//处理粘包
	for true {
		if *length == 0 && msgbuf.Len() >= 4 {
			binary.Read(msgbuf,binary.LittleEndian,&ulength)
			*length = int(ulength)
			if *length > 10240 {
				fmt.Printf("Message too length: %d\n", length)
			}
		}
		if *length > 0 && msgbuf.Len() >= *length {
			result = msgbuf.Next(*length)
			ans = append(ans, result)
			*length = 0
		} else {
			break
		}
	}
	//返回最终结果,nil or result

	Global.Count = Global.Count + len(ans)
	//fmt.Println(ans)
	return ans
}
