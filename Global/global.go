package Global

import (
	"wGame/Model"
	"net"
	"sync"
)

type Conns struct {
	RWlock             sync.RWMutex
	Conn               map[string]net.Conn
	PlayersChannel     map[string]chan int64
	PlayersChannelAck  map[string]chan int64           //读取buffer中reqack类型数据的信号
 	ConnCount          int
 	RoundNum           int64                         //当前回合数
	ConnStatus         map[string]int
	HaveReConn         bool                          //有重新连接请求
	StartStore         bool
	FlagRoundNum       int64
	ReconnData         []byte
}


var (
	AllDataSlice       chan Model.ReqEx               //客户端发送的req存储在这个channel中
	Forwardingsignal   chan int                       //转发signal的channel
	Forwardtimer       chan int                       //转发计时器channel
	ConnEstablish      chan int                       //连接建立channel
	Count              int                          //测试使用数据，记录发送的数据包次数
	PlayerNum          int                            //需要连接上的玩家数
	DebugLogger        chan string                    //输出到log文件的Logger channel
	Connstruct         Conns
	ReConnData         chan []byte
	Once               sync.Once
)

const LogFileName = "LogFile.log"

func init() {
	AllDataSlice     = make(chan Model.ReqEx, 5)
	Forwardingsignal = make(chan int, 1)
	Forwardtimer     = make(chan int, 1)
	ConnEstablish    = make(chan int, 1)
	PlayerNum        = 2
	DebugLogger      = make(chan string, 10)
	ReConnData       = make(chan []byte,1)

	Once             = sync.Once{}
	Connstruct.Conn  = make(map[string]net.Conn,5)
	Connstruct.PlayersChannel    = make(map[string]chan int64, 5)
	Connstruct.PlayersChannelAck = make(map[string]chan int64, 5)
	Connstruct.ConnCount = 0
	Connstruct.RoundNum  = 0
	Connstruct.RWlock    = sync.RWMutex{}
	Connstruct.ConnStatus= make(map[string]int, 5)
	Connstruct.HaveReConn=false
	Connstruct.StartStore=false
	Connstruct.FlagRoundNum=-1
	Connstruct.ReconnData = make([]byte,0)
}