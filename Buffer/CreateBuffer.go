package Buffer

import (
	"wGame/Model"
	"sync"
	"fmt"
)

type Node struct {
	Value    Model.ReqEx
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

func PushIntoQueue(reqex Model.ReqEx,top *Node, tail *Node, size *int,mutex *sync.Mutex) []*Node {
	//if top == nil {
	//	fmt.Println("top is nil")
	//}
	node := Node{reqex,nil}
	if *size <= length {
		//fmt.Println(reqex)
		if top == nil && *size == 0 {
			mutex.Lock()
			top  = &node
			tail = &node
			*size++
			mutex.Unlock()
		} else {
			mutex.Lock()
			tail.Next = &node
			tail = tail.Next
			*size++
			mutex.Unlock()
		}
	}
	if top == nil {
		fmt.Println("ffffffffffffnil")
	}
	//fmt.Println("top:",top.Value)
	//fmt.Println("tail:",tail.Value)
	//fmt.Println("size:",*size)
	//fmt.Println()
	return []*Node{top,tail}
}

func PopFromQueue(top *Node,size *int,mutex *sync.Mutex) []*Node {
	if *size > 0 && top != nil {
		mutex.Lock()
		temp := top
		top = top.Next
		*size--
		mutex.Unlock()
		//fmt.Println("success")
		return []*Node{temp,top}
	}
	//fmt.Println("failed")
	return nil
}