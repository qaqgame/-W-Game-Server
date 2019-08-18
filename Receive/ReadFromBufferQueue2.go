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
//LOOP:
//	for true {
//		select {
//		case <-cxt.Done():
//			fmt.Println("return out3")
//			return
//		default:
//			select {
//			case <-cxt.Done():
//				fmt.Println("return out4")
//				return
//			//case data := <-bufferchange2:
//			//	connbuffer2 = data
//			//	//fmt.Println("bufferlength2",connbuffer2.Size)
//			//	select {
//			//	case <-Global.Connstruct.PlayersChannelAck[remoteaddr]:
//			//	LOOP1:
//			//		for true {
//			//			temp = nil
//			//			temp = Buffer.PopFromQueue(connbuffer2)
//			//			//if connbuffer.Size > 0 {
//			//			//	fmt.Println("success > 0")
//			//			//}
//			//			if temp != nil {
//			//				bufferchangeBack2 <- connbuffer2
//			//				value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
//			//				if value != nil {
//			//					Global.Forwardtimer <- 1
//			//				}
//			//				Global.Connstruct.RWlock.RLock()
//			//				//fmt.Println("value11111:",value,value.RoundNum,Global.Connstruct.RoundNum,value.DataType)
//			//				if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
//			//					//fmt.Println("112211")
//			//					Global.Connstruct.RWlock.RUnlock()
//			//					//fmt.Println("221122")
//			//					Global.AllDataSlice <- *value
//			//					continue LOOP
//			//				} else {
//			//					if value.RoundNum < Global.Connstruct.RoundNum {
//			//						Global.Connstruct.RWlock.RUnlock()
//			//						_ = Buffer.PopFromQueue(connbuffer2)
//			//						bufferchangeBack2<-connbuffer2
//			//						continue LOOP1
//			//					}
//			//					continue LOOP1
//			//				}
//			//			}else {
//			//				//time.Sleep(5*time.Millisecond)
//			//				continue LOOP1
//			//			}
//			//		}
//			//
//			//	default:
//			//		continue LOOP
//			//	}
//			default:
//				select {
//				case <-Global.Connstruct.PlayersChannelAck[remoteaddr]:
//				LOOP2:
//					for true {
//						temp = nil
//						temp= Buffer.PopFromQueue(connbuffer2)
//						//if connbuffer.Size > 0 {
//						//	fmt.Println("success > 0")
//						//}
//						if temp != nil {
//							//bufferchangeBack2 <- connbuffer2
//							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
//							//if value != nil {
//							//	Global.Forwardtimer <- 1
//							//}
//							Global.Connstruct.RWlock.RLock()
//							//fmt.Println("value11111:",value,value.RoundNum,Global.Connstruct.RoundNum,value.DataType)
//							if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
//								//fmt.Println("112211")
//								Global.Connstruct.RWlock.RUnlock()
//								//fmt.Println("221122")
//								Global.AllDataSlice <- *value
//								continue LOOP
//							} else {
//								if value.RoundNum < Global.Connstruct.RoundNum {
//									Global.Connstruct.RWlock.RUnlock()
//									_ = Buffer.PopFromQueue(connbuffer2)
//									//bufferchangeBack2<-connbuffer2
//									continue LOOP2
//								}
//								continue LOOP2
//							}
//						}else {
//							//time.Sleep(5*time.Millisecond)
//							continue LOOP2
//						}
//					}
//
//				default:
//					continue LOOP
//				}
//			}
//		}
//	}
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
					//case data := <-bufferchange2:
					//	temp = nil
					//	connbuffer2 = data
					//	temp,connbuffer2 = Buffer.PopFromQueue(connbuffer2)
					//	if temp != nil {
					//		bufferchangeBack2 <- connbuffer2
					//		value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
					//		Global.Connstruct.RWlock.RLock()
					//		if value != nil {
					//			Global.Forwardtimer <- 1
					//		}
					//		fmt.Println("value22222:",value.RoundNum,Global.Connstruct.RoundNum,value.DataType)
					//		if value != nil && value.RoundNum == Global.Connstruct.RoundNum {
					//			Global.Connstruct.RWlock.RUnlock()
					//			Global.AllDataSlice <- *value
					//			break LOOP
					//		} else {
					//			if value.RoundNum < Global.Connstruct.RoundNum {
					//				Global.Connstruct.RWlock.RUnlock()
					//				_,connbuffer2 = Buffer.PopFromQueue(connbuffer2)
					//				continue LOOP
					//			}
					//			continue LOOP
					//		}
					//	} else {
					//		time.Sleep(5*time.Millisecond)
					//		continue LOOP
					//	}
					default:
						temp = nil
						temp = Buffer.PopFromQueue(connbuffer2)
						if temp != nil {
							//bufferchangeBack2 <- connbuffer2
							//fmt.Println("cnt",string(temp.Cnt))
							value := Buffer.ParserBufferQueue(temp.Cnt,remoteaddr)
							Global.Connstruct.RWlock.RLock()
							//if value != nil {
							//	Global.Forwardtimer <- 1
							//}
							//fmt.Println("value22222:",value,value.RoundNum,Global.Connstruct.RoundNum,value.DataType)
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
							//fmt.Println("bad end")
							//time.Sleep(5*time.Millisecond)
							continue LOOP
						}
					}
				}
			}
		}
	}
}