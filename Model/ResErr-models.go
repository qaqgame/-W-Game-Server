package Model

type ResErr struct {
	Result     string     `json:"result"`
	Content    ErrCnt     `json:"content"`
}

type ErrCnt struct {
	ErrorID    string     `json:"errorID"`
	ErrorMsg   string     `json:"errorMsg"`
}
