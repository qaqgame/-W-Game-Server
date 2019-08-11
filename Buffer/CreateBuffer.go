package Buffer

import (
	"sync"
	"wGame/Model"
	"wGame/Parser"
)

type Node struct {
	Cnt      []byte
	Next     *Node
}

const length = 100

func InitQueue() (*Node,*Node,*int,*sync.Mutex) {
	var top  *Node = nil
	var tail *Node = nil
	var size       = 0
	var mutex      = sync.Mutex{}
	return top,tail,&size,&mutex
}

func PushIntoQueue(cnt []byte,top *Node, tail *Node, size *int,mutex *sync.Mutex) []*Node {
	//if top == nil {
	//	fmt.Println("top is nil")
	//}
	value := cnt
	node := Node{value,nil}
	mutex.Lock()
	defer mutex.Unlock()
	if *size <= length {
		if top == nil && *size == 0 {
			top  = &node
			tail = &node
			*size++
		} else {
			tail.Next = &node
			tail = tail.Next
			*size++
		}
	}
	return []*Node{top,tail}
}

func PopFromQueue(top *Node,size *int,mutex *sync.Mutex) []*Node {
	mutex.Lock()
	defer mutex.Unlock()
	if *size > 0 && top != nil {
		temp := top
		top = top.Next
		*size--
		return []*Node{temp,top}
	}
	return nil
}

func ParserBufferQueue(cnt []byte,remoteAddr string) *Model.ReqEx {
	req := Parser.ParserReq(cnt)
	if req == nil {
		return nil
	}
	//重新组装新格式
	var reqex Model.ReqEx
	reqex.Request    = *req
	reqex.UserId     = req.UserID
	reqex.RemoteAddr = remoteAddr

	return &reqex
}