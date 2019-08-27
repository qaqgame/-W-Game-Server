package Receive

import (
	"context"
	"wGame/Buffer"
	"fmt"
	"wGame/Global"
)

func ReadFromBufferQueue2(cxt context.Context, remoteaddr string, connbuffer2 *Buffer.ConnBuffer) {
	var temp *Buffer.Node = nil
	if connbuffer2 == nil {
		fmt.Println("11111")
	}
	for true {
		select {
		case <-Global.Connstruct.PlayersChannelAck[remoteaddr]:
		LOOP:
			for true {
				select {
				case <-cxt.Done():
					return
				default:
					select {
					case <-cxt.Done():
						return

					default:
						temp = nil
						temp = Buffer.PopFromQueue(connbuffer2)
						if temp != nil {
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
							Global.Connstruct.RWlock.RLock()
							if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
								Global.Connstruct.RWlock.RUnlock()
								Global.AllDataSlice <- *value
								break LOOP
							} else {
								if value.RoundNum < Global.Connstruct.RoundNum {
									Global.Connstruct.RWlock.RUnlock()
									Buffer.PopFromQueue(connbuffer2)
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