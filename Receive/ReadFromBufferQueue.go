package Receive

import (
	"wGame/Buffer"
	"context"
	"wGame/Global"
	"fmt"
)

func ReadFromBufferQueue(cxt context.Context, remoteaddr string,connbuffer *Buffer.ConnBuffer) {
	var temp *Buffer.Node = nil
	if connbuffer == nil {
		fmt.Println("oooo")
	}
	for true {
		select {
		case <- Global.Connstruct.PlayersChannel[remoteaddr]:
		LOOP:
			for true {
				select {
				//退出线程
				case <-cxt.Done():
					return
				default:
					select {
					case <-cxt.Done():
						return
					default:
						temp = nil
						temp = Buffer.PopFromQueue(connbuffer)
						if temp != nil {
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
							if value != nil {
								Global.Forwardtimer <- 1
							}
							Global.Connstruct.RWlock.RLock()
							if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
								Global.Connstruct.RWlock.RUnlock()
								Global.AllDataSlice <- *value
								break LOOP
							} else {
								if value.RoundNum < Global.Connstruct.RoundNum {
									Global.Connstruct.RWlock.RUnlock()
									Buffer.PopFromQueue(connbuffer)
									continue LOOP
								}
								Global.Connstruct.RWlock.RUnlock()
								continue LOOP
							}
						}else {
							continue LOOP
						}
					}
				}
			}
		}
	}
}
