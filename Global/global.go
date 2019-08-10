package Global

import (
	"wGame/Model"
	"net"
)

//type Conns struct {
//	RWlock          sync.RWMutex
//	Conn            map[string]net.Conn
//	PlayerChannel   map[string]chan int
//  ConnCount       int
//}

var (
	AllDataSlice       chan Model.ReqEx               //客户端发送的req存储在这个channel中
	Conn               map[string]net.Conn            //存储conn的map
	Forwardingsignal   chan int                       //转发signal的channel
	Forwardtimer       chan int                       //转发计时器channel
	ConnCount          int                            //连接数
	ConnEstablish      chan int                       //连接建立channel
	Count              int                          //测试使用数据，记录发送的数据包次数
	PlayerNum          int                            //需要连接上的玩家数
	PlayersChannel     map[string]chan int            //每个玩家独自的channel
	DebugLogger        chan string                    //输出到log文件的Logger channel
	//Connstruct         Conns
)

const LogFileName = "LogFile.log"

func init() {
	AllDataSlice     = make(chan Model.ReqEx, 5)
	Conn             = make(map[string]net.Conn,5)
	Forwardingsignal = make(chan int, 1)
	Forwardtimer     = make(chan int, 1)
	ConnCount        = 0
	ConnEstablish    = make(chan int, 1)
	PlayerNum        = 5
	PlayersChannel   = make(map[string]chan int, 5)
	DebugLogger      = make(chan string, 10)
	//Connstruct.Conn  = Conn
	//Connstruct.PlayerChannel = PlayersChannel
}