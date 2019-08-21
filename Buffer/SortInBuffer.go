package Buffer

import (
	"wGame/Global"
)

func SortInBuffer(connbuffer *ConnBuffer, curr *Node) {
	temp := connbuffer.Top
	if curr.RoundNum < Global.Connstruct.RoundNum {
		return
	}
	if temp == nil && connbuffer.Size<100 {
		connbuffer.Size = 0
		if curr.RoundNum >= Global.Connstruct.RoundNum {
			temp = curr
			connbuffer.Tail = curr
			connbuffer.Size++
			return
		}
	}
	for {
		if temp.Next == nil {
			if temp.RoundNum < curr.RoundNum {
				temp.Next = curr
				connbuffer.Tail = curr
				break
			} else if temp.RoundNum > curr.RoundNum {
				curr.Next = temp
				connbuffer.Top = curr
				break
				//
			} else {
				return
			}
		} else {
			if temp.RoundNum < curr.RoundNum && temp.Next.RoundNum > curr.RoundNum {
				curr.Next = temp.Next
				temp.Next = curr
				break
			} else if temp.RoundNum == curr.RoundNum || temp.Next.RoundNum == curr.RoundNum {
				return
			}
		}
		temp = temp.Next
	}
	connbuffer.Size++
	return
}
