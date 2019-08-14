package Buffer

func SortInBuffer(connbuffer *ConnBuffer, curr *Node) *ConnBuffer {
	temp := connbuffer.Top
	for {
		if temp.Next == nil {
			if temp.RoundNum < curr.RoundNum {
				temp.Next = curr
				connbuffer.Tail = curr
			} else {
				curr.Next = temp
				connbuffer.Top = curr
			}
			break
		} else {
			if temp.RoundNum < curr.RoundNum && temp.Next.RoundNum > curr.RoundNum {
				curr.Next = temp.Next
				temp.Next = curr
			}
		}
		temp = temp.Next
	}
	connbuffer.Size++
	return connbuffer
}
