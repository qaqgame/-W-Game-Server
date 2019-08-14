package Receive

import (
	"wGame/Buffer"
	"context"
	"wGame/Global"
	"fmt"
)

func ReadFromBufferQueue(cxt context.Context, remoteaddr string,connbuffer *Buffer.ConnBuffer, bufferchange chan *Buffer.ConnBuffer, bufferchangeBack chan *Buffer.ConnBuffer) {
	var temp *Buffer.Node = nil
	if connbuffer == nil {
		fmt.Println("oooo")
	}
	for true {
		select {
		case <- Global.Connstruct.PlayersChannel[remoteaddr]:
		LOOP:
			for true {
				//退出线程
				select {
				case <-cxt.Done():
					return
				default:
					select {
					case <-cxt.Done():
						return
					case data := <-bufferchange:
						//fmt.Println("get")
						//fmt.Println("read1")
						temp = nil
						connbuffer = data
						temp,connbuffer = Buffer.PopFromQueue(connbuffer)
						//if connbuffer.Size > 0 {
						//	fmt.Println("success > 0")
						//}
						if temp != nil {
							bufferchangeBack <- connbuffer
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
							Global.Connstruct.RWlock.RLock()
							//fmt.Println("value11111:",value,value.RoundNum,Global.Connstruct.RoundNum,value.DataType)
							if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
								//fmt.Println("112211")
								Global.Connstruct.RWlock.RUnlock()
								//fmt.Println("221122")
								Global.AllDataSlice <- *value
								break LOOP
							} else {
								if value.RoundNum < Global.Connstruct.RoundNum {
									Global.Connstruct.RWlock.RUnlock()
									_,connbuffer = Buffer.PopFromQueue(connbuffer)
									continue LOOP
								}
								continue LOOP
							}
						}else {
							//time.Sleep(5*time.Millisecond)
							continue LOOP
						}
					default:
						//fmt.Println("get2",connbuffer)
						//fmt.Println("read2")
						temp = nil
						temp, connbuffer = Buffer.PopFromQueue(connbuffer)
						//if connbuffer.Size > 0 {
						//	fmt.Println("success > 0")
						//}
						if temp != nil {
							bufferchangeBack <- connbuffer
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
							Global.Connstruct.RWlock.RLock()
							//fmt.Println("value11111:",value)
							if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
								Global.Connstruct.RWlock.RUnlock()
								Global.AllDataSlice <- *value
								break LOOP
							} else {
								if value.RoundNum < Global.Connstruct.RoundNum {
									Global.Connstruct.RWlock.RUnlock()
									_,connbuffer = Buffer.PopFromQueue(connbuffer)
									continue LOOP
								}
								continue LOOP
							}
						}else {
							//time.Sleep(5*time.Millisecond)
							continue LOOP
						}
					}
				}
			}
		}
	}
}
