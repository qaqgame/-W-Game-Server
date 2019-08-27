package Model

//状态同步数据格式&&重连时发回的数据格式
//datatype=6       用户发送的状态同步信息
//datatype=7       服务端发送给客户端的重连数据
type StateSync struct {
	DataType    int        `json:"datatype"`
	Result      string     `json:"result"`
	RoundNum    int64      `json:"Roundnum"`
	AllStatus   []Status   `json:"AllStatus"`
}

//用户的状态信息
type Status struct {
	UserID      string     `json:"UserID"`
	Position    string     `json:"Position"`
}
