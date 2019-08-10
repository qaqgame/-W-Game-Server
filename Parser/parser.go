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
		fmt.Println("json error: ",string(data))
		//loginfo := Log.GetTransferInfo()
		//Global.DebugLogger <- loginfo + err.Error()
		return nil
	}
	return &request
}
