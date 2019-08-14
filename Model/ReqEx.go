package Model

type ReqEx struct {
	RoundNum    int64
	DataType    int
	Request     Req
	UserId      string
	RemoteAddr  string
}
