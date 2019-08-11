package Receive

import (
	"wGame/Buffer"
	"sync"
	"context"
	"wGame/Global"
	"time"
)

func ReadFromBufferQueue(cxt context.Context, remoteaddr string, top *Buffer.Node, size *int, bufferchange chan *Buffer.Node, bufferchangeBack chan *Buffer.Node, mutex *sync.Mutex) {
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
						top = data
						temp := Buffer.PopFromQueue(top,size,mutex)
						if temp != nil {
							bufferchangeBack <- temp[1]
							value := Buffer.ParserBufferQueue(temp[0].Cnt,remoteaddr)
							if value == nil {
								continue LOOP
							}
							Global.AllDataSlice <- *value
							break LOOP
						}else {
							time.Sleep(5*time.Millisecond)
							continue LOOP
						}
					default:
						temp := Buffer.PopFromQueue(top,size,mutex)
						if temp != nil {
							bufferchangeBack <- temp[1]
							value := Buffer.ParserBufferQueue(temp[0].Cnt,remoteaddr)
							if value == nil {
								continue LOOP
							}
							Global.AllDataSlice <- *value
							break LOOP
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
