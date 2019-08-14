package Model

type ResErr struct {
	DataType   int        `json:"datatype"`
	Result     string     `json:"result"`
	RoundNum   int64      `json:"Roundnum"`
	Content    ErrCnt     `json:"content"`
}

type ErrCnt struct {
	ErrorID    string     `json:"errorID"`
	ErrorMsg   string     `json:"errorMsg"`
}
