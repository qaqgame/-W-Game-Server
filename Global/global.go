package Global

import (
	"wGame/Model"
	"bufio"
)

var (
	AllDataSlice       chan Model.ReqEx               //客户端发送的req存储在这个channel中
	//Conns              []*bufio.ReadWriter          //存储连接的slice,考虑要不要换成map？？
	Conns              map[*bufio.ReadWriter]string   //存储连接的slice,考虑要不要换成map？？
	Forwardingsignal   chan int                       //转发signal的channel
	Forwardtimer       chan int                       //转发计时器channel
	ConnCount          int                            //连接数
	ConnEstablish      chan int                       //连接建立channel
	ConnBufferFull     chan int                       //缓冲区满
	Count int64
)

func init() {
	AllDataSlice     = make(chan Model.ReqEx,5)
	//Conns            = make([]*bufio.ReadWriter,0)
	Conns            = make(map[*bufio.ReadWriter]string,5)
	Forwardingsignal = make(chan int,1)
	Forwardtimer     = make(chan int,1)
	ConnCount        = 0
	ConnEstablish    = make(chan int,1)
	ConnBufferFull   = make(chan int,1)
}