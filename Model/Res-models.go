package Model

type Res struct {
	DataType      int        `json:"datatype"`
	Result        string     `json:"result"`
	RoundNum      int64      `json:"Roundnum"`
	Content       []Cnt      `json:"content"`
}

type Cnt struct {
	UserID        string     `json:"UserID"`
	Opinions      []Opinion  `json:"Opinions"`
}

type ReqAck struct {
	DataType      int        `json:"datatype"`
	Result        string     `json:"result"`
	RoundNum      int64      `json:"Roundnum"`
}