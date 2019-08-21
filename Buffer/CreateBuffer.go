package Buffer

import (
	"wGame/Model"
	"wGame/Parser"
	"fmt"
	"sync"
	"wGame/Global"
)

type Node struct {
	RoundNum int64
	Cnt      []byte
	Next     *Node
}

const length = 100
type ConnBuffer struct {
	Top       *Node
	Tail      *Node
	Size      int
	mutex   sync.Mutex
}

func InitQueue() *ConnBuffer {
	var cb ConnBuffer
	cb.Top = nil
	cb.Tail = nil
	cb.Size = 0

	cb.mutex = sync.Mutex{}

	return &cb
}

func PushIntoQueue(cnt []byte,connbuffer *ConnBuffer,roundnum int64) {
	value := cnt
	node := new(Node)
	node.Cnt = value
	node.Next = nil
	node.RoundNum = roundnum
	//fmt.Println("data:",string(node.Cnt),node.RoundNum)
	connbuffer.mutex.Lock()
	defer connbuffer.mutex.Unlock()
	if connbuffer.Size < length {
		if connbuffer.Top == nil{
			if roundnum >= Global.Connstruct.RoundNum {
				connbuffer.Top  = node
				connbuffer.Tail = node
				connbuffer.Size++
				return
			} else {
				return
			}
		} else {
			if roundnum >= Global.Connstruct.RoundNum {
				SortInBuffer(connbuffer,node)
			}
			return
		}
	}
	//TODO-buffer full
	/* else if *size >= length {
		for i := 0; i < 20; i++ {
			top = top.Next
			*size--
		}
		tail.Next = &node
		tail = tail.Next
		*size++
	}*/
	return
}

func PopFromQueue(connbuffer *ConnBuffer) (*Node) {
	connbuffer.mutex.Lock()
	defer connbuffer.mutex.Unlock()
	if connbuffer.Size > 0 && connbuffer.Top != nil {
		temp := connbuffer.Top
		connbuffer.Top = connbuffer.Top.Next
		connbuffer.Size--
		return temp
	}
	return nil
}

func ParserBufferQueue(cnt []byte,remoteAddr string) *Model.ReqEx {
	req := Parser.ParserReq(cnt)
	if req == nil {
		fmt.Println(string(cnt))
		return nil
	}
	//重新组装新格式
	var reqex Model.ReqEx
	reqex.Request    = *req
	reqex.UserId     = req.UserID
	reqex.RemoteAddr = remoteAddr
	reqex.DataType   = req.DataType
	reqex.RoundNum   = req.RoundNum

	return &reqex
}