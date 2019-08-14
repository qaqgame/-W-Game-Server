package Model

type Req struct {
	DataType       int         `json:"datatype"`
	UserID         string      `json:"UserID"`
	RoundNum       int64       `json:"Roundnum"`
	Opinions       []Opinion   `json:"Opinions"`
	Hash           HashValue   `json:"hash"`
}

type Opinion struct {
	Type           string      `json:"type"`
	Desc           string      `json:"desc"`
	Target         string      `json:"target"`
	Pos            string      `json:"position"`
	FrameNum       int         `json:"framenum"`
}

type HashValue struct {
	Player1Hash    string      `json:"player1Hash"`
	Player2Hash    string      `json:"player2Hash"`
	Player3Hash    string      `json:"player3Hash"`
	Player4Hash    string      `json:"player4Hash"`
	Player5Hash    string      `json:"player5Hash"`
}

type ReqMini struct {
	DataType       int         `json:"datatype"`
	RoundNum       int64       `json:"Roundnum"`
}