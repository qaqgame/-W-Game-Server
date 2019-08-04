package Parser

import (
	"wGame/Model"
	"encoding/json"
	"fmt"
)

func ParserReq(data []byte) *Model.Req {
	var request Model.Req
	err := json.Unmarshal(data,&request)
	if err != nil {
		fmt.Println("Unmarshal JSON Error: ",err)
		return nil
	}
	return &request
}
