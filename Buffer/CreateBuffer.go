package Buffer

import (
	"wGame/Model"
	"wGame/Parser"
	"fmt"
	"sync"
)

type Node struct {
	RoundNum int64
	DataType int
	Cnt      []byte
	Next     *Node
}

const length = 100
type ConnBuffer struct {
	Top       *Node
	Tail      *Node
	Size      int
	RWmutex   sync.RWMutex
}

func InitQueue() *ConnBuffer {
	var cb ConnBuffer
	cb.Top = nil
	cb.Tail = nil
	cb.Size = 0

	cb.RWmutex = sync.RWMutex{}

	return &cb
}

func PushIntoQueue(cnt []byte,connbuffer *ConnBuffer,roundnum int64) *ConnBuffer {

	value := cnt
	//fmt.Println("Push value",string(value),"size",connbuffer.Size)
	//fmt.Println("size",connbuffer.Size,roundnum)
	//node := Node{value,nil}
	node := new(Node)
	node.Cnt = value
	node.Next = nil
	node.RoundNum = roundnum
	connbuffer.RWmutex.RLock()
	//defer connbuffer.mutex.Unlock()
	if connbuffer.Size < length {
		if connbuffer.Top == nil && connbuffer.Size == 0 {
			connbuffer.RWmutex.RUnlock()

			connbuffer.RWmutex.Lock()
			connbuffer.Top  = node
			connbuffer.Tail = node
			connbuffer.Size++
			connbuffer.RWmutex.Unlock()

			//fmt.Println("value", string(connbuffer.Top.Cnt))
			return connbuffer
		} else {
			connbuffer.RWmutex.RUnlock()

			connbuffer.RWmutex.Lock()
			//fmt.Println("value", string(connbuffer.Top.Cnt))
			connbuffer = SortInBuffer(connbuffer,node)
			//connbuffer.Tail.Next = node
			//connbuffer.Tail = connbuffer.Tail.Next
			//connbuffer.Size++
			connbuffer.RWmutex.Unlock()

			return connbuffer
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

	connbuffer.RWmutex.RUnlock()
	return connbuffer
}

func PopFromQueue(connbuffer *ConnBuffer) (*Node,*ConnBuffer) {
	connbuffer.RWmutex.RLock()
	//defer connbuffer.mutex.Unlock()
	if connbuffer.Size > 0 && connbuffer.Top != nil {
		temp := connbuffer.Top
		connbuffer.RWmutex.RUnlock()

		connbuffer.RWmutex.Lock()
		connbuffer.Top = connbuffer.Top.Next
		connbuffer.Size--
		connbuffer.RWmutex.Unlock()

		return temp,connbuffer
	}

	connbuffer.RWmutex.RUnlock()
	return nil,connbuffer
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