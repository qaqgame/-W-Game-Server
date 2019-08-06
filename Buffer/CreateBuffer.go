package Buffer

import (
	"wGame/Model"
)

type Node struct {
	Value    Model.ReqEx
	Next     *Node
}

const length = 100

func InitQueue() (*Node,*Node,*int) {
	var top  *Node = nil
	var tail *Node = nil
	var size       = 0
	return top,tail,&size
}

func PushIntoQueue(reqex Model.ReqEx,top *Node, tail *Node, size *int) []*Node {
	node := Node{reqex,nil}
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
	//fmt.Println("top:",top.Value)
	//fmt.Println("tail:",tail.Value)
	//fmt.Println("size:",*size)
	//fmt.Println()
	return []*Node{top,tail}
}

func PopFromQueue(top *Node,size *int) []*Node {
	if *size > 0 {
		temp := top
		top = top.Next
		*size--
		//fmt.Println("success")
		return []*Node{temp,top}
	}
	//fmt.Println("failed")
	return nil
}