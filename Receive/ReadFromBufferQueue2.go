package Receive

import (
	"context"
	"wGame/Buffer"
	"fmt"
	"wGame/Global"
	"time"
)

func ReadFromBufferQueue2(cxt context.Context, remoteaddr string, connbuffer2 *Buffer.ConnBuffer, bufferchange2 chan *Buffer.ConnBuffer, bufferchangeBack2 chan *Buffer.ConnBuffer) {
	var temp *Buffer.Node = nil
	if connbuffer2 == nil {
		fmt.Println("11111")
	}
	for true {
		select {
		case <-Global.Connstruct.PlayersChannelAck[remoteaddr]:
			//fmt.Println("reqack")
		LOOP:
			for true {
				select {
				case <-cxt.Done():
					return
				default:
					select {
					case <-cxt.Done():
						return
					case data := <-bufferchange2:
						temp = nil
						connbuffer2 = data
						temp,connbuffer2 = Buffer.PopFromQueue(connbuffer2)
						if temp != nil {
							bufferchangeBack2 <- connbuffer2
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
							Global.Connstruct.RWlock.RLock()
							//fmt.Println("value22222:",value,value.RoundNum,Global.Connstruct.RoundNum,value.DataType)
							if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
								Global.Connstruct.RWlock.RUnlock()
								Global.AllDataSlice <- *value
								break LOOP
							} else {
								if value.RoundNum < Global.Connstruct.RoundNum {
									Global.Connstruct.RWlock.RUnlock()
									_,connbuffer2 = Buffer.PopFromQueue(connbuffer2)
									continue LOOP
								}
								continue LOOP
							}
						} else {
							time.Sleep(5*time.Millisecond)
							continue LOOP
						}
					default:
						temp = nil
						temp, connbuffer2 = Buffer.PopFromQueue(connbuffer2)
						if temp != nil {
							bufferchangeBack2 <- connbuffer2
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
							Global.Connstruct.RWlock.RLock()
							//fmt.Println("value22222:",value,value.RoundNum,Global.Connstruct.RoundNum,value.DataType)
							if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
								Global.Connstruct.RWlock.RUnlock()
								Global.AllDataSlice <- *value
								break LOOP
							} else {
								if value.RoundNum < Global.Connstruct.RoundNum {
									Global.Connstruct.RWlock.RUnlock()
									_,connbuffer2 = Buffer.PopFromQueue(connbuffer2)
									continue LOOP
								}
								continue LOOP
							}
						}else {
							time.Sleep(5*time.Millisecond)
							continue LOOP
						}
					}
				}
			}
		}
	}
}