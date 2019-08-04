package Model

type Req struct {
	UserID         string      `json:"UserID"`
	Opinions       []Opinion   `json:"Opinions"`
	Hash           HashValue   `json:"hash"`
}

type Opinion struct {
	Type           string      `json:"type"`
	Desc           string      `json:"desc"`
	Target         string      `json:"target"`
	Pos            string      `json:"pos"`
}

type HashValue struct {
	Player1Hash    string      `json:"player1-hash"`
	Player2Hash    string      `json:"player2-hash"`
	Player3Hash    string      `json:"player3-hash"`
	Player4Hash    string      `json:"player4-hash"`
	Player5Hash    string      `json:"player5-hash"`
}
