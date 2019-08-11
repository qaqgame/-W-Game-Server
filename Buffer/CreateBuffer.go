package Buffer

import (
	"wGame/Model"
	"wGame/Parser"
	"fmt"
	"sync"
)

type Node struct {
	Cnt      []byte
	Next     *Node
}

const length = 100
type ConnBuffer struct {
	Top       *Node
	Tail      *Node
	Size      int
	mutex   sync.RWMutex
}

func InitQueue() *ConnBuffer {
	var cb ConnBuffer
	cb.Top = nil
	cb.Tail = nil
	cb.Size = 0
	cb.mutex = sync.RWMutex{}

	return &cb
}

func PushIntoQueue(cnt []byte,connbuffer *ConnBuffer) *ConnBuffer {

	value := cnt
	//fmt.Println("value",string(value))
	//node := Node{value,nil}
	node := new(Node)
	node.Cnt = value
	node.Next = nil
	connbuffer.mutex.RLock()
	//defer connbuffer.mutex.Unlock()
	if connbuffer.Size < length {
		if connbuffer.Top == nil && connbuffer.Size == 0 {
			connbuffer.mutex.RUnlock()

			connbuffer.mutex.Lock()
			connbuffer.Top  = node
			connbuffer.Tail = node
			connbuffer.Size++
			connbuffer.mutex.Unlock()

			//fmt.Println("value", string(connbuffer.Top.Cnt))
			return connbuffer
		} else {
			connbuffer.mutex.RUnlock()

			connbuffer.mutex.Lock()
			//fmt.Println("value", string(connbuffer.Top.Cnt))
			connbuffer.Tail.Next = node
			connbuffer.Tail = connbuffer.Tail.Next
			connbuffer.Size++
			connbuffer.mutex.Unlock()

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

	connbuffer.mutex.RUnlock()
	return connbuffer
}

func PopFromQueue(connbuffer *ConnBuffer) (*Node,*ConnBuffer) {
	connbuffer.mutex.RLock()
	//defer connbuffer.mutex.Unlock()
	if connbuffer.Size > 0 && connbuffer.Top != nil {
		temp := connbuffer.Top
		connbuffer.mutex.RUnlock()

		connbuffer.mutex.Lock()
		connbuffer.Top = connbuffer.Top.Next
		connbuffer.Size--
		connbuffer.mutex.Unlock()

		return temp,connbuffer
	}

	connbuffer.mutex.RUnlock()
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

	return &reqex
}