package Receive

import (
	"wGame/Buffer"
	"context"
	"wGame/Global"
	"time"
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
						temp = nil
						connbuffer = data
						temp,connbuffer = Buffer.PopFromQueue(connbuffer)
						if temp != nil {
							bufferchangeBack <- connbuffer
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
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
						//fmt.Println("get2",connbuffer)
						temp = nil

						temp, connbuffer = Buffer.PopFromQueue(connbuffer)
						if temp != nil {
							bufferchangeBack <- connbuffer
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
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
