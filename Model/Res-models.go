package Model

type Res struct {
	Result        string     `json:"result"`
	Content       []Cnt      `json:"content"`
}

type Cnt struct {
	UserID        string     `json:"UserID"`
	Opinions      []Opinion  `json:"Opinions"`
}
