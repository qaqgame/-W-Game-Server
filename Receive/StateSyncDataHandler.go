package Receive

import (
	"wGame/Global"
	"time"
	"fmt"
)

var StateSync        = make(chan []byte,5)
var bufdata          = make([][]byte,0)
var reconnTimer      = make(chan int,1)
var ReconnForWarding = make(chan int,1)
var timerEnd         = make(chan int,1)


func StateSyncDataHandler() {
	fmt.Println("run StatSyncDataHandler")
	count := 0
LOOP:
	for {
		select {
		case data := <-StateSync:
			fmt.Println("get signal ",Global.Connstruct.ConnCount)
			bufdata = append(bufdata, data)
			count++
			if count == Global.Connstruct.ConnCount {
				if checkValue(bufdata) {
					fmt.Println("send type7 signal")
					Global.ReConnData <- bufdata[0]
				} else {
					//错误情况处理
					fmt.Println("send type7 signal1")
					Global.ReConnData <- bufdata[0]
				}
				count = 0
				bufdata = nil
			}
			break LOOP
		case <-ReconnForWarding:
			if len(bufdata) == 0 || bufdata == nil {
				break
			}
			if checkValue(bufdata) {
				Global.ReConnData <- bufdata[0]
			} else {
				//错误情况处理
				Global.ReConnData <- bufdata[0]
			}
			count = 0
			bufdata = nil
		}
	}
}

func checkValue(v [][]byte) bool {
	temp := v[0]
	for _,d := range v {
		if ifEqual(d,temp) {
			continue
		} else {
			return false
		}
	}
	return true
}

func ifEqual(v1, v2 []byte) bool {
	if string(v1) == string(v2) {
		return true
	}
	return false
}

func timer() {
	fmt.Println("Reconn timer Open")
	timeDur := time.Duration(500*time.Millisecond)
	t := time.NewTimer(timeDur)
	defer t.Stop()
LOOP:
	for true {
		select {
		case <-timerEnd:
			break LOOP
		case <-reconnTimer:
			t.Reset(timeDur)
		case <-t.C:
			ReconnForWarding <- 1
		}
	}
}